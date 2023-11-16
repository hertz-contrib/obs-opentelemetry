package tracing_test

import (
	"context"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/obs-opentelemetry/testutil"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestMetricsExample(t *testing.T) {
	// test util
	tracerProvider, meterProvider, registry := testutil.OtelTestProvider()
	defer func(tracerProvider *sdktrace.TracerProvider, ctx context.Context) {
		_ = tracerProvider.Shutdown(ctx)
	}(tracerProvider, context.Background())
	otel.SetMeterProvider(meterProvider)

	// server example
	tracer, cfg := hertztracing.NewServerTracer()
	h := server.Default(tracer, server.WithHostPorts(":39888"))
	h.Use(hertztracing.ServerMiddleware(cfg))
	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		hlog.CtxDebugf(c, "message received successfully")
		ctx.JSON(consts.StatusOK, "pong")
	})
	go h.Spin()

	<-time.After(time.Millisecond * 500)

	// client example
	c, _ := client.NewClient()
	c.Use(hertztracing.ClientMiddleware())
	_, body, err := c.Get(context.Background(), nil, "http://localhost:39888/ping?foo=bar")
	require.NoError(t, err)
	assert.NotNil(t, body)

	// diff metrics
	assert.NoError(t, testutil.GatherAndCompare(
		registry, "testdata/hertz_request_metrics.txt",
		"http_server_request_count_total", "http_client_request_count_total"),
	)
}
