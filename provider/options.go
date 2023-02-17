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

package provider

import (
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/contrib/propagators/ot"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

// Option opts for opentelemetry tracer provider
type Option interface {
	apply(cfg *config)
}

type option func(cfg *config)

func (fn option) apply(cfg *config) {
	fn(cfg)
}

type config struct {
	enableTracing bool
	enableMetrics bool

	exportInsecure          bool
	exportEnableCompression bool
	exportEndpoint          string
	exportHeaders           map[string]string

	resource          *resource.Resource
	sdkTracerProvider *sdktrace.TracerProvider

	sampler sdktrace.Sampler

	resourceAttributes []attribute.KeyValue
	resourceDetectors  []resource.Detector

	textMapPropagator propagation.TextMapPropagator

	meterProvider *metric.MeterProvider
}

func newConfig(opts []Option) *config {
	cfg := defaultConfig()

	for _, opt := range opts {
		opt.apply(cfg)
	}

	return cfg
}

func defaultConfig() *config {
	return &config{
		enableTracing: true,
		enableMetrics: true,
		sampler:       sdktrace.AlwaysSample(),
		textMapPropagator: propagation.NewCompositeTextMapPropagator(
			propagation.NewCompositeTextMapPropagator(
				b3.New(),
				ot.OT{},
				propagation.Baggage{},
				propagation.TraceContext{},
			),
		),
	}
}

// WithServiceName configures `service.name` resource attribute
func WithServiceName(serviceName string) Option {
	return option(func(cfg *config) {
		cfg.resourceAttributes = append(cfg.resourceAttributes, semconv.ServiceNameKey.String(serviceName))
	})
}

// WithDeploymentEnvironment configures `deployment.environment` resource attribute
func WithDeploymentEnvironment(env string) Option {
	return option(func(cfg *config) {
		cfg.resourceAttributes = append(cfg.resourceAttributes, semconv.DeploymentEnvironmentKey.String(env))
	})
}

// WithServiceNamespace configures `service.namespace` resource attribute
func WithServiceNamespace(namespace string) Option {
	return option(func(cfg *config) {
		cfg.resourceAttributes = append(cfg.resourceAttributes, semconv.ServiceNamespaceKey.String(namespace))
	})
}

// WithResourceAttribute configures resource attribute.
func WithResourceAttribute(rAttr attribute.KeyValue) Option {
	return option(func(cfg *config) {
		cfg.resourceAttributes = append(cfg.resourceAttributes, rAttr)
	})
}

// WithResourceAttributes configures resource attributes.
func WithResourceAttributes(rAttrs []attribute.KeyValue) Option {
	return option(func(cfg *config) {
		cfg.resourceAttributes = rAttrs
	})
}

// WithResource configures resource
func WithResource(resource *resource.Resource) Option {
	return option(func(cfg *config) {
		cfg.resource = resource
	})
}

// WithExportEndpoint configures export endpoint
func WithExportEndpoint(endpoint string) Option {
	return option(func(cfg *config) {
		cfg.exportEndpoint = endpoint
	})
}

// WithEnableTracing enable tracing
func WithEnableTracing(enableTracing bool) Option {
	return option(func(cfg *config) {
		cfg.enableTracing = enableTracing
	})
}

// WithEnableMetrics enable metrics
func WithEnableMetrics(enableMetrics bool) Option {
	return option(func(cfg *config) {
		cfg.enableMetrics = enableMetrics
	})
}

// WithTextMapPropagator configures propagation
func WithTextMapPropagator(p propagation.TextMapPropagator) Option {
	return option(func(cfg *config) {
		cfg.textMapPropagator = p
	})
}

// WithResourceDetector configures resource detector
func WithResourceDetector(detector resource.Detector) Option {
	return option(func(cfg *config) {
		cfg.resourceDetectors = append(cfg.resourceDetectors, detector)
	})
}

// WithHeaders configures gRPC requests headers for exported telemetry data
func WithHeaders(headers map[string]string) Option {
	return option(func(cfg *config) {
		cfg.exportHeaders = headers
	})
}

// WithInsecure disables client transport security for the exporter's gRPC
func WithInsecure() Option {
	return option(func(cfg *config) {
		cfg.exportInsecure = true
	})
}

// WithEnableCompression enable gzip transport compression
func WithEnableCompression() Option {
	return option(func(cfg *config) {
		cfg.exportEnableCompression = true
	})
}

// WithSampler configures sampler
func WithSampler(sampler sdktrace.Sampler) Option {
	return option(func(cfg *config) {
		cfg.sampler = sampler
	})
}

// WithSdkTracerProvider configures sdkTracerProvider
func WithSdkTracerProvider(sdkTracerProvider *sdktrace.TracerProvider) Option {
	return option(func(cfg *config) {
		cfg.sdkTracerProvider = sdkTracerProvider
	})
}

// WithMeterProvider configures MeterProvider
func WithMeterProvider(meterProvider *metric.MeterProvider) Option {
	return option(func(cfg *config) {
		cfg.meterProvider = meterProvider
	})
}
