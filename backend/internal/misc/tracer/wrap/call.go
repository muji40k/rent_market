package wrap

import (
	"context"
	"encoding/json"
	"rent_service/misc/contextholder"
	"syscall"

	"runtime"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var (
	cpuTimeMetric      metric.Float64Counter
	cpuTimeTotalMetric metric.Float64Counter
	memoryMetric       metric.Float64Counter
)

func add(m metric.Float64Counter, value float64) {
	if nil != m {
		m.Add(context.Background(), value)
	}
}

func RegisterMetrics(meter metric.Meter) error {
	var err error

	cpuTimeMetric, err = meter.Float64Counter("tracer.cpu.time",
		metric.WithDescription("Cpu time spent on tracing"),
		// metric.WithUnit("{call}"),
	)

	if nil == err {
		cpuTimeTotalMetric, err = meter.Float64Counter("tracer.cpu.total",
			metric.WithDescription("Total Cpu time"),
			// metric.WithUnit("{call}"),
		)
	}

	if nil == err {
		memoryMetric, err = meter.Float64Counter("tracer.memory.usage",
			metric.WithDescription("Memory spent on tracing"),
			// metric.WithUnit("{call}"),
		)
	}

	return err
}

func SpanCall[E error](
	hl *contextholder.Holder,
	tracer trace.Tracer,
	name string,
	f func(span trace.Span) error,
	ew func(error) E,
) error {
	_, err := SpanCallValue(hl, tracer, name, func(span trace.Span) (error, error) {
		return nil, f(span)
	}, ew)

	return err
}

func SpanCallValue[T any, E error](
	hl *contextholder.Holder,
	tracer trace.Tracer,
	name string,
	f func(span trace.Span) (T, error),
	ew func(error) E,
) (T, error) {
	timer := timerWrap()
	memer := memWrap()
	timer.startGlobal()
	memer.startGlobal()

	var span trace.Span
	var out T
	serr := hl.Push(func(ctx context.Context) (context.Context, error) {
		var nctx context.Context
		nctx, span = tracer.Start(ctx, name)
		return nctx, nil
	})
	defer func() {
		if nil != span {
			span.End()
		}

		if nil == serr {
			hl.Pop()
		}
	}()

	var err error

	if nil != serr {
		err = ew(serr)
	}

	if nil == err {
		timer.startLocal()
		memer.startLocal()

		span.AddEvent("Actual start")
		out, err = f(span)
		span.AddEvent("Actual end")

		timer.stopLocal()
		memer.stopLocal()
	}

	if nil == err {
		span.SetStatus(codes.Ok, "Call success")
	} else if nil != serr {
		span.SetStatus(codes.Error, "Error with context")
		span.SetAttributes(attribute.String("ContextError", serr.Error()))
	} else {
		span.SetStatus(codes.Error, "Error during call")
		span.SetAttributes(attribute.String("Error", err.Error()))
	}

	timer.stopGlobal()
	memer.stopGlobal()
	add(cpuTimeMetric, timer.Global-timer.Local)
	add(cpuTimeTotalMetric, timer.Global)
	add(memoryMetric, memer.Global-memer.Local)

	return out, err
}

func AttributeJSON(name string, value interface{}) attribute.KeyValue {
	parsed, err := json.Marshal(value)

	if nil == err {
		return attribute.String(name, string(parsed))
	} else {
		return attribute.String("ParseError", err.Error())
	}
}

type mWrap[T any] struct {
	take   func() T
	sub    func(*T, *T) T
	gstart T
	lstart T
	Local  T
	Global T
}

func newW[T any](take func() T, sub func(*T, *T) T) mWrap[T] {
	return mWrap[T]{
		take: take,
		sub:  sub,
	}
}

func timerWrap() mWrap[float64] {
	return newW(func() float64 {
		var usage syscall.Rusage
		err := syscall.Getrusage(syscall.RUSAGE_SELF, &usage)

		if nil != err {
			return 0
		} else {
			return float64(usage.Utime.Nano()) + float64(usage.Stime.Nano())
		}
	}, func(a, b *float64) float64 {
		return *a - *b
	})
}

func memWrap() mWrap[float64] {
	return newW(func() float64 {
		var stats runtime.MemStats
		runtime.ReadMemStats(&stats)
		return float64(stats.TotalAlloc)
	}, func(a, b *float64) float64 {
		return *a - *b
	})
}

func (self *mWrap[T]) startGlobal() {
	self.gstart = self.take()
}

func (self *mWrap[T]) startLocal() {
	self.lstart = self.take()
}

func (self *mWrap[T]) stopLocal() {
	r := self.take()
	self.Local = self.sub(&r, &self.lstart)
}

func (self *mWrap[T]) stopGlobal() {
	r := self.take()
	self.Global = self.sub(&r, &self.gstart)
}

