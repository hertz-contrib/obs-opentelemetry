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
	"testing"

	"github.com/sirupsen/logrus"
)

func TestOtelSeverityText(t *testing.T) {
	type args struct {
		lv logrus.Level
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "warn",
			args: args{
				lv: logrus.WarnLevel,
			},
			want: "WARN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OtelSeverityText(tt.args.lv); got != tt.want {
				t.Errorf("OtelSeverityText() = %v, want %v", got, tt.want)
			}
		})
	}
}
