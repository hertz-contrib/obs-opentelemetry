# opentelemetry (这是一个社区驱动的项目)

[English](README.md) | 中文

适用于  [Hertz](https://github.com/cloudwego/hertz) 的 [Opentelemetry](https://opentelemetry.io/).

## 特性

#### 提供者

- [x] 集成的默认 opentelemetry 程序，达到开箱即用
- [x] 支持设置环境变量

### 遥测工具

#### 链路追踪

- [x] 支持在 hertz 服务端和客户端中的 http 链路追踪
- [x] 支持通过 http header 自动透明传输对等服务

#### Metrics
- [x] 支持Hertz http 指标 [Rate, Errors, Duration]
- [x] 支持服务拓扑图度量[服务拓扑图]。
- [x] 支持go runtime 度量

#### 日志

- [x] 在logrus的基础上扩展 Hertz 日志工具 
- [x] 实现跟踪自动关联日志

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

##  追踪相关日志

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

#### 结合 context 使用日志

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

#### 日志格式示例

```log
{"level":"debug","msg":"message received successfully: my request","span_id":"445ef16484a171b8","time":"2022-07-04T06:27:35+08:00","trace_flags":"01","trace_id":"e9e579b32c9d6b0598f8f33d65689e06"}
```

## 示例

[Executable Example](https://github.com/cloudwego/hertz-examples/tree/main/opentelemetry)

## 现已支持的 Mertrics

### RPC Metrics

#### Hertz Server

下面的表格为 RPC server metric 的配置项。

| 名称                   | 指标数据模型 | 单位     | 单位(UCUM) | 描述                  | 状态     | Streaming                                                    |
|------|------------|------|-------------------------------------------|-------------|--------|-----------|
| `http.server.duration` | Histogram    | 毫秒(ms) | `ms`       | 测量请求RPC的持续时间 | 推荐使用 | 并不适用， 虽然streaming RPC可能将这个指标记录为*批处理开始到批处理结束*，但在实际使用中很难解释。 |

#### Hertz Client

下面的表格为 RPC server metric 的配置项,这些适用于传统的RPC使用，不支持 streaming RPC

| Name | Instrument | Unit | Unit (UCUM) | Description | Status | Streaming |
|------|------------|------|-------------------------------------------|-------------|--------|-----------|
| `http.client.duration` | Histogram  | millseconds | `ms`        | 测量请求RPC的持续时间 | 推荐使用 | 并不适用， 虽然streaming RPC可能将这个指标记录为*批处理开始到批处理结束*，但在实际使用中很难解释。 |


### R.E.D
R.E.D (Rate, Errors, Duration) 方法定义了你应该为你架构中的每个微服务测量的三个关键指标。我们可以根据`http.server.duration`来计算R.E.D。

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

| **库/框架**                                         | 版本    | 记录   |
| --- | --- | --- |
| go.opentelemetry.io/otel | v1.7.0 | <br /> |
| go.opentelemetry.io/otel/trace | v1.7.0 | <br /> |
| go.opentelemetry.io/otel/metric | v0.30.0 | <br /> |
| go.opentelemetry.io/otel/semconv | v1.7.0 |  |
| go.opentelemetry.io/contrib/instrumentation/runtime | v0.30.0 |  |
| hertz | v0.1.0 |  |

