package web

import (
	"rent_service/application"
	"rent_service/internal/misc/types/collection"
	"rent_service/server"
)

type Builder struct {
	host            string
	port            uint
	swaggerSpecsUrl string
	corsFillers     []server.CorsFiller
	controllers     []server.IController
}

func New() *Builder {
	return &Builder{}
}

func (self *Builder) WithHost(host string) *Builder {
	self.host = host
	return self
}

func (self *Builder) WithPort(port uint) *Builder {
	self.port = port
	return self
}

func (self *Builder) WithCorsFiller(filler server.CorsFiller) *Builder {
	self.corsFillers = append(self.corsFillers, filler)
	return self
}

func (self *Builder) WithCorsFillers(corsFillers []server.CorsFiller) *Builder {
	if nil != corsFillers {
		self.corsFillers = corsFillers
	}

	return self
}

func (self *Builder) WithController(controller server.IController) *Builder {
	self.controllers = append(self.controllers, controller)
	return self
}

func (self *Builder) WithControllers(controllers []server.IController) *Builder {
	if nil != controllers {
		self.controllers = controllers
	}

	return self
}

func (self *Builder) WithSwaggerSpecification(url string) *Builder {
	self.swaggerSpecsUrl = url

	return self
}

func (self *Builder) Build() (application.IApplication, error) {
	configurators := collection.ChainIterator(
		collection.MapIterator(
			func(host *string) server.Configurator {
				return server.WithHost(*host)
			},
			collection.SingleIterator(self.host),
		),
		collection.ChainIterator(
			collection.MapIterator(
				func(url *string) server.Configurator {
					return server.WithSwaggerSpecification(*url)
				},
				collection.FilterIterator(
					func(url *string) bool {
						return "" != *url
					},
					collection.SingleIterator(self.swaggerSpecsUrl),
				),
			),
			collection.ChainIterator(
				collection.MapIterator(
					func(port *uint) server.Configurator {
						return server.WithPort(*port)
					},
					collection.SingleIterator(self.port),
				),
				collection.ChainIterator(
					collection.MapIterator(
						func(cors *[]server.CorsFiller) server.Configurator {
							return server.WithCors(*cors...)
						},
						collection.SingleIterator(self.corsFillers),
					),
					collection.MapIterator(
						func(controller *server.IController) server.Configurator {
							return server.WithController(*controller)
						},
						collection.FilterIterator(
							func(controller *server.IController) bool {
								return nil != *controller
							},
							collection.SliceIterator(self.controllers),
						),
					),
				),
			),
		),
	)

	return server.New(collection.Collect(configurators)...), nil
}

