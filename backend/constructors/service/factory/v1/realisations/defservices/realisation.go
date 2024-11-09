package defservices

import (
	builder "rent_service/builders/service/factory/v1/defservices"
	constructor "rent_service/constructors/service/factory/v1"
	v1 "rent_service/internal/factory/services/v1"
	"rent_service/internal/logic/delivery"
	"rent_service/internal/logic/services/implementations/defservices/codegen/simple"
	"rent_service/internal/logic/services/implementations/defservices/paymentcheckers"
	"rent_service/internal/logic/services/implementations/defservices/photoregistry/implementations/defregistry/storages/local"
	rcontext "rent_service/internal/repository/context/v1"
	"strings"

	"github.com/google/uuid"
)

type Parser func() (Config, error)

type Config struct {
	CodegenLength uint
	Photo         struct {
		Main    string
		Temp    string
		BaseUrl string
	}
}

func New(
	parser Parser,
	deliveries delivery.ICreator,
	checkers map[uuid.UUID]paymentcheckers.IRegistrationChecker,
) constructor.Provider {
	return func() (string, constructor.Realisation) {
		return "default", newConstructor(parser, deliveries, checkers)
	}
}

func newConstructor(
	parser Parser,
	deliveries delivery.ICreator,
	checkers map[uuid.UUID]paymentcheckers.IRegistrationChecker,
) constructor.Realisation {
	return func(context *rcontext.Context) (v1.IFactory, error) {
		var sfactory v1.IFactory
		conf, err := parser()

		if nil == err {
			sfactory, err = builder.New().
				WithCodegen(simple.New(conf.CodegenLength)).
				WithDeliveryCreator(deliveries).
				WithPaymentCheckers(checkers).
				WithPhotoStorage(local.New(
					conf.Photo.Temp,
					conf.Photo.Main,
					func(path string) string {
						return strings.Replace(
							path,
							conf.Photo.Main,
							conf.Photo.BaseUrl,
							-1,
						)
					},
				)).
				WithRepositoryContext(context).
				Build()
		}

		return sfactory, err
	}
}

