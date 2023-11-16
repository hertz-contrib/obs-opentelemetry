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
	"time"

	"go.opentelemetry.io/otel/metric"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
	serverconfig "github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/tracer"
	"github.com/cloudwego/hertz/pkg/common/tracer/stats"
	"github.com/hertz-contrib/obs-opentelemetry/tracing/internal"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var _ tracer.Tracer = (*serverTracer)(nil)

type serverTracer struct {
	config            *Config
	counters          map[string]metric.Int64Counter
	histogramRecorder map[string]metric.Float64Histogram
}

func NewServerTracer(opts ...Option) (serverconfig.Option, *Config) {
	cfg := newConfig(opts)
	st := &serverTracer{
		config:            cfg,
		counters:          make(map[string]metric.Int64Counter),
		histogramRecorder: make(map[string]metric.Float64Histogram),
	}

	st.createMeasures()

	return server.WithTracer(st), cfg
}

func (s *serverTracer) createMeasures() {
	serverRequestCountMeasure, err := s.config.meter.Int64Counter(
		ServerRequestCount,
		metric.WithUnit("count"),
		metric.WithDescription("measures Incoming request count total"),
	)
	handleErr(err)

	serverLatencyMeasure, err := s.config.meter.Float64Histogram(
		ServerLatency,
		metric.WithUnit("ms"),
		metric.WithDescription("measures th incoming end to end duration"),
	)
	handleErr(err)

	s.counters[ServerRequestCount] = serverRequestCountMeasure
	s.histogramRecorder[ServerLatency] = serverLatencyMeasure
}

func (s *serverTracer) Start(ctx context.Context, _ *app.RequestContext) context.Context {
	tc := &internal.TraceCarrier{}
	tc.SetTracer(s.config.tracer)

	return internal.WithTraceCarrier(ctx, tc)
}

func (s *serverTracer) Finish(ctx context.Context, c *app.RequestContext) {
	// trace carrier from context
	tc := internal.TraceCarrierFromContext(ctx)
	if tc == nil {
		hlog.Warnf("get tracer container failed")
		return
	}

	ti := c.GetTraceInfo()
	st := ti.Stats()

	if st.Level() == stats.LevelDisabled {
		return
	}

	httpStart := st.GetEvent(stats.HTTPStart)
	if httpStart == nil {
		return
	}

	elapsedTime := float64(st.GetEvent(stats.HTTPFinish).Time().Sub(httpStart.Time())) / float64(time.Millisecond)

	// span
	span := tc.Span()
	if span == nil || !span.IsRecording() {
		return
	}

	// span attributes from original http request
	if httpReq, err := adaptor.GetCompatRequest(c.GetRequest()); err == nil {
		span.SetAttributes(semconv.NetAttributesFromHTTPRequest("tcp", httpReq)...)
		span.SetAttributes(semconv.EndUserAttributesFromHTTPRequest(httpReq)...)
		span.SetAttributes(semconv.HTTPServerAttributesFromHTTPRequest("", s.config.serverHttpRouteFormatter(c), httpReq)...)
		span.SetStatus(semconv.SpanStatusFromHTTPStatusCode(c.Response.StatusCode()))
	}

	injectStatsEventsToSpan(span, st)

	if panicMsg, panicStack, httpErr := parseHTTPError(ti); httpErr != nil || len(panicMsg) > 0 {
		recordErrorSpanWithStack(span, httpErr, panicMsg, panicStack)
	}

	span.End(oteltrace.WithTimestamp(getEndTimeOrNow(ti)))

	metricsAttributes := extractMetricsAttributesFromSpan(span)
	s.counters[ServerRequestCount].Add(ctx, 1, metric.WithAttributes(metricsAttributes...))
	s.histogramRecorder[ServerLatency].Record(ctx, elapsedTime, metric.WithAttributes(metricsAttributes...))
}
