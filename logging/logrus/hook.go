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
	"github.com/cloudwego-contrib/cwgo-pkg/telemetry/instrumentation/otellogrus"

	"github.com/sirupsen/logrus"
)

// Ref to https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/logs/overview.md#json-formats

var _ logrus.Hook = (*TraceHook)(nil)

// TraceHookConfig trace hook config
type TraceHookConfig struct {
	recordStackTraceInSpan bool
	enableLevels           []logrus.Level
	errorSpanLevel         logrus.Level
}

// TraceHook trace hook
type TraceHook struct {
	hook otellogrus.TraceHook
}

// NewTraceHook create trace hook
func NewTraceHook(cfg *TraceHookConfig) *TraceHook {
	return &TraceHook{hook: *otellogrus.NewTraceHook(otellogrus.NewTraceHookConfig(
		cfg.recordStackTraceInSpan,
		cfg.enableLevels,
		cfg.errorSpanLevel))}
}

// Levels get levels
func (h *TraceHook) Levels() []logrus.Level {
	return h.hook.Levels()
}

// Fire logrus hook fire
func (h *TraceHook) Fire(entry *logrus.Entry) error {
	return h.hook.Fire(entry)
}
