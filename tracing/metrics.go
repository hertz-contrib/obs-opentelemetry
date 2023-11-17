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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// Server HTTP metrics
const (
	ServerRequestCount = "http.server.request_count" // measures the incoming request count total
	ServerLatency      = "http.server.duration"      // measures th incoming end to end duration
)

// Client HTTP metrics.
const (
	ClientRequestCount = "http.client.request_count" // measures the client request count total
	ClientLatency      = "http.client.duration"      // measures the duration outbound HTTP requests
)

var (
	HTTPMetricsAttributes = []attribute.Key{
		semconv.HTTPHostKey,
		semconv.HTTPRouteKey,
		semconv.HTTPMethodKey,
		semconv.HTTPStatusCodeKey,
	}

	PeerMetricsAttributes = []attribute.Key{
		semconv.PeerServiceKey,
		PeerServiceNamespaceKey,
		PeerDeploymentEnvironmentKey,
		RequestProtocolKey,
	}

	// MetricResourceAttributes resource attributes
	MetricResourceAttributes = []attribute.Key{
		semconv.ServiceNameKey,
		semconv.ServiceNamespaceKey,
		semconv.DeploymentEnvironmentKey,
		semconv.ServiceInstanceIDKey,
		semconv.ServiceVersionKey,
		semconv.TelemetrySDKLanguageKey,
		semconv.TelemetrySDKVersionKey,
		semconv.ProcessPIDKey,
		semconv.HostNameKey,
		semconv.HostIDKey,
	}
)

func extractMetricsAttributesFromSpan(span oteltrace.Span) []attribute.KeyValue {
	var attrs []attribute.KeyValue
	readOnlySpan, ok := span.(trace.ReadOnlySpan)
	if !ok {
		return attrs
	}

	// span attributes
	for _, attr := range readOnlySpan.Attributes() {
		if matchAttributeKey(attr.Key, HTTPMetricsAttributes) {
			attrs = append(attrs, attr)
		}
		if matchAttributeKey(attr.Key, PeerMetricsAttributes) {
			attrs = append(attrs, attr)
		}
	}

	// span resource attributes
	for _, attr := range readOnlySpan.Resource().Attributes() {
		if matchAttributeKey(attr.Key, MetricResourceAttributes) {
			attrs = append(attrs, attr)
		}
	}

	// status code
	attrs = append(attrs, StatusKey.String(readOnlySpan.Status().Code.String()))

	return attrs
}

func matchAttributeKey(key attribute.Key, toMatchKeys []attribute.Key) bool {
	for _, attrKey := range toMatchKeys {
		if attrKey == key {
			return true
		}
	}
	return false
}
