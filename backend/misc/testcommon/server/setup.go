package server

import (
	"net/http"
	"rent_service/builders/mothers/test/application/server"
	"rent_service/internal/logic/context/v1"
	"rent_service/internal/misc/types/collection"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/stretchr/testify/assert"
)

type Context struct {
	Server   *gin.Engine
	Inserter *server.Inserter
}

func New(
	factories v1.Factories,
	controllers ...server.ControllerCreator,
) Context {
	return Context{
		server.TestServer(factories, controllers...),
		server.NewInserter(),
	}
}

func (self *Context) Close() {
	if nil != self.Inserter {
		self.Inserter.Close()
	}
}

func (self *Context) SetUp(
	t provider.T,
	factories v1.Factories,
	controllers ...server.ControllerCreator,
) {
	t.WithNewStep("Create server", func(sCtx provider.StepCtx) {
		self.Server = server.TestServer(factories, controllers...)
	})

	t.WithNewStep("Create session database helper", func(sCtx provider.StepCtx) {
		self.Inserter = server.NewInserter()
	})
}

func (self *Context) TearDown(t provider.T) {
	t.WithNewStep("Close session database connection", func(sCtx provider.StepCtx) {
		if nil != self.Inserter {
			self.Inserter.Close()
		}
	})
}

type Provider interface {
	assert.TestingT
	httpexpect.Logger
}

func (self *Context) GetClient(
	t Provider,
	printers ...func(Provider) httpexpect.Printer,
) *httpexpect.Expect {
	var iter collection.Iterator[func(Provider) httpexpect.Printer]

	if 0 != len(printers) {
		iter = collection.SliceIterator(printers)
	} else {
		iter = collection.SingleIterator(func(t Provider) httpexpect.Printer {
			return httpexpect.NewCompactPrinter(t)
		})
	}

	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(self.Server),
			Jar:       httpexpect.NewCookieJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: collection.Collect(collection.MapIterator(
			func(ctor *func(Provider) httpexpect.Printer) httpexpect.Printer {
				return (*ctor)(t)
			},
			iter,
		)),
	})
}

