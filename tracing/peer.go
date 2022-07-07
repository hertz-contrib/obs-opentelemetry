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
	"strings"

	"github.com/cloudwego/hertz/pkg/protocol"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

func injectPeerServiceToMetadata(_ context.Context, attrs []attribute.KeyValue) map[string]string {
	serviceName, serviceNamespace, deploymentEnv := getServiceFromResourceAttributes(attrs)

	md := make(map[string]string, 3)

	if serviceName != "" {
		md[semconvAttributeKeyToHTTPHeader(string(semconv.ServiceNameKey))] = serviceName
	}

	if serviceNamespace != "" {
		md[semconvAttributeKeyToHTTPHeader(string(semconv.ServiceNamespaceKey))] = serviceNamespace
	}

	if deploymentEnv != "" {
		md[semconvAttributeKeyToHTTPHeader(string(semconv.DeploymentEnvironmentKey))] = deploymentEnv
	}

	return md
}

func extractPeerServiceAttributesFromMetadata(headers *protocol.RequestHeader) []attribute.KeyValue {
	var attrs []attribute.KeyValue

	serviceName, serviceNamespace, deploymentEnv := headers.Get(semconvAttributeKeyToHTTPHeader(string(semconv.ServiceNameKey))),
		headers.Get(semconvAttributeKeyToHTTPHeader(string(semconv.ServiceNamespaceKey))),
		headers.Get(semconvAttributeKeyToHTTPHeader(string(semconv.DeploymentEnvironmentKey)))

	if serviceName != "" {
		attrs = append(attrs, semconv.PeerServiceKey.String(serviceName))
	}

	if serviceNamespace != "" {
		attrs = append(attrs, PeerServiceNamespaceKey.String(serviceNamespace))
	}

	if deploymentEnv != "" {
		attrs = append(attrs, PeerDeploymentEnvironmentKey.String(deploymentEnv))
	}

	return attrs
}

func semconvAttributeKeyToHTTPHeader(key string) string {
	return strings.ReplaceAll(key, ".", "-")
}
