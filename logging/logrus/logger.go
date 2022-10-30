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
	hertzlogrus "github.com/hertz-contrib/logger/logrus"
)

// Logger an alias to github.com/hertz-contrib/logger/logrus Logger
type Logger = hertzlogrus.Logger

// NewLogger create logger with otel hook
func NewLogger(opts ...Option) *Logger {
	cfg := defaultConfig()

	// apply options
	for _, opt := range opts {
		opt.apply(cfg)
	}

	// default trace hooks
	cfg.hooks = append(cfg.hooks, NewTraceHook(cfg.traceHookConfig))

	// attach hook
	for _, hook := range cfg.hooks {
		cfg.logger.AddHook(hook)
	}

	return hertzlogrus.NewLogger(
		hertzlogrus.WithLogger(cfg.logger),
	)
}
