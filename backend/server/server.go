package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"rent_service/logger"
	"rent_service/misc/contextholder"
	"rent_service/server/errstructs"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swag "github.com/go-openapi/runtime/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	// "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
type EngineExtender func(engine *gin.Engine)

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
func WithExtender(ext EngineExtender) Configurator {
	return func(server *Server) {
		ext(server.engine)
	}
}

func SwaggerSpecificationExtender(url string) EngineExtender {
	return func(engine *gin.Engine) {
		h := swag.SwaggerUI(swag.SwaggerUIOpts{
			BasePath: "/",
			Path:     "docs", // AAAAAAAAAAAAAA
			SpecURL:  url,
		}, nil)

		engine.GET("/", func(c *gin.Context) {
			c.Request.URL.Path = "/docs" // So... where is no other context to create deep copy
			h.ServeHTTP(c.Writer, c.Request)
		})
	}
}

type CorsFiller func(config *cors.Config)

func CorsExtender(fillers ...CorsFiller) EngineExtender {
	return func(engine *gin.Engine) {
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

		engine.Use(cors.New(config))
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

			duration := time.Since(start)
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

func TracerExtender(tr trace.TracerProvider, hl *contextholder.Holder) EngineExtender {
	return func(engine *gin.Engine) {
		if nil != hl {
			// Context holder is going to be wrapped around, so asyncronous
			// execution will break the stack
			var mu sync.Mutex
			engine.Use(func(ctx *gin.Context) {
				mu.Lock()
				defer mu.Unlock()
				ctx.Next()
			})
		}

		engine.Use(otelgin.Middleware(
			"rent_service",
			otelgin.WithTracerProvider(tr),
		))

		if nil != hl {
			tracer := tr.Tracer("rent_service")
			engine.Use(func(ctx *gin.Context) {
				var span trace.Span
				var err error
				var perr error
				serr := hl.Start(ctx.Request.Context())
				defer func() {
					if nil != span {
						span.End()
					}

					if nil == perr {
						hl.Pop()
					}

					if nil == serr {
						hl.Pop()
					}
				}()

				if nil != serr {
					err = serr
				}

				if nil == err {
					perr = hl.Push(func(ctx context.Context) (context.Context, error) {
						var nctx context.Context
						nctx, span = tracer.Start(ctx, "HTTP handler")
						return nctx, nil
					})
					err = perr
				}

				if nil == err {
					span.AddEvent("Actual start")
					ctx.Next()
					span.AddEvent("Actual end")
					span.SetStatus(codes.Ok, "No error occured")
				} else {
					span.SetStatus(codes.Error, "Wrap error")
					span.SetAttributes(attribute.String("Error", err.Error()))
					ctx.JSON(
						http.StatusInternalServerError,
						errstructs.NewInternalErr(err),
					)
				}
			})
		}
	}
}

func New(config ...Configurator) *Server {
	out := Server{
		engine:     gin.New(),
		host:       "localhost",
		port:       80,
		corsConfig: cors.Config{},
	}

	_ = out.engine.SetTrustedProxies(nil)
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
	logger.Log(self.logger, logger.INFO, "Server detached")
}

func (self *Server) Clear() {
	logger.Log(self.logger, logger.INFO, "Server shutting down")

	if err := self.srv.Shutdown(context.Background()); err != nil {
		logger.Logf(self.logger, logger.ERROR, "Shutdown error: %s", err)
	}

	logger.Log(self.logger, logger.INFO, "Server down")
}

func (self *Server) Extend(controller IController) *Server {
	controller.Register(self.engine)

	return self
}

