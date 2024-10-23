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
	"reflect"
	"testing"

	"github.com/cloudwego-contrib/cwgo-pkg/telemetry/instrumentation/otelhertz"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/contrib/propagators/ot"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func TestExtract(t *testing.T) {
	ctx := context.Background()
	bags, _ := baggage.Parse("foo=bar")
	ctx = baggage.ContextWithBaggage(ctx, bags)
	ctx = metainfo.WithValue(ctx, "foo", "bar")

	headers := &protocol.RequestHeader{}
	headers.Set("foo", "bar")

	type args struct {
		ctx      context.Context
		c        *Config
		metadata *protocol.RequestHeader
	}
	tests := []struct {
		name  string
		args  args
		want  baggage.Baggage
		want1 trace.SpanContext
	}{
		{
			name: "extract successful",
			args: args{
				ctx:      ctx,
				c:        otelhertz.DefaultConfig(),
				metadata: headers,
			},
			want:  bags,
			want1: trace.SpanContext{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Extract(tt.args.ctx, tt.args.c, tt.args.metadata)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Extract() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Extract() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestInject(t *testing.T) {
	cfg := otelhertz.NewConfig([]Option{WithTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		b3.New(),
		ot.OT{},
		propagation.Baggage{},
		propagation.TraceContext{},
	))}...)

	ctx := context.Background()

	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    [16]byte{1},
		SpanID:     [8]byte{2},
		TraceFlags: 0,
		TraceState: trace.TraceState{},
		Remote:     false,
	})

	ctx = trace.ContextWithSpanContext(ctx, spanContext)
	md := &protocol.RequestHeader{}

	type args struct {
		ctx      context.Context
		c        *Config
		metadata *protocol.RequestHeader
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "inject valid",
			args: args{
				ctx:      ctx,
				c:        cfg,
				metadata: md,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Inject(tt.args.ctx, tt.args.c, tt.args.metadata)
			assert.NotEmpty(t, tt.args.metadata)
			assert.Equal(t, "01000000000000000000000000000000-0200000000000000-0", md.Get("b3"))
			assert.Equal(t, "00-01000000000000000000000000000000-0200000000000000-00", md.Get("traceparent"))
		})
	}
}
