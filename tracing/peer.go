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

package tracing

import (
	"context"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

func injectPeerServiceToMetaInfo(ctx context.Context, attrs []attribute.KeyValue) map[string]string {
	md := metainfo.GetAllValues(ctx)
	if md == nil {
		md = make(map[string]string)
	}

	serviceName, serviceNamespace, deploymentEnv := getServiceFromResourceAttributes(attrs)

	if serviceName != "" {
		md[string(semconv.ServiceNameKey)] = serviceName
	}

	if serviceNamespace != "" {
		md[string(semconv.ServiceNamespaceKey)] = serviceNamespace
	}

	if deploymentEnv != "" {
		md[string(semconv.DeploymentEnvironmentKey)] = deploymentEnv
	}

	return md
}

func extractPeerServiceAttributesFromMetaInfo(md map[string]string) []attribute.KeyValue {
	var attrs []attribute.KeyValue

	for k, v := range md {
		switch k {
		case string(semconv.ServiceNameKey):
			attrs = append(attrs, semconv.PeerServiceKey.String(v))
		case string(semconv.ServiceNamespaceKey):
			attrs = append(attrs, PeerServiceNamespaceKey.String(v))
		case string(semconv.DeploymentEnvironmentKey):
			attrs = append(attrs, PeerDeploymentEnvironmentKey.String(v))
		}
	}

	return attrs
}
