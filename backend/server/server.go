package server

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swag "github.com/go-openapi/runtime/middleware"
)

type Server struct {
	engine     *gin.Engine
	host       string
	port       uint
	corsConfig cors.Config
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
func WithController(controller IController) Configurator {
	return func(server *Server) {
		server.Extend(controller)
	}
}
func WithSwaggerSpecification(url string) Configurator {
	return func(server *Server) {
		h := swag.SwaggerUI(swag.SwaggerUIOpts{
			BasePath: "/",
			Path:     "docs", // AAAAAAAAAAAAAA
			SpecURL:  url,
		}, nil)

		server.engine.GET("/", func(c *gin.Context) {
			c.Request.URL.Path = "/docs" // So... where is no other context to create deep copy
			h.ServeHTTP(c.Writer, c.Request)
		})
	}
}

type CorsFiller func(config *cors.Config)

func WithCors(fillers ...CorsFiller) Configurator {
	return func(server *Server) {
		config := cors.Config{
			AllowHeaders: []string{
				"Origin",
				"Content-Length",
				"Content-Type",
			},
			AllowCredentials: false,
			AllowAllOrigins:  true,
			MaxAge:           12 * time.Hour,
		}

		for _, filler := range fillers {
			filler(&config)
		}

		server.engine.Use(cors.New(config))
	}
}

func New(config ...Configurator) *Server {
	out := Server{
		engine:     gin.New(),
		host:       "localhost",
		port:       80,
		corsConfig: cors.Config{},
	}

	out.engine.SetTrustedProxies(nil)
	out.engine.Use(gin.Recovery())

	for _, f := range config {
		f(&out)
	}

	return &out
}

func (self *Server) Run() {
	self.engine.Run(fmt.Sprintf("%v:%v", self.host, self.port))
}

func (self *Server) Clear() {}

func (self *Server) Extend(controller IController) *Server {
	controller.Register(self.engine)

	return self
}

