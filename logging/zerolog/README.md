# Notes
- If you're using the global logger in `github.com/rs/zerolog/log.Logger` to do log. You should add the middleware to pass the OTEL extra information (key, value) into the `context.Context`.
- Example:
```go
import (
    "context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
	"github.com/hertz-contrib/obs-opentelemetry/logging/zerolog"
)

func OtelZerologMiddleware() app.HandlerFunc {
	return func(ctx context.Context, reqCtx *app.RequestContext) {
		ctx = log.Logger.WithContext(ctx)

		span := trace.SpanFromContext(ctx)

		ctx = context.WithValue(ctx, zerolog.ExtraKey(zerolog.TraceIDKey), span.SpanContext().TraceID())
		ctx = context.WithValue(ctx, zerolog.ExtraKey(zerolog.SpanIDKey), span.SpanContext().SpanID())
		ctx = context.WithValue(ctx, zerolog.ExtraKey(zerolog.TraceFlagsKey), span.SpanContext().TraceFlags())

		reqCtx.Next(ctx)
	}
```
