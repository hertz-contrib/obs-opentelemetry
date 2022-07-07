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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
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

type Config struct {
	tracer trace.Tracer
	meter  metric.Meter

	tracerProvider    trace.TracerProvider
	meterProvider     metric.MeterProvider
	textMapPropagator propagation.TextMapPropagator

	recordSourceOperation bool
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
		tracerProvider:    otel.GetTracerProvider(),
		meterProvider:     global.MeterProvider(),
		textMapPropagator: otel.GetTextMapPropagator(),
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
