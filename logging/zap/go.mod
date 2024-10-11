module github.com/hertz-contrib/obs-opentelemetry/logging/zap

go 1.21

require (
	github.com/cloudwego-contrib/cwgo-pkg/log/logging/zap v0.0.0
	github.com/cloudwego-contrib/cwgo-pkg/telemetry/instrumentation/otelzap v0.0.0
	github.com/cloudwego/hertz v0.9.2
	github.com/stretchr/testify v1.9.0
	go.opentelemetry.io/otel v1.28.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.28.0
	go.opentelemetry.io/otel/sdk v1.28.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/cloudwego-contrib/cwgo-pkg/telemetry/instrumentation/otelzap => D:\Projects\Go\cwgo-pkg\telemetry\instrumentation\otelzap

replace github.com/cloudwego-contrib/cwgo-pkg/log/logging/zap => D:\Projects\Go\cwgo-pkg\log\logging\zap


