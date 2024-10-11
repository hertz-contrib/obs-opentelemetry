module github.com/hertz-contrib/obs-opentelemetry/logging/logrus

go 1.21

require (
	github.com/cloudwego-contrib/cwgo-pkg/telemetry/instrumentation/otellogrus v0.0.0
	github.com/cloudwego/hertz v0.9.2
	github.com/sirupsen/logrus v1.9.3
	go.opentelemetry.io/otel v1.28.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.28.0
	go.opentelemetry.io/otel/sdk v1.28.0
)

require (
	github.com/cloudwego-contrib/cwgo-pkg/log/logging/logrus v0.0.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
)

replace github.com/cloudwego-contrib/cwgo-pkg/telemetry/instrumentation/otellogrus => D:\Projects\Go\cwgo-pkg\telemetry/instrumentation/otellogrus

replace github.com/cloudwego-contrib/cwgo-pkg/log/logging/logrus => D:\Projects\Go\cwgo-pkg\log\logging\logrus
