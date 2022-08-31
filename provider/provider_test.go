package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	semconv140 "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func Test_newResource(t *testing.T) {
	type args struct {
		cfg *config
	}
	tests := []struct {
		name              string
		args              args
		wantResources     []attribute.KeyValue
		unwantedResources []attribute.KeyValue
	}{
		{
			name: "with conflict schema version",
			args: args{
				cfg: &config{
					resourceAttributes: []attribute.KeyValue{
						semconv140.ServiceNameKey.String("test-semconv-resource"),
					},
				},
			},
			wantResources: []attribute.KeyValue{
				semconv.ServiceNameKey.String("test-semconv-resource"),
			},
			unwantedResources: []attribute.KeyValue{
				semconv.ServiceNameKey.String("unknown_service:___Test_newResource_in_github_com_hertz_contrib_obs_opentelemetry_provider.test"),
			},
		},
		{
			name: "resource override",
			args: args{
				cfg: &config{
					resource: resource.Default(),
					resourceAttributes: []attribute.KeyValue{
						semconv.ServiceNameKey.String("test-resource"),
					},
				},
			},
			wantResources: nil,
			unwantedResources: []attribute.KeyValue{
				semconv.ServiceNameKey.String("test-resource"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newResource(tt.args.cfg)
			for _, res := range tt.wantResources {
				assert.Contains(t, got.Attributes(), res)
			}
			for _, unwantedResource := range tt.unwantedResources {
				assert.NotContains(t, got.Attributes(), unwantedResource)
			}
		})
	}
}
