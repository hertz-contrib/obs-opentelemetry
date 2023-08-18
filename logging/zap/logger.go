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
	"context"
	"errors"
	"io"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ hlog.FullLogger = (*Logger)(nil)

// Ref to https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/logs/README.md#json-formats
const (
	traceIDKey    = "trace_id"
	spanIDKey     = "span_id"
	traceFlagsKey = "trace_flags"
	logEventKey   = "log"
)

var (
	logSeverityTextKey = attribute.Key("otel.log.severity.text")
	logMessageKey      = attribute.Key("otel.log.message")
)

type Logger struct {
	l      *zap.SugaredLogger
	config *config
}

func NewLogger(opts ...Option) *Logger {
	config := defaultConfig()

	// apply options
	for _, opt := range opts {
		opt.apply(config)
	}

	logger := zap.New(
		zapcore.NewCore(config.coreConfig.enc, config.coreConfig.ws, config.coreConfig.lvl),
		config.zapOpts...)

	return &Logger{
		l:      logger.Sugar(),
		config: config,
	}
}

func (l *Logger) Log(level hlog.Level, kvs ...interface{}) {
	switch level {
	case hlog.LevelTrace, hlog.LevelDebug:
		l.l.Debug(kvs...)
	case hlog.LevelInfo:
		l.l.Info(kvs...)
	case hlog.LevelNotice, hlog.LevelWarn:
		l.l.Warn(kvs...)
	case hlog.LevelError:
		l.l.Error(kvs...)
	case hlog.LevelFatal:
		l.l.Fatal(kvs...)
	default:
		l.l.Warn(kvs...)
	}
}

func (l *Logger) Logf(level hlog.Level, format string, kvs ...interface{}) {
	logger := l.l.With()
	switch level {
	case hlog.LevelTrace, hlog.LevelDebug:
		logger.Debugf(format, kvs...)
	case hlog.LevelInfo:
		logger.Infof(format, kvs...)
	case hlog.LevelNotice, hlog.LevelWarn:
		logger.Warnf(format, kvs...)
	case hlog.LevelError:
		logger.Errorf(format, kvs...)
	case hlog.LevelFatal:
		logger.Fatalf(format, kvs...)
	default:
		logger.Warnf(format, kvs...)
	}
}

func (l *Logger) CtxLogf(level hlog.Level, ctx context.Context, format string, kvs ...interface{}) {
	var zlevel zapcore.Level
	span := trace.SpanFromContext(ctx)

	sl := l.l.With(
		traceIDKey, span.SpanContext().TraceID(), spanIDKey, span.SpanContext().SpanID(), traceFlagsKey, span.SpanContext().TraceFlags())
	switch level {
	case hlog.LevelDebug, hlog.LevelTrace:
		zlevel = zap.DebugLevel
		sl.Debugf(format, kvs...)
	case hlog.LevelInfo:
		zlevel = zap.InfoLevel
		sl.Infof(format, kvs...)
	case hlog.LevelNotice, hlog.LevelWarn:
		zlevel = zap.WarnLevel
		sl.Warnf(format, kvs...)
	case hlog.LevelError:
		zlevel = zap.ErrorLevel
		sl.Errorf(format, kvs...)
	case hlog.LevelFatal:
		zlevel = zap.FatalLevel
		sl.Fatalf(format, kvs...)
	default:
		zlevel = zap.WarnLevel
		sl.Warnf(format, kvs...)
	}

	if !span.IsRecording() {
		l.Logf(level, format, kvs...)
		return
	}

	msg := getMessage(format, kvs)

	attrs := []attribute.KeyValue{
		logMessageKey.String(msg),
		logSeverityTextKey.String(OtelSeverityText(zlevel)),
	}
	span.AddEvent(logEventKey, trace.WithAttributes(attrs...))

	// set span status
	if zlevel >= l.config.traceConfig.errorSpanLevel {
		span.SetStatus(codes.Error, msg)
		span.RecordError(errors.New(msg), trace.WithStackTrace(l.config.traceConfig.recordStackTraceInSpan))
	}
}

func (l *Logger) Trace(v ...interface{}) {
	l.Log(hlog.LevelTrace, v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.Log(hlog.LevelDebug, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.Log(hlog.LevelInfo, v...)
}

func (l *Logger) Notice(v ...interface{}) {
	l.Log(hlog.LevelNotice, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.Log(hlog.LevelWarn, v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.Log(hlog.LevelError, v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.Log(hlog.LevelFatal, v...)
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	l.Logf(hlog.LevelTrace, format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Logf(hlog.LevelDebug, format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Logf(hlog.LevelInfo, format, v...)
}

func (l *Logger) Noticef(format string, v ...interface{}) {
	l.Logf(hlog.LevelInfo, format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Logf(hlog.LevelWarn, format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Logf(hlog.LevelError, format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Logf(hlog.LevelFatal, format, v...)
}

func (l *Logger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(hlog.LevelDebug, ctx, format, v...)
}

func (l *Logger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(hlog.LevelDebug, ctx, format, v...)
}

func (l *Logger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(hlog.LevelInfo, ctx, format, v...)
}

func (l *Logger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(hlog.LevelWarn, ctx, format, v...)
}

func (l *Logger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(hlog.LevelWarn, ctx, format, v...)
}

func (l *Logger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(hlog.LevelError, ctx, format, v...)
}

func (l *Logger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(hlog.LevelFatal, ctx, format, v...)
}

func (l *Logger) SetLevel(level hlog.Level) {
	var lvl zapcore.Level
	switch level {
	case hlog.LevelTrace, hlog.LevelDebug:
		lvl = zap.DebugLevel
	case hlog.LevelInfo:
		lvl = zap.InfoLevel
	case hlog.LevelWarn, hlog.LevelNotice:
		lvl = zap.WarnLevel
	case hlog.LevelError:
		lvl = zap.ErrorLevel
	case hlog.LevelFatal:
		lvl = zap.FatalLevel
	default:
		lvl = zap.WarnLevel
	}
	l.config.coreConfig.lvl.SetLevel(lvl)
}

func (l *Logger) SetOutput(writer io.Writer) {
	ws := zapcore.AddSync(writer)
	log := zap.New(
		zapcore.NewCore(l.config.coreConfig.enc, ws, l.config.coreConfig.lvl),
		l.config.zapOpts...,
	).Sugar()
	l.config.coreConfig.ws = ws
	l.l = log
}

func (l *Logger) Sync() {
	_ = l.l.Sync()
}
