# opentelemetry (This is a community driven project)

English | [中文](README_CN.md)

[Opentelemetry](https://opentelemetry.io/) for [Hertz](https://github.com/cloudwego/hertz)

## Feature
#### Provider
- [x] Out-of-the-box default opentelemetry provider
- [x] Support setting via environment variables

### Instrumentation

#### Tracing
- [x] Support server and client hertz http tracing
- [x] Support automatic transparent transmission of peer service through http headers

#### Metrics
- [x] Support hertz http metrics [R.E.D]
- [x] Support service topology map metrics [Service Topology Map]
- [x] Support go runtime metrics

#### Logging
- [x] Extend hertz logger based on logrus
- [x] Implement tracing auto associated logs

## Configuration via environment variables
- [Exporter](https://opentelemetry.io/docs/reference/specification/protocol/exporter/)
- [SDK](https://opentelemetry.io/docs/reference/specification/sdk-environment-variables/#general-sdk-configuration)

## Server usage
```go
import (
    ...
    "github.com/hertz-contrib/obs-opentelemetry/provider"
    "github.com/hertz-contrib/obs-opentelemetry/tracing"
)


func main()  {
    serviceName := "echo"
	
    p := provider.NewOpenTelemetryProvider(
        provider.WithServiceName(serviceName),
        provider.WithExportEndpoint("localhost:4317"),
        provider.WithInsecure(),
    )
    defer p.Shutdown(context.Background())

    tracer, cfg := hertztracing.NewServerTracer()
    h := server.Default(tracer)
    h.Use(hertztracing.ServerMiddleware(cfg))
    
    ...
	
    h.Spin()
}

```

## Client usage
```go
import (
    ...
    "github.com/hertz-contrib/obs-opentelemetry/provider"
    "github.com/hertz-contrib/obs-opentelemetry/tracing"
)

func main(){
    serviceName := "echo-client"
	
    p := provider.NewOpenTelemetryProvider(
        provider.WithServiceName(serviceName),
        provider.WithExportEndpoint("localhost:4317"),
        provider.WithInsecure(),
    )
    defer p.Shutdown(context.Background())

    c, _ := client.NewClient()
    c.Use(hertztracing.ClientMiddleware())

    ...   
	
}

```

## Tracing associated Logs

#### set logger impl
```go
import (
    hertzlogrus "github.com/hertz-contrib/obs-opentelemetry/logging/logrus"
)

func init()  {
    hlog.SetLogger(hertzlogrus.NewLogger())
    hlog.SetLevel(hlog.LevelDebug)

}
```

#### log with context

```go
h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
    req := &api.Request{Message: "my request"}
    resp, err := client.Echo(c, req)
    if err != nil {
        hlog.Errorf(err.Error())
    }
    hlog.CtxDebugf(c, "message received successfully: %s", req.Message)
    ctx.JSON(consts.StatusOK, resp)
})
```

#### view log

```log
{"level":"debug","msg":"message received successfully: my request","span_id":"445ef16484a171b8","time":"2022-07-04T06:27:35+08:00","trace_flags":"01","trace_id":"e9e579b32c9d6b0598f8f33d65689e06"}
```


## Example

[Executable Example](https://github.com/cloudwego/hertz-examples/tree/main/opentelemetry)

## Supported Metrics

### HTTP Request Metrics

#### Hertz Server

Below is a table of HTTP server metric instruments.

| Name                          | Instrument Type | Unit        | Unit  | Description                                                                  |
|-------------------------------|---------------------------------------------------|--------------|-------------------------------------------|------------------------------------------------------------------------------|
| `http.server.duration`        | Histogram                                         | milliseconds | `ms`                                      | measures the duration inbound HTTP requests |


#### Hertz Client

Below is a table of HTTP client metric instruments.

| Name                        | Instrument Type ([*](README.md#instrument-types)) | Unit         | Unit ([UCUM](README.md#instrument-units)) | Description                                              |
|-----------------------------|---------------------------------------------------|--------------|-------------------------------------------|----------------------------------------------------------|
| `http.client.duration`      | Histogram                                         | milliseconds | `ms`                                      | measures the duration outbound HTTP requests             |


### R.E.D
The RED Method defines the three key metrics you should measure for every microservice in your architecture. We can calculate RED based on `http.server.duration`.

#### Rate
the number of requests, per second, you services are serving.

eg: QPS
```
sum(rate(http_server_duration_count{}[5m])) by (service_name, http_method)
```

#### Errors
the number of failed requests per second.

eg: Error ratio
```
sum(rate(http_server_duration_count{status_code="Error"}[5m])) by (service_name, http_method) / sum(rate(http_server_duration_count{}[5m])) by (service_name, http_method)
```

#### Duration
distributions of the amount of time each request takes

eg: P99 Latency
```
histogram_quantile(0.99, sum(rate(http_server_duration_bucket{}[5m])) by (le, service_name, http_method))
```

### Service Topology Map
The `http.server.duration` will record the peer service and the current service dimension. Based on this dimension, we can aggregate the service topology map
```
sum(rate(http_server_duration_count{}[5m])) by (service_name, peer_service)
```

### Runtime Metrics
| Name | Instrument | Unit | Unit (UCUM)) | Description |
|------|------------|------|-------------------------------------------|-------------|
| `process.runtime.go.cgo.calls` | Sum | - | - | Number of cgo calls made by the current process. |
| `process.runtime.go.gc.count` | Sum | - | - | Number of completed garbage collection cycles. |
| `process.runtime.go.gc.pause_ns` | Histogram | nanosecond | `ns` | Amount of nanoseconds in GC stop-the-world pauses. |
| `process.runtime.go.gc.pause_total_ns` | Histogram | nanosecond | `ns` | Cumulative nanoseconds in GC stop-the-world pauses since the program started. |
| `process.runtime.go.goroutines` | Gauge | - | - | measures duration of outbound HTTP Request. | 
| `process.runtime.go.lookups` | Sum | - | - | Number of pointer lookups performed by the runtime. |
| `process.runtime.go.mem.heap_alloc` | Gauge | bytes | `bytes` | Bytes of allocated heap objects. |
| `process.runtime.go.mem.heap_idle` | Gauge | bytes | `bytes` | Bytes in idle (unused) spans. |
| `process.runtime.go.mem.heap_inuse` | Gauge | bytes | `bytes` | Bytes in in-use spans. |
| `process.runtime.go.mem.heap_objects` | Gauge | - | - | Number of allocated heap objects. |
| `process.runtime.go.mem.live_objects` | Gauge | - | - | Number of live objects is the number of cumulative Mallocs - Frees. |
| `process.runtime.go.mem.heap_released` | Gauge | bytes | `bytes` | Bytes of idle spans whose physical memory has been returned to the OS. |
| `process.runtime.go.mem.heap_sys` | Gauge | bytes | `bytes` | Bytes of idle spans whose physical memory has been returned to the OS. |
| `runtime.uptime` | Sum | ms | `ms` | Milliseconds since application was initialized. |


## Compatibility
The sdk of OpenTelemetry is fully compatible with 1.X opentelemetry-go. [see](https://github.com/open-telemetry/opentelemetry-go#compatibility)


maintained by: [CoderPoet](https://github.com/CoderPoet)


## Dependencies
| **Library/Framework** | **Versions** | **Notes** |
| --- |---------| --- |
| go.opentelemetry.io/otel | v1.9.0  | <br /> |
| go.opentelemetry.io/otel/trace | v1.9.0  | <br /> |
| go.opentelemetry.io/otel/metric | v0.31.0 | <br /> |
| go.opentelemetry.io/contrib/instrumentation/runtime | v0.30.0 |  |
| hertz | v0.4.1  |  |


