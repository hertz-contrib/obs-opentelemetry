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

	"github.com/cloudwego/hertz/pkg/protocol"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var _ propagation.TextMapCarrier = &metadataProvider{}

type metadataProvider struct {
	metadata map[string]string
	headers  *protocol.RequestHeader
}

// Get a value from metadata by key
func (m *metadataProvider) Get(key string) string {
	return m.headers.Get(key)
}

// Set a value to metadata by k/v
func (m *metadataProvider) Set(key, value string) {
	m.headers.Set(key, value)
}

// Keys Iteratively get all keys of metadata
func (m *metadataProvider) Keys() []string {
	out := make([]string, 0, len(m.metadata))

	m.headers.VisitAll(func(key, value []byte) {
		out = append(out, string(key))
	})

	return out
}

// Inject injects span context into the hertz metadata info
func Inject(ctx context.Context, c *Config, headers *protocol.RequestHeader) {
	c.GetTextMapPropagator().Inject(ctx, &metadataProvider{headers: headers})
}

// Extract returns the baggage and span context
func Extract(ctx context.Context, c *Config, headers *protocol.RequestHeader) (baggage.Baggage, trace.SpanContext) {
	ctx = c.GetTextMapPropagator().Extract(ctx, &metadataProvider{headers: headers})
	return baggage.FromContext(ctx), trace.SpanContextFromContext(ctx)
}
