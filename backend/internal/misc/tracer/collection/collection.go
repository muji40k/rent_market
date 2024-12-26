package collection

import (
	"context"
	"fmt"
	"reflect"
	"rent_service/internal/misc/types/collection"
	"rent_service/misc/contextholder"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type traceCollection[T any] struct {
	tracer     trace.Tracer
	hl         *contextholder.Holder
	collection collection.Collection[T]
}

type traceIterator[T any] struct {
	tracer trace.Tracer
	hl     *contextholder.Holder
	iter   collection.Iterator[T]
	i      uint
	id     uuid.UUID
}

func TraceCollection[T any](
	hl *contextholder.Holder,
	tracer trace.Tracer,
	collection collection.Collection[T],
) collection.Collection[T] {
	if nil == collection {
		return nil
	}

	return &traceCollection[T]{tracer, hl, collection}
}

func (self *traceCollection[T]) Iter() collection.Iterator[T] {
	return TraceIterator(self.hl, self.tracer, self.collection.Iter())
}

func TraceIterator[T any](
	hl *contextholder.Holder,
	tracer trace.Tracer,
	iter collection.Iterator[T],
) collection.Iterator[T] {
	if nil == iter {
		return nil
	}

	return &traceIterator[T]{tracer, hl, iter, 0, uuid.Must(uuid.NewRandom())}
}

func (self *traceIterator[T]) Next() (T, bool) {
	var span trace.Span
	err := self.hl.Push(func(ctx context.Context) (context.Context, error) {
		var nctx context.Context
		nctx, span = self.tracer.Start(
			ctx,
			fmt.Sprintf("Iterator [%v] Next: %v", self.id, self.i),
		)
		return nctx, nil
	})
	defer func() {
		if nil != span {
			span.End()
		}

		if nil == err {
			self.hl.Pop()
		}
	}()
	self.i++

	if nil != span {
		span.SetAttributes(attribute.String(
			"Type", fmt.Sprintf("%T", *(new(T))),
		))
		span.AddEvent("Actual start")
	}
	value, next := self.iter.Next()
	span.AddEvent("Actual end")
	span.SetAttributes(attribute.Bool("Next", next))

	if nil == err {
		if !next || !reflect.ValueOf(value).IsZero() {
			span.SetStatus(codes.Ok, "Wrap success")
		} else {
			span.SetStatus(codes.Error, "Empty value found, can be error")
		}
	} else {
		span.SetStatus(codes.Error, "Wrap error")
		span.SetAttributes(attribute.String("Error", err.Error()))
	}

	return value, next
}

func (self *traceIterator[T]) Skip() bool {
	var span trace.Span
	err := self.hl.Push(func(ctx context.Context) (context.Context, error) {
		var nctx context.Context
		nctx, span = self.tracer.Start(
			ctx,
			fmt.Sprintf("Iterator [%v] Skip: %v", self.id, self.i),
		)
		return nctx, nil
	})
	defer func() {
		if nil == err {
			self.hl.Pop()
		}

		if nil != span {
			span.End()
		}
	}()
	self.i++

	if nil != span {
		span.SetAttributes(attribute.String(
			"Type", fmt.Sprintf("%T", *(new(T))),
		))
		span.AddEvent("Actual start")
	}
	next := self.iter.Skip()
	if nil != span {
		span.AddEvent("Actual end")
		span.SetAttributes(attribute.Bool("Next", next))

		if nil == err {
			span.SetStatus(codes.Ok, "Wrap success")
		} else {
			span.SetStatus(codes.Error, "Wrap error")
			span.SetAttributes(attribute.String("Error", err.Error()))
		}
	}

	return next
}

