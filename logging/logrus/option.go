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

// Option logrus hook option
type Option = otellogrus.Option

// WithLogger configures logger
func WithLogger(logger *logrus.Logger) Option {
	return otellogrus.WithLogger(logger)
}

// WithHook configures hook
func WithHook(hook logrus.Hook) Option {
	return otellogrus.WithHook(hook)
}

// WithTraceHookConfig configures trace hook config
func WithTraceHookConfig(hookConfig *TraceHookConfig) Option {
	return otellogrus.WithTraceHookConfig(otellogrus.NewTraceHookConfig(
		hookConfig.recordStackTraceInSpan,
		hookConfig.enableLevels,
		hookConfig.errorSpanLevel))
}

// WithTraceHookLevels configures hook levels
func WithTraceHookLevels(levels []logrus.Level) Option {
	return otellogrus.WithTraceHookLevels(levels)
}

// WithTraceHookErrorSpanLevel configures trace hook error span level
func WithTraceHookErrorSpanLevel(level logrus.Level) Option {
	return otellogrus.WithTraceHookErrorSpanLevel(level)
}

// WithRecordStackTraceInSpan configures whether record stack trace in span
func WithRecordStackTraceInSpan(recordStackTraceInSpan bool) Option {
	return otellogrus.WithRecordStackTraceInSpan(recordStackTraceInSpan)
}
