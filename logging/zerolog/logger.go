// Copyright 2024 CloudWeGo Authors.
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
	*hertzzerolog.Logger
	config *config
}

// Ref to https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/logs/README.md#json-formats
const (
	traceIDKey    = "trace_id"
	spanIDKey     = "span_id"
	traceFlagsKey = "trace_flags"
)

type ExtraKey string

func NewLogger(opts ...Option) *Logger {
	cfg := defaultConfig()

	// apply options
	for _, opt := range opts {
		opt.apply(cfg)
	}
	logger := *cfg.logger
	zerologLogger := logger.Unwrap().
		Hook(cfg.defaultZerologHookFn())

	for i := range cfg.hooks {
		zerologLogger.Hook(cfg.hooks[i])
	}

	return &Logger{
		Logger: hertzzerolog.From(zerologLogger),
		config: cfg,
	}
}

func (l *Logger) CtxLogf(level hlog.Level, ctx context.Context, format string, kvs ...any) {
	var zlevel zerolog.Level

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

	span := trace.SpanFromContext(ctx)
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
