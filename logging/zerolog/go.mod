module github.com/hertz-contrib/obs-opentelemetry/logging/zerolog

go 1.21

require (
	github.com/cloudwego-contrib/cwgo-pkg/log/logging/zerolog v0.0.0-00010101000000-000000000000
	github.com/cloudwego-contrib/cwgo-pkg/telemetry/instrumentation/otelzerolog v0.0.0-00010101000000-000000000000
	github.com/cloudwego/hertz v0.9.2
	github.com/rs/zerolog v1.31.0
	github.com/stretchr/testify v1.9.0
	go.opentelemetry.io/otel v1.28.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.28.0
	go.opentelemetry.io/otel/sdk v1.28.0
)

require (
	github.com/cloudwego/kitex v0.11.3 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/cloudwego-contrib/cwgo-pkg/telemetry/instrumentation/otelzerolog => github.com/smx-Morgan/cwgo-pkg/telemetry/instrumentation/otelzerolog v0.0.0-20241019002536-84cf43046703

replace github.com/cloudwego-contrib/cwgo-pkg/log/logging/zerolog => github.com/smx-Morgan/cwgo-pkg/log/logging/zerolog v0.0.0-20241019002536-84cf43046703
