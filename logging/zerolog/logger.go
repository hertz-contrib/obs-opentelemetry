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

package zerolog

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzzerolog "github.com/hertz-contrib/logger/zerolog"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Logger struct {
	hertzzerolog.Logger
	config *config
}

// Ref to https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/logs/README.md#json-formats
const (
	traceIDKey    = "trace_id"
	spanIDKey     = "span_id"
	traceFlagsKey = "trace_flags"
)

type ExtraKey string

var extraKeys = []ExtraKey{traceIDKey, spanIDKey, traceFlagsKey}

func NewLogger(opts ...Option) *Logger {
	config := defaultConfig()

	// apply options
	for _, opt := range opts {
		opt.apply(config)
	}
	logger := *config.logger
	logger.Unwrap().Hook(zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, message string) {
		ctx := e.GetCtx()
		e.Any(traceIDKey, ctx.Value(ExtraKey(traceIDKey)))
		e.Any(spanIDKey, ctx.Value(ExtraKey(spanIDKey)))
		e.Any(traceFlagsKey, ctx.Value(ExtraKey(traceFlagsKey)))
	}))

	return &Logger{
		Logger: logger,
		config: config,
	}
}

func (l *Logger) CtxLogf(level hlog.Level, ctx context.Context, format string, kvs ...any) {
	var zlevel zerolog.Level
	span := trace.SpanFromContext(ctx)

	ctx = context.WithValue(ctx, ExtraKey(traceIDKey), span.SpanContext().TraceID())
	ctx = context.WithValue(ctx, ExtraKey(spanIDKey), span.SpanContext().SpanID())
	ctx = context.WithValue(ctx, ExtraKey(traceFlagsKey), span.SpanContext().TraceFlags())

	switch level {
	case hlog.LevelDebug, hlog.LevelTrace:
		zlevel = zerolog.DebugLevel
		l.Logger.CtxDebugf(ctx, format, kvs...)
	case hlog.LevelInfo:
		zlevel = zerolog.InfoLevel
		l.Logger.CtxInfof(ctx, format, kvs...)
	case hlog.LevelNotice, hlog.LevelWarn:
		zlevel = zerolog.WarnLevel
		l.Logger.CtxWarnf(ctx, format, kvs...)
	case hlog.LevelError:
		zlevel = zerolog.ErrorLevel
		l.Logger.CtxErrorf(ctx, format, kvs...)
	case hlog.LevelFatal:
		zlevel = zerolog.FatalLevel
		l.Logger.CtxFatalf(ctx, format, kvs...)
	default:
		zlevel = zerolog.WarnLevel
		l.Logger.CtxWarnf(ctx, format, kvs...)
	}

	if !span.IsRecording() {
		l.Logger.Logf(level, format, kvs...)
		return
	}

	// set span status
	if zlevel >= l.config.traceConfig.errorSpanLevel {
		msg := getMessage(format, kvs)
		span.SetStatus(codes.Error, "")
		span.RecordError(errors.New(msg), trace.WithStackTrace(l.config.traceConfig.recordStackTraceInSpan))
	}
}

func (l *Logger) CtxTracef(ctx context.Context, format string, v ...any) {
	l.CtxLogf(hlog.LevelDebug, ctx, format, v...)
}

func (l *Logger) CtxDebugf(ctx context.Context, format string, v ...any) {
	l.CtxLogf(hlog.LevelDebug, ctx, format, v...)
}

func (l *Logger) CtxInfof(ctx context.Context, format string, v ...any) {
	l.CtxLogf(hlog.LevelInfo, ctx, format, v...)
}

func (l *Logger) CtxNoticef(ctx context.Context, format string, v ...any) {
	l.CtxLogf(hlog.LevelWarn, ctx, format, v...)
}

func (l *Logger) CtxWarnf(ctx context.Context, format string, v ...any) {
	l.CtxLogf(hlog.LevelWarn, ctx, format, v...)
}

func (l *Logger) CtxErrorf(ctx context.Context, format string, v ...any) {
	l.CtxLogf(hlog.LevelError, ctx, format, v...)
}

func (l *Logger) CtxFatalf(ctx context.Context, format string, v ...any) {
	// l.CtxLogf(hlog.LevelFatal, ctx, format, v...)
	l.Logger.CtxFatalf(ctx, format, v...)
}
