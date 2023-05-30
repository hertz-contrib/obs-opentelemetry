module github.com/hertz-contrib/obs-opentelemetry/logging/logrus

go 1.19

require (
	github.com/cloudwego/hertz v0.6.4
	github.com/hertz-contrib/logger/logrus v1.0.0
	github.com/sirupsen/logrus v1.9.0
	go.opentelemetry.io/otel v1.16.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.16.0
	go.opentelemetry.io/otel/sdk v1.16.0
	go.opentelemetry.io/otel/trace v1.16.0
)

require (
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel/metric v1.16.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
)
