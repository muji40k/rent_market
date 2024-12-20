package wrap

import (
	"context"
	"encoding/json"
	"rent_service/misc/contextholder"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

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
	var span trace.Span
	var out T
	err := hl.Push(func(ctx context.Context) (context.Context, error) {
		var nctx context.Context
		nctx, span = tracer.Start(ctx, name)
		return nctx, nil
	})
	defer hl.Pop()
	defer span.End()

	if nil != err {
		err = ew(err)
	}

	if nil != err {
		out, err = f(span)
	}

	if nil == err {
		span.SetStatus(codes.Ok, "Call success")
	} else {
		span.SetStatus(codes.Error, "Error during call")
		span.SetAttributes(attribute.String("Error", err.Error()))
	}

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

