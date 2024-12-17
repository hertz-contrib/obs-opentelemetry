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

package zap

import (
	"bytes"
	"context"
	"testing"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.uber.org/zap"
)

func noopProvider(ctx context.Context) func() {
	provider := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(provider)

	exp := tracetest.NewNoopExporter()

	bsp := sdktrace.NewSimpleSpanProcessor(exp)
	provider.RegisterSpanProcessor(bsp)

	return func() {
		if err := provider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}
}

func stdoutProvider(ctx context.Context) func() {
	provider := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(provider)

	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		panic(err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exp)
	provider.RegisterSpanProcessor(bsp)

	return func() {
		if err := provider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}
}

// TestLogger test logger work with opentelemetry
func TestLogger(t *testing.T) {
	ctx := context.Background()

	buf := new(bytes.Buffer)

	shutdown := stdoutProvider(ctx)
	defer shutdown()

	logger := NewLogger(
		WithTraceErrorSpanLevel(zap.WarnLevel),
		WithRecordStackTraceInSpan(true),
	)
	defer logger.Sync()

	hlog.SetLogger(logger)
	hlog.SetOutput(buf)
	hlog.SetLevel(hlog.LevelDebug)

	logger.Info("log from origin zap")
	assert.Contains(t, buf.String(), "log from origin zap")
	buf.Reset()

	tracer := otel.Tracer("test otel std logger")

	ctx, span := tracer.Start(ctx, "root")

	hlog.CtxInfof(ctx, "hello %s", "world")
	assert.Contains(t, buf.String(), "trace_id")
	assert.Contains(t, buf.String(), "span_id")
	assert.Contains(t, buf.String(), "trace_flags")
	buf.Reset()

	span.End()

	ctx, child1 := tracer.Start(ctx, "child1")

	hlog.CtxTracef(ctx, "trace %s", "this is a trace log")
	hlog.CtxDebugf(ctx, "debug %s", "this is a debug log")
	hlog.CtxInfof(ctx, "info %s", "this is a info log")

	child1.End()
	assert.Equal(t, codes.Unset, child1.(sdktrace.ReadOnlySpan).Status().Code)

	ctx, child2 := tracer.Start(ctx, "child2")
	hlog.CtxNoticef(ctx, "notice %s", "this is a notice log")
	hlog.CtxWarnf(ctx, "warn %s", "this is a warn log")
	hlog.CtxErrorf(ctx, "error %s", "this is a error log")

	child2.End()
	assert.Equal(t, codes.Error, child2.(sdktrace.ReadOnlySpan).Status().Code)

	_, errSpan := tracer.Start(ctx, "error")

	hlog.Info("no trace context")

	errSpan.End()
}

// TestLogLevel test SetLevel
func TestLogLevel(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewLogger(
		WithTraceErrorSpanLevel(zap.WarnLevel),
		WithRecordStackTraceInSpan(true),
	)
	defer logger.Sync()

	// output to buffer
	logger.SetOutput(buf)

	logger.Debug("this is a debug log")
	assert.NotContains(t, buf.String(), "this is a debug log")

	logger.SetLevel(hlog.LevelDebug)

	logger.Debugf("this is a debug log %s", "msg")
	assert.Contains(t, buf.String(), "this is a debug log")
}

func TestWithExtraKeys(t *testing.T) {
	ctx := context.Background()

	buf := new(bytes.Buffer)

	shutdown := noopProvider(ctx)
	defer shutdown()

	logger := NewLogger(
		WithTraceErrorSpanLevel(zap.WarnLevel),
		WithLogger(hertzzap.NewLogger(hertzzap.WithExtraKeys([]hertzzap.ExtraKey{"logger"}))),
	)
	defer logger.Sync()

	hlog.SetLogger(logger)
	hlog.SetOutput(buf)

	tracer := otel.Tracer("test otel std logger")

	ctx, span := tracer.Start(ctx, "root")
	ctx = context.WithValue(ctx, hertzzap.ExtraKey("logger"), "zap")

	hlog.CtxInfof(ctx, "info %s", "this is a info log")
	assert.Contains(t, buf.String(), "\"logger\":\"zap\"")

	span.End()
}
