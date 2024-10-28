package server

import (
	"rent_service/application"
	builder "rent_service/builders/application/web"
	"rent_service/constructors"
	constructor "rent_service/constructors/application/v1"
	aconstructor "rent_service/constructors/application/v1/realisations/server/authenticators"
	scontext "rent_service/internal/logic/context/v1"
	"rent_service/server"
	"rent_service/server/authenticator"
)

type Config struct {
	Host       string
	Port       uint
	SwaggerURL string
}

type Parser func() (Config, error)
type Extender func(
	auth authenticator.IAuthenticator,
	context *scontext.Context,
) (server.CorsFiller, server.IController)

func New(
	parser Parser,
	authConstructor *aconstructor.Constructor,
	extenders ...Extender,
) constructor.Provider {
	return func() (string, constructor.Realisation) {
		return "web", newConstructor(parser, authConstructor, extenders...)
	}
}

type cleanerWrap struct {
	app     application.IApplication
	cleaner *constructors.Cleaner
}

func (self *cleanerWrap) Run() {
	self.app.Run()
}

func (self *cleanerWrap) Clear() {
	self.app.Clear()
	self.cleaner.Clean()
}

func newConstructor(
	parser Parser,
	authConstructor *aconstructor.Constructor,
	extenders ...Extender,
) constructor.Realisation {
	return func(context *scontext.Context) (application.IApplication, error) {
		var wrap = &cleanerWrap{nil, constructors.NewCleaner()}
		var auth authenticator.IAuthenticator
		config, err := parser()

		if nil == err {
			auth, err = authConstructor.Construct(context, wrap.cleaner)
		}

		if nil == err {
			builder := builder.New().
				WithHost(config.Host).
				WithPort(config.Port).
				WithSwaggerSpecification(config.SwaggerURL)

			for _, e := range extenders {
				cors, cont := e(auth, context)

				if nil != cors {
					builder.WithCorsFiller(cors)
				}

				if nil != cont {
					builder.WithController(cont)
				}
			}

			wrap.app, err = builder.Build()
		}

		return wrap, err
	}
}

