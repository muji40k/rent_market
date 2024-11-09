package main

import (
	v1 "rent_service/constructors/v1"

	// Context
	scontext "rent_service/internal/logic/context/v1"

	// Static initialization
	delivery_composite "rent_service/internal/logic/delivery/implementations/composite"
	delivery_dummy "rent_service/internal/logic/delivery/implementations/dummy"
	"rent_service/internal/logic/services/implementations/defservices/paymentcheckers"
	checker_dummy "rent_service/internal/logic/services/implementations/defservices/paymentcheckers/dummy"
	"rent_service/internal/repository/implementation/sql/hashers/md5"
	"rent_service/internal/repository/implementation/sql/repositories/user"

	// Application construction
	application_parser "rent_service/constructors/parsers/env/application/v1"
	server_parser "rent_service/constructors/parsers/env/application/v1/server"
	auth_parser "rent_service/constructors/parsers/env/application/v1/server/authenticator"
	auth_apikey_parser "rent_service/constructors/parsers/env/application/v1/server/authenticator/apikey"
	auth_apikey_repo_parser "rent_service/constructors/parsers/env/application/v1/server/authenticator/apikey/repository"
	auth_apikey_repo_sql_parser "rent_service/constructors/parsers/env/application/v1/server/authenticator/apikey/repository/sql"

	application_constructor "rent_service/constructors/application/v1"
	server_realisation "rent_service/constructors/application/v1/realisations/server"
	auth_constructor "rent_service/constructors/application/v1/realisations/server/authenticators"
	auth_apikey_realisation "rent_service/constructors/application/v1/realisations/server/authenticators/realisations/apikey"
	auth_apikey_repo_constructor "rent_service/constructors/application/v1/realisations/server/authenticators/realisations/apikey/repositories"
	auth_apikey_repo_sql_realisation "rent_service/constructors/application/v1/realisations/server/authenticators/realisations/apikey/repositories/realisations/sql"

	// Service construction
	service_parser "rent_service/constructors/parsers/env/service/v1"
	service_default_parser "rent_service/constructors/parsers/env/service/v1/defservices"

	service_constructor "rent_service/constructors/service/factory/v1"
	service_default_realisation "rent_service/constructors/service/factory/v1/realisations/defservices"

	// Repository construction
	repository_parser "rent_service/constructors/parsers/env/repository/v1"
	repository_psql_parser "rent_service/constructors/parsers/env/repository/v1/psql"

	repository_constructor "rent_service/constructors/repository/factory/v1"
	repository_psql_realisation "rent_service/constructors/repository/factory/v1/realisations/psql"

	// Server stuff
	"rent_service/server"
	"rent_service/server/authenticator"
	"rent_service/server/controllers/categories"
	"rent_service/server/controllers/deliveries"
	"rent_service/server/controllers/deliverycompanies"
	"rent_service/server/controllers/instances"
	"rent_service/server/controllers/payments"
	"rent_service/server/controllers/paymethods"
	"rent_service/server/controllers/periods"
	"rent_service/server/controllers/photos"
	"rent_service/server/controllers/pickuppoints"
	"rent_service/server/controllers/products"
	"rent_service/server/controllers/profiles"
	prequests "rent_service/server/controllers/provision/requests"
	preturns "rent_service/server/controllers/provision/returns"
	"rent_service/server/controllers/provisions"
	rrequests "rent_service/server/controllers/rent/requests"
	rreturns "rent_service/server/controllers/rent/returns"
	"rent_service/server/controllers/rents"
	"rent_service/server/controllers/roles"
	"rent_service/server/controllers/sessions"
	"rent_service/server/controllers/storedinstances"
	"rent_service/server/controllers/users"
	"rent_service/server/headers"

	"github.com/google/uuid"
)

var hashers = map[string]user.Hasher{
	"md5": md5.Hash,
}

var delivery_creator = delivery_composite.New(
	delivery_composite.Pair(delivery_dummy.Id, delivery_dummy.New()),
)

var cdummy = checker_dummy.New()
var payment_checkers = map[uuid.UUID]paymentcheckers.IRegistrationChecker{
	cdummy.MethodId(): cdummy,
}

func main() {
	aconstructor := application_constructor.New(
		application_parser.Parser,
		server_realisation.New(
			server_parser.Parser,
			auth_constructor.New(
				auth_parser.Parser,
				auth_apikey_realisation.New(
					auth_apikey_parser.Parser,
					auth_apikey_repo_constructor.New(
						auth_apikey_repo_parser.Parser,
						auth_apikey_repo_sql_realisation.New(
							auth_apikey_repo_sql_parser.Parser,
						),
					),
				),
			),
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return categories.CorsFiller, categories.New(c)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return deliveries.CorsFiller, deliveries.New(a, c)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return deliverycompanies.CorsFiller, deliverycompanies.New(c, a)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return instances.CorsFiller, instances.New(a, c, c, c, c)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return payments.CorsFiller, payments.New(a, c)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return paymethods.CorsFiller, paymethods.New(c)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return periods.CorsFiller, periods.New(c)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return photos.CorsFiller, photos.New(c, a)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return pickuppoints.CorsFiller, pickuppoints.New(c, c, c)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return products.CorsFiller, products.New(c, c, c)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return profiles.CorsFiller, profiles.New(a, c, c, c)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return prequests.CorsFiller, prequests.New(a, c, c)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return preturns.CorsFiller, preturns.New(a, c, c)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return provisions.CorsFiller, provisions.New(a, c)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return rrequests.CorsFiller, rrequests.New(a, c, c)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return rreturns.CorsFiller, rreturns.New(a, c, c)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return rents.CorsFiller, rents.New(a, c)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return roles.CorsFiller, roles.New(a, c)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return sessions.CorsFiller, sessions.New(a)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return storedinstances.CorsFiller, storedinstances.New(a, c)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return users.CorsFiller, users.New(a, c, c)
			},
			func(a authenticator.IAuthenticator, c *scontext.Context) (server.CorsFiller, server.IController) {
				return headers.CorsFiller, nil
			},
		),
	)

	sconstructor := service_constructor.New(
		service_parser.Parser,
		service_default_realisation.New(
			service_default_parser.Parser,
			&delivery_creator,
			payment_checkers,
		),
	)

	rconstructor := repository_constructor.New(
		repository_parser.Parser,
		repository_psql_realisation.New(
			repository_psql_parser.Parser,
			hashers,
		),
	)

	app, cleaner, err := v1.Construct(aconstructor, sconstructor, rconstructor)

	defer cleaner()

	if nil == err {
		app.Run()
	} else {
		panic(err)
	}
}

