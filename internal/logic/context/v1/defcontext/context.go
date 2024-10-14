package defcontext

import (
	"rent_service/internal/factory/services/v1/deffactory"
	"rent_service/internal/logic/context/v1"
)

func New(repositories interface{}) v1.Context {
	factory := deffactory.New()
	return v1.New(factory.ToFactories())
}

