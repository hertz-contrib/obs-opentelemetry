# Zerolog + OpenTelemetry + Hezrt

## [Document](https://www.cloudwego.io/docs/hertz/tutorials/observability/open-telemetry/)

## Example

1. See this [example](https://github.com/cloudwego/hertz-examples/tree/main/opentelemetry)
2. Small change in herzt server

```go
import (
	// ...

	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzZerolog "github.com/hertz-contrib/logger/zerolog"
	hertzOtelZerolog "github.com/hertz-contrib/obs-opentelemetry/logging/zerolog"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// ...
	
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(serviceName),
		// Support setting ExportEndpoint via environment variables: OTEL_EXPORTER_OTLP_ENDPOINT
		provider.WithExportEndpoint("localhost:4317"),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

	tracer, cfg := hertztracing.NewServerTracer()
	h := server.Default(tracer)
	h.Use(hertztracing.ServerMiddleware(cfg))

	logger := hertzZerolog.New(
		hertzZerolog.WithOutput(w),             // allows to specify output
		hertzZerolog.WithLevel(hlog.LevelInfo), // option with log level
		hertzZerolog.WithCaller(),              // option with caller
		hertzZerolog.WithTimestamp(),           // option with timestamp
		hertzZerolog.WithFormattedTimestamp(time.RFC3339),
	)

	log.Logger = logger.Unwrap() // log.Output(w).With().Caller().Logger()
	log.Logger = log.Level(zerolog.InfoLevel)

	otelLogger := hertzOtelZerolog.NewLogger(hertzOtelZerolog.WithLogger(logger))
	log.Logger = otelLogger.Unwrap()
	hlog.SetLogger(otelLogger)

	// ...
}
```

