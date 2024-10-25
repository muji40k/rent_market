package main

import (
	"rent_service/internal/factory/repositories/v1/psql"
	"rent_service/internal/factory/services/v1/deffactory"
	cv1 "rent_service/internal/logic/context/v1"
	delivery_dummy "rent_service/internal/logic/delivery/implementations/dummy"
	simple_codegen "rent_service/internal/logic/services/implementations/defservices/misc/codegen/simple"
	checker_dummy "rent_service/internal/logic/services/implementations/defservices/misc/paymentcheckers/dummy"
	"rent_service/internal/logic/services/implementations/defservices/misc/photoregistry/storages/local"
	"rent_service/internal/logic/services/implementations/defservices/payment"
	rv1 "rent_service/internal/repository/context/v1"
	"rent_service/internal/repository/implementation/sql/hashers/md5"
	simple_setter "rent_service/internal/repository/implementation/sql/technical/implementations/simple"
	"rent_service/server"
	"rent_service/server/authenticator/implementations/apikey"
	"rent_service/server/authenticator/implementations/apikey/repositories/sql"
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
	"strings"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	db1, err1 := sqlx.Connect(
		"pgx",
		"postgres://postgres:postgres@localhost:5432/rent_market",
	)

	db2, err2 := sqlx.Connect(
		"pgx",
		"postgres://postgres:postgres@localhost:5432/authentication",
	)

	defer func() {
		if nil != db1 {
			db1.Close()
		}
	}()
	defer func() {
		if nil != db2 {
			db2.Close()
		}
	}()

	if nil != err1 {
		panic(err1)
	}

	if nil != err2 {
		panic(err2)
	}

	rcontext := rv1.New(psql.New(
		db1,
		simple_setter.New("web_backend"),
		md5.Hash,
	).ToFactories())

	scontext := cv1.New(deffactory.New(
		&rcontext,
		simple_codegen.New(6),
		local.New(
			"../server/temp/",
			"../server/media/",
			func(path string) string {
				return strings.Replace(path, "../server/media/", "http://localhost/static/", -1)
			},
		),
		delivery_dummy.New(),
		map[uuid.UUID]payment.IRegistrationChecker{
			checker_dummy.Id: checker_dummy.New(),
		},
	).ToFactories())

	authenticator := apikey.New(&scontext, sql.New(db2))

	var s = server.New(
		server.WithPort(12345),
		server.WithCors(
			categories.CorsFiller,
			deliveries.CorsFiller,
			deliverycompanies.CorsFiller,
			instances.CorsFiller,
			payments.CorsFiller,
			paymethods.CorsFiller,
			periods.CorsFiller,
			photos.CorsFiller,
			pickuppoints.CorsFiller,
			products.CorsFiller,
			profiles.CorsFiller,
			prequests.CorsFiller,
			preturns.CorsFiller,
			provisions.CorsFiller,
			rrequests.CorsFiller,
			rreturns.CorsFiller,
			rents.CorsFiller,
			roles.CorsFiller,
			sessions.CorsFiller,
			storedinstances.CorsFiller,
			users.CorsFiller,
			headers.CorsFiller,
		),
		server.WithController(categories.New(&scontext)),
		server.WithController(deliveries.New(authenticator, &scontext)),
		server.WithController(deliverycompanies.New(&scontext, authenticator)),
		server.WithController(instances.New(
			authenticator,
			&scontext,
			&scontext,
			&scontext,
			&scontext,
		)),
		server.WithController(payments.New(authenticator, &scontext)),
		server.WithController(paymethods.New(&scontext)),
		server.WithController(periods.New(&scontext)),
		server.WithController(photos.New(&scontext, authenticator)),
		server.WithController(pickuppoints.New(&scontext, &scontext, &scontext)),
		server.WithController(products.New(&scontext, &scontext, &scontext)),
		server.WithController(profiles.New(
			authenticator,
			&scontext,
			&scontext,
			&scontext,
		)),
		server.WithController(prequests.New(
			authenticator,
			&scontext,
			&scontext,
		)),
		server.WithController(preturns.New(
			authenticator,
			&scontext,
			&scontext,
		)),
		server.WithController(provisions.New(authenticator, &scontext)),
		server.WithController(rrequests.New(
			authenticator,
			&scontext,
			&scontext,
		)),
		server.WithController(rreturns.New(
			authenticator,
			&scontext,
			&scontext,
		)),
		server.WithController(rents.New(authenticator, &scontext)),
		server.WithController(roles.New(authenticator, &scontext)),
		server.WithController(sessions.New(authenticator)),
		server.WithController(storedinstances.New(authenticator, &scontext)),
		server.WithController(users.New(authenticator, &scontext, &scontext)),
	)

	s.Run()
}

