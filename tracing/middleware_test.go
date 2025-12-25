// Copyright 2022 CloudWeGo Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tracing

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/tracer/stats"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func TestServerMiddleware(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	tracer, cfg := NewServerTracer(WithCustomResponseHandler(func(c context.Context, ctx *app.RequestContext) {
		ctx.Header("trace-id", oteltrace.SpanFromContext(c).SpanContext().TraceID().String())
	}))
	h := server.Default(tracer, server.WithHostPorts("127.0.0.1:6666"))
	h.Use(ServerMiddleware(cfg))
	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
	})

	go h.Spin()
	time.Sleep(100 * time.Millisecond)
	resp, err := http.Get("http://127.0.0.1:6666/ping")
	assert.Nil(t, err)
	assert.True(t, len(resp.Header.Get("trace-id")) != 0)
}

func TestServerMiddlewareDisableTrace(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	tracer, cfg := NewServerTracer(WithCustomResponseHandler(func(c context.Context, ctx *app.RequestContext) {
		ctx.Header("trace-id", oteltrace.SpanFromContext(c).SpanContext().TraceID().String())
	}))
	h := server.Default(tracer,
		server.WithHostPorts("127.0.0.1:16666"),
		server.WithTraceLevel(stats.LevelDisabled),
	)
	h.Use(ServerMiddleware(cfg))
	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
	})

	go h.Spin()
	time.Sleep(100 * time.Millisecond)
	resp, err := http.Get("http://127.0.0.1:16666/ping")
	assert.Nil(t, err)
	assert.True(t, len(resp.Header.Get("trace-id")) == 0)
}

// TestServerMiddlewareWithShouldIgnore tests that WithShouldIgnore option works correctly
// and ignored paths do not generate traces.
func TestServerMiddlewareWithShouldIgnore(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))

	tracer, cfg := NewServerTracer(
		WithShouldIgnore(func(ctx context.Context, c *app.RequestContext) bool {
			return strings.HasPrefix(string(c.Path()), "/health")
		}),
		WithCustomResponseHandler(func(c context.Context, ctx *app.RequestContext) {
			span := oteltrace.SpanFromContext(c)
			if span.SpanContext().IsValid() {
				ctx.Header("trace-id", span.SpanContext().TraceID().String())
			}
		}),
	)

	h := server.Default(tracer, server.WithHostPorts("127.0.0.1:26666"))
	h.Use(ServerMiddleware(cfg))
	h.GET("/health", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(200, "ok")
	})
	h.GET("/api/test", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(200, "test")
	})

	go h.Spin()
	time.Sleep(100 * time.Millisecond)

	// Request to ignored path should not have trace-id
	resp, err := http.Get("http://127.0.0.1:26666/health")
	assert.Nil(t, err)
	assert.True(t, len(resp.Header.Get("trace-id")) == 0)

	// Request to normal path should have trace-id
	resp, err = http.Get("http://127.0.0.1:26666/api/test")
	assert.Nil(t, err)
	assert.True(t, len(resp.Header.Get("trace-id")) != 0)
}

// TestServerTracerNoDataRace tests that serverTracer does not cause data races
// under concurrent requests. This test should be run with -race flag.
// It specifically tests the fix for:
// 1. shouldIgnore check removed from Start() method
// 2. extractMetricsAttributesFromSpan called before span.End()
func TestServerTracerNoDataRace(t *testing.T) {
	// Use a custom SpanProcessor that simulates async export behavior
	// to detect data races between extractMetricsAttributesFromSpan()
	// and the exporter accessing span attributes concurrently.
	exporter, err := stdouttrace.New(stdouttrace.WithWriter(io.Discard))
	if err != nil {
		t.Fatal(err)
	}

	// Use BatchSpanProcessor with very aggressive settings to maximize
	// the chance of race detection
	bsp := sdktrace.NewBatchSpanProcessor(exporter,
		sdktrace.WithBatchTimeout(1*time.Millisecond), // Very short timeout
		sdktrace.WithMaxExportBatchSize(1),            // Export immediately when 1 span is ready
		sdktrace.WithMaxQueueSize(1),                  // Small queue to force frequent exports
	)
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(bsp))
	otel.SetTracerProvider(tp)

	// Use a noop MeterProvider to avoid polluting metrics in other tests
	mp := sdkmetric.NewMeterProvider()
	otel.SetMeterProvider(mp)
	defer func() {
		_ = mp.Shutdown(context.Background())
	}()

	tracer, cfg := NewServerTracer(
		WithShouldIgnore(func(ctx context.Context, c *app.RequestContext) bool {
			path := string(c.Path())
			return strings.HasPrefix(path, "/health")
		}),
	)

	h := server.Default(tracer, server.WithHostPorts("127.0.0.1:36666"))
	h.Use(ServerMiddleware(cfg))
	h.GET("/health", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(200, "ok")
	})
	h.GET("/api/test", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(200, "test")
	})

	go h.Spin()
	time.Sleep(100 * time.Millisecond)

	// Send concurrent requests continuously while BatchSpanProcessor exports in background
	var wg sync.WaitGroup
	done := make(chan struct{})

	// Start multiple request goroutines
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				default:
					_, _ = http.Get("http://127.0.0.1:36666/api/test")
				}
			}
		}()
	}

	// Let the test run for enough time to trigger potential races
	time.Sleep(500 * time.Millisecond)
	close(done)
	wg.Wait()

	// Shutdown triggers final export, which may also race with any pending operations
	_ = tp.Shutdown(context.Background())
}
