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

	"github.com/cloudwego-contrib/cwgo-pkg/telemetry/instrumentation/otelhertz"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"go.opentelemetry.io/otel/propagation"
)

// Option opts for opentelemetry tracer provider
type Option = otelhertz.Option

type ConditionFunc func(ctx context.Context, c *app.RequestContext) bool

type Config = otelhertz.Config

// WithRecordSourceOperation configures record source operation dimension
func WithRecordSourceOperation(recordSourceOperation bool) Option {
	return otelhertz.WithRecordSourceOperation(recordSourceOperation)
}

// WithTextMapPropagator configures propagation
func WithTextMapPropagator(p propagation.TextMapPropagator) Option {
	return otelhertz.WithTextMapPropagator(p)
}

// WithCustomResponseHandler configures CustomResponseHandler
func WithCustomResponseHandler(h app.HandlerFunc) Option {
	return otelhertz.WithCustomResponseHandler(h)
}

// WithClientHttpRouteFormatter configures clientHttpRouteFormatter
func WithClientHttpRouteFormatter(clientHttpRouteFormatter func(req *protocol.Request) string) Option {
	return otelhertz.WithClientHttpRouteFormatter(clientHttpRouteFormatter)
}

// WithServerHttpRouteFormatter configures serverHttpRouteFormatter
func WithServerHttpRouteFormatter(serverHttpRouteFormatter func(c *app.RequestContext) string) Option {
	return otelhertz.WithServerHttpRouteFormatter(serverHttpRouteFormatter)
}

// WithClientSpanNameFormatter configures clientSpanNameFormatter
func WithClientSpanNameFormatter(clientSpanNameFormatter func(req *protocol.Request) string) Option {
	return otelhertz.WithClientSpanNameFormatter(clientSpanNameFormatter)
}

// WithServerSpanNameFormatter configures serverSpanNameFormatter
func WithServerSpanNameFormatter(serverSpanNameFormatter func(c *app.RequestContext) string) Option {
	return otelhertz.WithServerSpanNameFormatter(serverSpanNameFormatter)
}

// WithShouldIgnore allows you to define the condition for enabling distributed tracing
func WithShouldIgnore(condition ConditionFunc) Option {
	return otelhertz.WithShouldIgnore(otelhertz.ConditionFunc(condition))
}
