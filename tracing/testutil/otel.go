// Copyright 2022 CloudWeGo Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testutil

import (
	"os"

	"github.com/cloudwego-contrib/cwgo-pkg/telemetry/meter/global"
	"github.com/cloudwego-contrib/cwgo-pkg/telemetry/semantic"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"

	cwmetric "github.com/cloudwego-contrib/cwgo-pkg/telemetry/meter/metric"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	otelmetric "go.opentelemetry.io/otel/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// OtelTestProvider get otel test provider
func OtelTestProvider() (*sdktrace.TracerProvider, otelmetric.MeterProvider, *prometheus.Registry) {
	// prometheus registry
	registry := prometheus.NewRegistry()

	// init tracer
	tracerProvider, err := initTracer()
	if err != nil {
		panic(err)
	}

	meterProvider, err := initMeterProvider(registry)
	if err != nil {
		panic(err)
	}
	meter := meterProvider.Meter(
		"github.com/cloudwego-contrib/telemetry-opentelemetry",
		otelmetric.WithInstrumentationVersion(semantic.SemVersion()),
	)
	serverRequestCountMeasure, err := meter.Int64Counter(
		semantic.BuildMetricName("http", "server", semantic.RequestCount),
		otelmetric.WithUnit("count"),
		otelmetric.WithDescription("measures Incoming request count total"),
	)
	HandleErr(err)

	serverLatencyMeasure, err := meter.Float64Histogram(
		semantic.BuildMetricName("http", "server", semantic.ServerLatency),
		otelmetric.WithUnit("ms"),
		otelmetric.WithDescription("measures th incoming end to end duration"),
	)
	HandleErr(err)

	measureServer := cwmetric.NewMeasure(
		cwmetric.WithCounter(semantic.HTTPCounter, cwmetric.NewOtelCounter(serverRequestCountMeasure)),
		cwmetric.WithRecorder(semantic.HTTPLatency, cwmetric.NewOtelRecorder(serverLatencyMeasure)),
	)
	global.SetTracerMeasure(measureServer)
	return tracerProvider, meterProvider, registry
}

// GatherAndCompare compare metrics with registry
func GatherAndCompare(registry *prometheus.Registry, expectedFilePath string, metricName ...string) error {
	file, err := os.Open(expectedFilePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	err = testutil.GatherAndCompare(registry, file, metricName...)
	if err != nil {
		return err
	}
	return nil
}

func initMeterProvider(registry *prometheus.Registry) (otelmetric.MeterProvider, error) {
	exporter, err := initMetricExporter(registry)
	if err != nil {
		return nil, err
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	return provider, nil
}

func initMetricExporter(registry *prometheus.Registry) (*otelprom.Exporter, error) {
	return otelprom.New(
		otelprom.WithRegisterer(registry),
	)
}

func initTracer() (*sdktrace.TracerProvider, error) {
	// Create stdout exporter to be able to retrieve
	// the collected spans.
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("test-server"),
			semconv.ServiceNamespaceKey.String("test-ns"),
			semconv.DeploymentEnvironmentKey.String("test-env"),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, err
}

func HandleErr(err error) {
	if err != nil {
		otel.Handle(err)
	}
}
