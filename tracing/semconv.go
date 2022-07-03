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

import "go.opentelemetry.io/otel/attribute"

// Attribute keys that can be added to a span.
const (
	ReadBytesKey  = attribute.Key("http.read_bytes")  // if anything was read from the request body, the total number of bytes read
	ReadErrorKey  = attribute.Key("http.read_error")  // If an error occurred while reading a request, the string of the error (io.EOF is not recorded)
	WroteBytesKey = attribute.Key("http.wrote_bytes") // if anything was written to the response writer, the total number of bytes written
	WriteErrorKey = attribute.Key("http.write_error") // if an error occurred while writing a reply, the string of the error (io.EOF is not recorded)
)

const (
	// RequestProtocolKey protocol of the request.
	//
	// Type: string
	// Required: Always
	// Examples:
	// http: 'http'
	// rpc: 'grpc', 'java_rmi', 'wcf', 'kitex'
	// db: mysql, postgresql
	// mq: 'rabbitmq', 'activemq', 'AmazonSQS'
	RequestProtocolKey = attribute.Key("request.protocol")
)

const (
	// PeerServiceNamespaceKey peer.service.namespace
	PeerServiceNamespaceKey = attribute.Key("peer.service.namespace")
	// PeerDeploymentEnvironmentKey peer.deployment.environment
	PeerDeploymentEnvironmentKey = attribute.Key("peer.deployment.environment")
)

const (
	StatusKey = attribute.Key("status.code")
)
