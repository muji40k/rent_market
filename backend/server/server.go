package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"rent_service/logger"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swag "github.com/go-openapi/runtime/middleware"
)

type Server struct {
	engine     *gin.Engine
	srv        *http.Server
	host       string
	port       uint
	corsConfig cors.Config
	logger     logger.ILogger
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

func WithLogger(log logger.ILogger) Configurator {
	return func(server *Server) {
		server.logger = log
		server.engine.Use(func(ctx *gin.Context) {
			start := time.Now()
			path := ctx.Request.URL.Path
			raw := ctx.Request.URL.RawQuery

			ctx.Next()

			duration := time.Now().Sub(start)
			client := ctx.ClientIP()
			method := ctx.Request.Method
			status := ctx.Writer.Status()
			bodySize := ctx.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}

			logger.Logf(log, logger.INFO,
				"Request processed: %v - %v %v; client: %v; body: %v; duration: %v",
				status, method, path, client, bodySize, duration,
			)
		})
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
	self.srv = &http.Server{
		Addr:    fmt.Sprintf("%v:%v", self.host, self.port),
		Handler: self.engine.Handler(),
	}

	go func() {
		if err := self.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logf(self.logger, logger.ERROR, "Listen error: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log(self.logger, logger.INFO, "Server shutting down")
}

func (self *Server) Clear() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := self.srv.Shutdown(ctx); err != nil {
		logger.Logf(self.logger, logger.ERROR, "Shutdown error: %s", err)
	}

	select {
	case <-ctx.Done():
		logger.Log(self.logger, logger.INFO, "Context timeout")
	}

	logger.Log(self.logger, logger.INFO, "Server down")
}

func (self *Server) Extend(controller IController) *Server {
	controller.Register(self.engine)

	return self
}

