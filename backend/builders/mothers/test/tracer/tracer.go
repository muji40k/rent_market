package tracer

import (
	"os"
	"rent_service/internal/misc/tracer/cleanstack"
	"rent_service/internal/misc/tracer/providers"
	"rent_service/internal/misc/tracer/wrap"

	"go.opentelemetry.io/contrib/instrumentation/host"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

const (
	TRACER_ENDPOINT string = "TRACER_ENDPOINT"
	METER_ENDPOINT  string = "METER_ENDPOINT"
)

func JaegerTracer(c *cleanstack.Cleaner) *tracesdk.TracerProvider {
	tendpoint := os.Getenv(TRACER_ENDPOINT)
	mendpoint := os.Getenv(METER_ENDPOINT)

	if "" == tendpoint {
		tendpoint = "localhost:4318"
	}

	if "" == mendpoint {
		mendpoint = "localhost:4318"
	}

	texp, err := providers.NewJaegerTraceExporter(tendpoint, c)

	if nil != err {
		panic(err)
	}

	mexp, err := providers.NewJaegerMetricExporter(mendpoint, c)

	if nil != err {
		panic(err)
	}

	tracer, err := providers.NewTraceProvider(texp, c)

	if nil != err {
		panic(err)
	}

	meter, err := providers.NewMeterProvider(mexp, c)

	if nil != err {
		panic(err)
	}

	host.Start(host.WithMeterProvider(meter))
	wrap.RegisterMetrics(meter.Meter("rent_service"))

	return tracer
}

