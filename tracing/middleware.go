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

	"github.com/cloudwego/hertz/pkg/common/tracer/stats"
	"go.opentelemetry.io/otel/attribute"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/hertz-contrib/obs-opentelemetry/tracing/internal"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type StringHeader protocol.RequestHeader

// Visit implements the metainfo.HTTPHeaderCarrier interface.
func (sh *StringHeader) Visit(f func(k, v string)) {
	(*protocol.RequestHeader)(sh).VisitAll(
		func(key, value []byte) {
			f(string(key), string(value))
		})
}

func ClientMiddleware(opts ...Option) client.Middleware {
	cfg := newConfig(opts)
	histogramRecorder := make(map[string]metric.Float64Histogram)
	counters := make(map[string]metric.Int64Counter)

	clientRequestCountMeasure, err := cfg.meter.Int64Counter(
		ClientRequestCount,
		metric.WithUnit("count"),
		metric.WithDescription("measures the client request count total"),
	)
	handleErr(err)

	clientLatencyMeasure, err := cfg.meter.Float64Histogram(
		ClientLatency,
		metric.WithUnit("ms"),
		metric.WithDescription("measures the duration outbound HTTP requests"),
	)
	handleErr(err)

	counters[ClientRequestCount] = clientRequestCountMeasure
	histogramRecorder[ClientLatency] = clientLatencyMeasure

	return func(next client.Endpoint) client.Endpoint {
		return func(ctx context.Context, req *protocol.Request, resp *protocol.Response) (err error) {
			if ctx == nil {
				ctx = context.Background()
			}

			start := time.Now()

			// trace start
			ctx, span := cfg.tracer.Start(
				ctx,
				clientSpanNaming(req),
				oteltrace.WithTimestamp(start),
				oteltrace.WithSpanKind(oteltrace.SpanKindClient),
			)
			defer span.End()

			// inject client service resource attributes (canonical service) to meta map
			md := injectPeerServiceToMetadata(ctx, span.(trace.ReadOnlySpan).Resource().Attributes())

			Inject(ctx, cfg, &req.Header)

			for k, v := range md {
				req.Header.Set(k, v)
			}

			err = next(ctx, req, resp)

			// end span
			if httpReq, err := adaptor.GetCompatRequest(req); err == nil {
				span.SetAttributes(semconv.NetAttributesFromHTTPRequest("tcp", httpReq)...)
				span.SetAttributes(semconv.EndUserAttributesFromHTTPRequest(httpReq)...)
				span.SetAttributes(semconv.HTTPServerAttributesFromHTTPRequest("", cfg.clientHttpRouteFormatter(req), httpReq)...)
			}

			// span attributes
			attrs := []attribute.KeyValue{
				semconv.HTTPURLKey.String(req.URI().String()),
				semconv.HTTPStatusCodeKey.Int(resp.StatusCode()),
			}
			span.SetAttributes(attrs...)

			// set span status with resp status code
			span.SetStatus(semconv.SpanStatusFromHTTPStatusCode(resp.StatusCode()))

			// extract metrics attr
			metricsAttributes := extractMetricsAttributesFromSpan(span)

			// record metrics
			counters[ClientRequestCount].Add(ctx, 1, metric.WithAttributes(metricsAttributes...))
			histogramRecorder[ClientLatency].Record(
				ctx,
				float64(time.Since(start))/float64(time.Millisecond),
				metric.WithAttributes(metricsAttributes...),
			)

			return
		}
	}
}

func ServerMiddleware(cfg *Config) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// get tracer carrier
		tc := internal.TraceCarrierFromContext(ctx)
		if tc == nil {
			hlog.CtxWarnf(ctx, "TraceCarrier not found in context")
			c.Next(ctx)
			return
		}

		sTracer := tc.Tracer()
		ti := c.GetTraceInfo()
		if ti.Stats().Level() == stats.LevelDisabled {
			c.Next(ctx)
			return
		}

		opts := []oteltrace.SpanStartOption{
			oteltrace.WithTimestamp(getStartTimeOrNow(ti)),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}

		peerServiceAttributes := extractPeerServiceAttributesFromMetadata(&c.Request.Header)

		// extract baggage and span context from header
		bags, spanCtx := Extract(ctx, cfg, &c.Request.Header)

		// set baggage
		ctx = baggage.ContextWithBaggage(ctx, bags)

		ctx, span := sTracer.Start(oteltrace.ContextWithRemoteSpanContext(ctx, spanCtx), serverSpanNaming(c), opts...)

		// peer service attributes
		span.SetAttributes(peerServiceAttributes...)

		// set span and attrs into tracer carrier for serverTracer finish
		tc.SetSpan(span)

		c.Next(ctx)

		if cfg.customResponseHandler != nil {
			// execute custom response handler
			cfg.customResponseHandler(ctx, c)
		}
	}
}
