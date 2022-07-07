module github.com/hertz-contrib/obs-opentelemetry/logging/logrus

go 1.17

require (
	github.com/cloudwego/hertz v0.1.0
	github.com/sirupsen/logrus v1.8.1
	go.opentelemetry.io/otel v1.4.1
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.4.1
	go.opentelemetry.io/otel/sdk v1.4.1
	go.opentelemetry.io/otel/trace v1.4.1
)

require (
	github.com/go-logr/logr v1.2.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	golang.org/x/sys v0.0.0-20220110181412-a018aaa089fe // indirect
)
