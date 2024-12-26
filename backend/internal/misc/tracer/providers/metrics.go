package providers

import (
	"rent_service/internal/misc/tracer/cleanstack"

	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func NewMeterProvider(exp metricsdk.Exporter, c *cleanstack.Cleaner) (*metricsdk.MeterProvider, error) {
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
		reader := metricsdk.NewPeriodicReader(exp)

		if nil != c {
			cleanstack.PushFE(c, reader.Shutdown)
		}

		return metricsdk.NewMeterProvider(
			metricsdk.WithReader(reader),
			metricsdk.WithResource(r),
		), nil
	}
}

