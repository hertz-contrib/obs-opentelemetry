# opentelemetry (这是一个社区驱动的项目)

[English](README.md) | 中文

适用于 [Hertz](https://github.com/cloudwego/hertz) 的 [Opentelemetry](https://opentelemetry.io/).

## 特性

#### Provider

- [x] 集成的默认 opentelemetry 程序，达到开箱即用
- [x] 支持设置环境变量

### Instrumentation

#### Tracing

- [x] 支持在 hertz 服务端和客户端之间启用 http 链路追踪
- [x] 支持通过设置 http header 以启动自动透明地传输对端服务

#### Metrics
- [x] 支持的 hertz http 指标有 [Rate, Errors, Duration]
- [x] 支持服务拓扑图指标 [服务拓扑图]
- [x] 支持 go runtime 指标

#### Logging

- [x] 在 logrus 的基础上适配了 hertz 日志工具
- [x] 实现了链路追踪自动关联日志的功能

## 通过环境变量来配置

- [Exporter](https://opentelemetry.io/docs/reference/specification/protocol/exporter/)
- [SDK](https://opentelemetry.io/docs/reference/specification/sdk-environment-variables/#general-sdk-configuration)

## 服务端使用示例

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

## 客户端使用示例

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

##  Tracing 和 Logging 进行关联

## 设置日志

```go
import (
    hertzlogrus "github.com/hertz-contrib/obs-opentelemetry/logging/logrus"
)

func init()  {
    hlog.SetLogger(hertzlogrus.NewLogger())
    hlog.SetLevel(hlog.LevelDebug)

}
```

#### 结合 context 使用 Logging

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

#### Logging 格式示例

```log
{"level":"debug","msg":"message received successfully: my request","span_id":"445ef16484a171b8","time":"2022-07-04T06:27:35+08:00","trace_flags":"01","trace_id":"e9e579b32c9d6b0598f8f33d65689e06"}
```

## 可以执行的示例

[Executable Example](https://github.com/cloudwego/hertz-examples/tree/main/opentelemetry)

## 现已支持的 Metrics

### HTTP Request Metrics

#### Hertz Server

下表列出了 HTTP 服务的指标

| 名称                          | Instrument Type | 单位        | 单位  | 描述                                                                  |
|-------------------------------|---------------------------------------------------|--------------|-------------------------------------------|------------------------------------------------------------------------------|
| `http.server.duration`        | Histogram                                         | milliseconds | `ms`                                      | 测量入站 HTTP 请求的耗时 |

#### Hertz Client

下表列出了 HTTP 客户端指标

| 名称                        | Instrument Type | 单位         | 单位 （UCUM） | 描述                                              |
|-----------------------------|---------------------------------------------------|--------------|-------------------------------------------|----------------------------------------------------------|
| `http.client.duration`      | Histogram                                         | milliseconds | `ms`                                      | 测量出站 HTTP 请求的耗时            |


### R.E.D
R.E.D (Rate, Errors, Duration) 定义了架构中的每个微服务测量的三个关键指标。OpenTelemetry 可以根据`http.server.duration`来计算R.E.D。

#### Rate

你的服务每秒钟所提供的请求数。

例如: QPS（Queries Per Second）每秒查询率

```
sum(rate(http_server_duration_count{}[5m])) by (service_name, http_method)
```

#### Errors

每秒失败的请求数。

例如：错误率

```
sum(rate(http_server_duration_count{status_code="Error"}[5m])) by (service_name, http_method) / sum(rate(http_server_duration_count{}[5m])) by (service_name, http_method)
```

#### Duration

每个请求所需时间的分布情况

例如：[P99 Latency](https://stackoverflow.com/questions/12808934/what-is-p99-latency) 

```
histogram_quantile(0.99, sum(rate(http_server_duration_bucket{}[5m])) by (le, service_name, http_method))
```

### 服务拓扑图

 `http.server.duration`将记录对等服务和当前服务维度。基于这个维度，我们可以汇总生成服务拓扑图
```
sum(rate(http_server_duration_count{}[5m])) by (service_name, peer_service)
```

### Runtime Metrics

| 名称                                   | 指标数据模型 | 单位       | 单位(UCUM) | 描述                                |
| -------------------------------------- | ------------ | ---------- | ---------- |-----------------------------------|
| `process.runtime.go.cgo.calls`         | Sum          | -          | -          | 当前进程调用的cgo数量                      |
| `process.runtime.go.gc.count`          | Sum          | -          | -          | 已完成的 gc 周期的数量                     |
| `process.runtime.go.gc.pause_ns`       | Histogram    | nanosecond | `ns`       | 在GC stop-the-world 中暂停的纳秒数量       |
| `process.runtime.go.gc.pause_total_ns` | Histogram    | nanosecond | `ns`       | 自程序启动以来，GC stop-the-world 的累计微秒计数 |
| `process.runtime.go.goroutines`        | Gauge        | -          | -          | 协程数量                              |
| `process.runtime.go.lookups`           | Sum          | -          | -          | 运行时执行的指针查询的数量                     |
| `process.runtime.go.mem.heap_alloc`    | Gauge        | bytes      | `bytes`    | 分配的堆对象的字节数                        |
| `process.runtime.go.mem.heap_idle`     | Gauge        | bytes      | `bytes`    | 空闲（未使用）的堆内存                       |
| `process.runtime.go.mem.heap_inuse`    | Gauge        | bytes      | `bytes`    | 已使用的堆内存                           |
| `process.runtime.go.mem.heap_objects`  | Gauge        | -          | -          | 已分配的堆对象数量                         |
| `process.runtime.go.mem.live_objects`  | Gauge        | -          | -          | 存活对象数量(Mallocs - Frees)           |
| `process.runtime.go.mem.heap_released` | Gauge        | bytes      | `bytes`    | 已交还给操作系统的堆内存                      |
| `process.runtime.go.mem.heap_sys`      | Gauge        | bytes      | `bytes`    | 从操作系统获得的堆内存                       |
| `runtime.uptime`                       | Sum          | ms         | `ms`       | 自应用程序被初始化以来的毫秒数                   |

##  兼容性

OpenTelemetry的 sdk 与1.x opentelemetry-go完全兼容，[详情查看](https://github.com/open-telemetry/opentelemetry-go#compatibility)


维护者: [CoderPoet](https://github.com/CoderPoet)

## 依赖

| **库/框架**                                         | 版本      | 记录   |
| --- |---------| --- |
| go.opentelemetry.io/otel | v1.9.0  | <br /> |
| go.opentelemetry.io/otel/trace | v1.9.0  | <br /> |
| go.opentelemetry.io/otel/metric | v0.31.0 | <br /> |
| go.opentelemetry.io/contrib/instrumentation/runtime | v0.30.0 |  |
| hertz | v0.4.1  |  |

