package providers

import (
	"rent_service/internal/misc/tracer/cleanstack"

	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func NewTraceProvider(exp tracesdk.SpanExporter, c *cleanstack.Cleaner) (*tracesdk.TracerProvider, error) {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("rent_service"),
		),
	)

	if err != nil {
		return nil, err
	} else {
		prc := tracesdk.NewBatchSpanProcessor(exp)

		if nil != c {
			cleanstack.PushFE(c, prc.Shutdown)
		}

		return tracesdk.NewTracerProvider(
			tracesdk.WithSpanProcessor(prc),
			tracesdk.WithResource(r),
		), nil
	}
}

