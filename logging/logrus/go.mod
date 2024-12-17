module github.com/hertz-contrib/obs-opentelemetry/logging/logrus

go 1.19

require (
	github.com/cloudwego/hertz v0.9.2
	github.com/hertz-contrib/logger/logrus v1.0.0
	github.com/sirupsen/logrus v1.9.3
	go.opentelemetry.io/otel v1.19.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.19.0
	go.opentelemetry.io/otel/sdk v1.19.0
	go.opentelemetry.io/otel/trace v1.19.0
)

require (
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	go.opentelemetry.io/otel/metric v1.19.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
)
