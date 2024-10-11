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
	"net/http"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/tracer/stats"
	"go.opentelemetry.io/otel"
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
