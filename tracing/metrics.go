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

// Server HTTP metrics. ref to https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/semantic_conventions/http-metrics.md#http-server
const (
	RequestCount          = "http.server.request_count"           // Incoming request count total
	RequestContentLength  = "http.server.request_content_length"  // Incoming request bytes total
	ResponseContentLength = "http.server.response_content_length" // Incoming response bytes total
	ServerLatency         = "http.server.duration"                // Incoming end to end duration, microseconds
)

// Client HTTP metrics. ref to https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/semantic_conventions/http-metrics.md#http-client
// http.client.duration	Histogram	milliseconds	ms	measures the duration outbound HTTP requests
// http.client.request.size	Histogram	bytes	By	measures the size of HTTP request messages (compressed)
// http.client.response.size	Histogram	bytes	By	measures the size of HTTP response messages (compressed)
const (
	ClientLatency      = "http.client.duration"
	ClientRequestSize  = "http.client.request.size"
	ClientResponseSize = "http.client.response.size"
)

var (
	HTTPMetricsAttributes = []attribute.Key{
		semconv.HTTPHostKey,
		semconv.HTTPRouteKey,
		semconv.HTTPMethodKey,
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
