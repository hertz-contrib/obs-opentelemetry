module github.com/hertz-contrib/obs-opentelemetry/logging/logrus

go 1.17

require (
	github.com/cloudwego/hertz v0.4.1
	github.com/hertz-contrib/logger/logrus v1.0.0
	github.com/sirupsen/logrus v1.9.0
	go.opentelemetry.io/otel v1.9.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.9.0
	go.opentelemetry.io/otel/sdk v1.9.0
	go.opentelemetry.io/otel/trace v1.9.0
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)
