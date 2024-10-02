package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	host   string
	port   uint
}

type IController interface {
	Register(engine *gin.Engine)
}

type Configurator func(server *Server)

func WithHost(host string) Configurator {
	return func(server *Server) {
		server.host = host
	}
}
func WithPort(port uint) Configurator {
	return func(server *Server) {
		server.port = port
	}
}

func New(config ...Configurator) Server {
	out := Server{
		engine: gin.New(),
		host:   "localhost",
		port:   80,
	}

	out.engine.Use(gin.Recovery())

	for _, f := range config {
		f(&out)
	}

	return out
}

func (self *Server) Run() {
	self.engine.Run(fmt.Sprintf("%v:%v", self.host, self.port))
}

func (self *Server) Extend(controller IController) {
	controller.Register(self.engine)
}

