package v1

import (
	"rent_service/application"
	"rent_service/constructors"
	acv1 "rent_service/constructors/application/v1"
	rfcv1 "rent_service/constructors/repository/factory/v1"
	sfcv1 "rent_service/constructors/service/factory/v1"
	rfv1 "rent_service/internal/factory/repositories/v1"
	sfv1 "rent_service/internal/factory/services/v1"
	sv1 "rent_service/internal/logic/context/v1"
	rv1 "rent_service/internal/repository/context/v1"
)

type Clean func()

func Construct(
	app *acv1.Constructor,
	service *sfcv1.Constructor,
	repository *rfcv1.Constructor,
) (application.IApplication, Clean, error) {
	var err error
	var rfactory rfv1.IFactory
	var rcontext *rv1.Context
	var sfactory sfv1.IFactory
	var scontext *sv1.Context
	var result application.IApplication

	cleaner := constructors.NewCleaner()
	rfactory, err = repository.Construct(cleaner)

	if nil == err {
		rcontext = rv1.New(rfactory.ToFactories())
		sfactory, err = service.Construct(rcontext, cleaner)
	}

	if nil == err {
		scontext = sv1.New(sfactory.ToFactories())
		result, err = app.Construct(scontext, cleaner)
	}

	return result, cleaner.Clean, err
}

