package providers

import (
	"context"
	"rent_service/internal/misc/tracer/cleanstack"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

func NewJaegerTraceExporter(endpoint string, c *cleanstack.Cleaner) (*otlptrace.Exporter, error) {
	exp, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)

	if nil != c {
		cleanstack.PushFE(c, exp.Shutdown)
	}

	return exp, err
}

func NewJaegerMetricExporter(endpoint string, c *cleanstack.Cleaner) (*otlpmetrichttp.Exporter, error) {
	exp, err := otlpmetrichttp.New(context.Background(),
		otlpmetrichttp.WithEndpoint(endpoint),
		otlpmetrichttp.WithInsecure(),
	)

	if nil != c {
		cleanstack.PushFE(c, exp.Shutdown)
	}

	return exp, err
}

