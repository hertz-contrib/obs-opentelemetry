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

	"github.com/bytedance/gopkg/cloud/metainfo"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var _ propagation.TextMapCarrier = &metadataProvider{}

type metadataProvider struct {
	metadata map[string]string
}

// Get a value from metadata by key
func (m *metadataProvider) Get(key string) string {
	if v, ok := m.metadata[key]; ok {
		return v
	}
	return ""
}

// Set a value to metadata by k/v
func (m *metadataProvider) Set(key, value string) {
	m.metadata[key] = value
}

// Keys Iteratively get all keys of metadata
func (m *metadataProvider) Keys() []string {
	out := make([]string, 0, len(m.metadata))
	for k := range m.metadata {
		out = append(out, k)
	}
	return out
}

// Inject injects span context into the kitex metadata info
func Inject(ctx context.Context, c *Config, metadata map[string]string) {
	c.textMapPropagator.Inject(ctx, &metadataProvider{metadata: metadata})
}

// Extract returns the baggage and span context
func Extract(ctx context.Context, c *Config, metadata map[string]string) (baggage.Baggage, trace.SpanContext) {
	ctx = c.textMapPropagator.Extract(ctx, &metadataProvider{metadata: CGIVariableToHTTPHeaderMetadata(metadata)})
	return baggage.FromContext(ctx), trace.SpanContextFromContext(ctx)
}

// CGIVariableToHTTPHeaderMetadata converts all CGI variable into HTTP header key.
// For example, `ABC_DEF` will be converted to `abc-def`.
func CGIVariableToHTTPHeaderMetadata(metadata map[string]string) map[string]string {
	res := make(map[string]string, len(metadata))
	for k, v := range metadata {
		res[metainfo.CGIVariableToHTTPHeader(k)] = v
	}
	return res
}

// HTTPHeaderToCGIVariableMetadata performs an CGI variable conversion.
// For example, an HTTP header key `abc-def` will result in `ABC_DEF`.
func HTTPHeaderToCGIVariableMetadata(metadata map[string]string) map[string]string {
	res := make(map[string]string, len(metadata))
	for k, v := range metadata {
		res[metainfo.HTTPHeaderToCGIVariable(k)] = v
	}
	return res
}
