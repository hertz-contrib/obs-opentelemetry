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
	"github.com/cloudwego-contrib/cwgo-pkg/telemetry/provider/otelprovider"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Option opts for opentelemetry tracer provider
type Option = otelprovider.Option

// WithServiceName configures `service.name` resource attribute
func WithServiceName(serviceName string) Option {
	return otelprovider.WithServiceName(serviceName)
}

// WithDeploymentEnvironment configures `deployment.environment` resource attribute
func WithDeploymentEnvironment(env string) Option {
	return otelprovider.WithDeploymentEnvironment(env)
}

// WithServiceNamespace configures `service.namespace` resource attribute
func WithServiceNamespace(namespace string) Option {
	return otelprovider.WithServiceNamespace(namespace)
}

// WithResourceAttribute configures resource attribute.
func WithResourceAttribute(rAttr attribute.KeyValue) Option {
	return otelprovider.WithResourceAttribute(rAttr)
}

// WithResourceAttributes configures resource attributes.
func WithResourceAttributes(rAttrs []attribute.KeyValue) Option {
	return otelprovider.WithResourceAttributes(rAttrs)
}

// WithResource configures resource
func WithResource(resource *resource.Resource) Option {
	return otelprovider.WithResource(resource)
}

// WithExportEndpoint configures export endpoint
func WithExportEndpoint(endpoint string) Option {
	return otelprovider.WithExportEndpoint(endpoint)
}

// WithEnableTracing enable tracing
func WithEnableTracing(enableTracing bool) Option {
	return otelprovider.WithEnableTracing(enableTracing)
}

// WithEnableMetrics enable metrics
func WithEnableMetrics(enableMetrics bool) Option {
	return otelprovider.WithEnableMetrics(enableMetrics)
}

// WithTextMapPropagator configures propagation
func WithTextMapPropagator(p propagation.TextMapPropagator) Option {
	return otelprovider.WithTextMapPropagator(p)
}

// WithResourceDetector configures resource detector
func WithResourceDetector(detector resource.Detector) Option {
	return otelprovider.WithResourceDetector(detector)
}

// WithHeaders configures gRPC requests headers for exported telemetry data
func WithHeaders(headers map[string]string) Option {
	return otelprovider.WithHeaders(headers)
}

// WithInsecure disables client transport security for the exporter's gRPC
func WithInsecure() Option {
	return otelprovider.WithInsecure()
}

// WithEnableCompression enable gzip transport compression
func WithEnableCompression() Option {
	return otelprovider.WithEnableCompression()
}

// WithSampler configures sampler
func WithSampler(sampler sdktrace.Sampler) Option {
	return otelprovider.WithSampler(sampler)
}

// WithSdkTracerProvider configures sdkTracerProvider
func WithSdkTracerProvider(sdkTracerProvider *sdktrace.TracerProvider) Option {
	return otelprovider.WithSdkTracerProvider(sdkTracerProvider)
}

// WithMeterProvider configures MeterProvider
func WithMeterProvider(meterProvider *metric.MeterProvider) Option {
	return otelprovider.WithMeterProvider(meterProvider)
}
