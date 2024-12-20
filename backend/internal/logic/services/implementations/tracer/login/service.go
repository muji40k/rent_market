package login

import (
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/login"
	"rent_service/internal/logic/services/types/token"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/misc/contextholder"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type service struct {
	wrapped login.IService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped login.IService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) login.IService {
	return &service{wrapped, hl, tracer}
}

func (self *service) Register(
	email string,
	password string,
	name string,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.Register",
		func(span trace.Span) error {
			span.SetAttributes(
				attribute.String("Email", email),
			)
			return self.wrapped.Register(email, password, name)
		},
		cmnerrors.Internal,
	)
}

func (self *service) Login(email string, password string) (token.Token, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.Login",
		func(span trace.Span) (token.Token, error) {
			span.SetAttributes(
				attribute.String("Email", email),
			)
			return self.wrapped.Login(email, password)
		},
		cmnerrors.Internal,
	)
}

