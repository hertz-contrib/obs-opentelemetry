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

package logrus

import (
	"errors"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Ref to https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/logs/overview.md#json-formats
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

var _ logrus.Hook = (*TraceHook)(nil)

type TraceHookConfig struct {
	recordStackTraceInSpan bool
	enableLevels           []logrus.Level
	errorSpanLevel         logrus.Level
}

type TraceHook struct {
	cfg *TraceHookConfig
}

func NewTraceHook(cfg *TraceHookConfig) *TraceHook {
	return &TraceHook{cfg: cfg}
}

func (h *TraceHook) Levels() []logrus.Level {
	return h.cfg.enableLevels
}

func (h *TraceHook) Fire(entry *logrus.Entry) error {
	if entry.Context == nil {
		return nil
	}

	span := trace.SpanFromContext(entry.Context)
	if !span.IsRecording() {
		return nil
	}

	// attach span context to log entry data fields
	entry.Data[traceIDKey] = span.SpanContext().TraceID()
	entry.Data[spanIDKey] = span.SpanContext().SpanID()
	entry.Data[traceFlagsKey] = span.SpanContext().TraceFlags()

	// attach log to span event attributes
	attrs := []attribute.KeyValue{
		logMessageKey.String(entry.Message),
		logSeverityTextKey.String(OtelSeverityText(entry.Level)),
	}
	span.AddEvent(logEventKey, trace.WithAttributes(attrs...))

	// set span status
	if entry.Level <= h.cfg.errorSpanLevel {
		span.SetStatus(codes.Error, entry.Message)
		span.RecordError(errors.New(entry.Message), trace.WithStackTrace(h.cfg.recordStackTraceInSpan))
	}

	return nil
}
