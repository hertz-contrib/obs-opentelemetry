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

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName = "github.com/hertz-contrib/obs-opentelemetry"
)

// Option opts for opentelemetry tracer provider
type Option interface {
	apply(cfg *Config)
}

type option func(cfg *Config)

func (fn option) apply(cfg *Config) {
	fn(cfg)
}

type ConditionFunc func(ctx context.Context, c *app.RequestContext) bool

type Config struct {
	tracer trace.Tracer
	meter  metric.Meter

	clientHttpRouteFormatter func(req *protocol.Request) string
	serverHttpRouteFormatter func(c *app.RequestContext) string

	clientSpanNameFormatter func(req *protocol.Request) string
	serverSpanNameFormatter func(c *app.RequestContext) string

	tracerProvider    trace.TracerProvider
	meterProvider     metric.MeterProvider
	textMapPropagator propagation.TextMapPropagator

	recordSourceOperation bool

	customResponseHandler app.HandlerFunc
	shouldIgnore          ConditionFunc
}

func newConfig(opts []Option) *Config {
	cfg := defaultConfig()

	for _, opt := range opts {
		opt.apply(cfg)
	}

	cfg.meter = cfg.meterProvider.Meter(
		instrumentationName,
		metric.WithInstrumentationVersion(SemVersion()),
	)

	cfg.tracer = cfg.tracerProvider.Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion(SemVersion()),
	)

	return cfg
}

func defaultConfig() *Config {
	return &Config{
		tracerProvider:        otel.GetTracerProvider(),
		meterProvider:         otel.GetMeterProvider(),
		textMapPropagator:     otel.GetTextMapPropagator(),
		customResponseHandler: func(c context.Context, ctx *app.RequestContext) {},
		clientHttpRouteFormatter: func(req *protocol.Request) string {
			return string(req.Path())
		},
		clientSpanNameFormatter: func(req *protocol.Request) string {
			return string(req.Method()) + " " + string(req.Path())
		},
		serverHttpRouteFormatter: func(c *app.RequestContext) string {
			// FullPath returns a matched route full path. For not found routes
			// returns an empty string.
			route := c.FullPath()
			// fall back to path
			if route == "" {
				route = string(c.Path())
			}
			return route
		},
		serverSpanNameFormatter: func(c *app.RequestContext) string {
			// Ref to https://github.com/open-telemetry/opentelemetry-specification/blob/ffddc289462dfe0c2041e3ca42a7b1df805706de/specification/trace/api.md#span
			// FullPath returns a matched route full path. For not found routes
			// returns an empty string.
			route := c.FullPath()
			// fall back to handler name
			if route == "" {
				route = string(c.Path())
			}
			return string(c.Method()) + " " + route
		},
		shouldIgnore: func(ctx context.Context, c *app.RequestContext) bool {
			return false
		},
	}
}

// WithRecordSourceOperation configures record source operation dimension
func WithRecordSourceOperation(recordSourceOperation bool) Option {
	return option(func(cfg *Config) {
		cfg.recordSourceOperation = recordSourceOperation
	})
}

// WithTextMapPropagator configures propagation
func WithTextMapPropagator(p propagation.TextMapPropagator) Option {
	return option(func(cfg *Config) {
		cfg.textMapPropagator = p
	})
}

// WithCustomResponseHandler configures CustomResponseHandler
func WithCustomResponseHandler(h app.HandlerFunc) Option {
	return option(func(cfg *Config) {
		cfg.customResponseHandler = h
	})
}

// WithClientHttpRouteFormatter configures clientHttpRouteFormatter
func WithClientHttpRouteFormatter(clientHttpRouteFormatter func(req *protocol.Request) string) Option {
	return option(func(cfg *Config) {
		cfg.clientHttpRouteFormatter = clientHttpRouteFormatter
	})
}

// WithServerHttpRouteFormatter configures serverHttpRouteFormatter
func WithServerHttpRouteFormatter(serverHttpRouteFormatter func(c *app.RequestContext) string) Option {
	return option(func(cfg *Config) {
		cfg.serverHttpRouteFormatter = serverHttpRouteFormatter
	})
}

// WithClientSpanNameFormatter configures clientSpanNameFormatter
func WithClientSpanNameFormatter(clientSpanNameFormatter func(req *protocol.Request) string) Option {
	return option(func(cfg *Config) {
		cfg.clientSpanNameFormatter = clientSpanNameFormatter
	})
}

// WithServerSpanNameFormatter configures serverSpanNameFormatter
func WithServerSpanNameFormatter(serverSpanNameFormatter func(c *app.RequestContext) string) Option {
	return option(func(cfg *Config) {
		cfg.serverSpanNameFormatter = serverSpanNameFormatter
	})
}

// WithShouldIgnore allows you to define the condition for enabling distributed tracing
func WithShouldIgnore(condition ConditionFunc) Option {
	return option(func(cfg *Config) {
		cfg.shouldIgnore = condition
	})
}
