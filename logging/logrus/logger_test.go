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

package logrus_test

import (
	"context"
	"testing"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	otelhertzlogrus "github.com/hertz-contrib/obs-opentelemetry/logging/logrus"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func stdoutProvider(ctx context.Context) func() {
	provider := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(provider)

	exp, err := stdouttrace.New()
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

func TestLogger(t *testing.T) {
	ctx := context.Background()
	shutdown := stdoutProvider(ctx)
	defer shutdown()

	logger := otelhertzlogrus.NewLogger(
		otelhertzlogrus.WithTraceHookErrorSpanLevel(logrus.WarnLevel),
		otelhertzlogrus.WithTraceHookLevels(logrus.AllLevels),
		otelhertzlogrus.WithRecordStackTraceInSpan(true),
	)

	logger.Logger().Info("log from origin logrus")

	hlog.SetLogger(logger)
	hlog.SetLevel(hlog.LevelDebug)

	tracer := otel.Tracer("test otel std logger")
	ctx, span := tracer.Start(ctx, "root")

	hlog.SetLogger(logger)
	hlog.SetLevel(hlog.LevelTrace)

	hlog.Trace("trace")
	hlog.Debug("debug")
	hlog.Info("info")
	hlog.Notice("notice")
	hlog.Warn("warn")
	hlog.Error("error")

	hlog.Tracef("log level: %s", "trace")
	hlog.Debugf("log level: %s", "debug")
	hlog.Infof("log level: %s", "info")
	hlog.Noticef("log level: %s", "notice")
	hlog.Warnf("log level: %s", "warn")
	hlog.Errorf("log level: %s", "error")

	hlog.CtxTracef(ctx, "log level: %s", "trace")
	hlog.CtxDebugf(ctx, "log level: %s", "debug")
	hlog.CtxInfof(ctx, "log level: %s", "info")
	hlog.CtxNoticef(ctx, "log level: %s", "notice")
	hlog.CtxWarnf(ctx, "log level: %s", "warn")
	hlog.CtxErrorf(ctx, "log level: %s", "error")

	span.End()

	ctx, child := tracer.Start(ctx, "child")
	hlog.CtxWarnf(ctx, "foo %s", "bar")
	child.End()

	ctx, errSpan := tracer.Start(ctx, "error")
	hlog.CtxErrorf(ctx, "error %s", "this is a error")
	hlog.Info("no trace context")
	errSpan.End()
}
