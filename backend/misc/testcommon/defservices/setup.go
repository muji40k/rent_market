package defservices

import (
	"rent_service/builders/mothers/test/service/defservices"
	"rent_service/internal/factory/services/v1/deffactory"
	rv1 "rent_service/internal/repository/context/v1"
	"rent_service/misc/nullable"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type Context struct {
	Factory       *deffactory.Factory
	PhotoRegistry *defservices.PhotoRegistry
}

func New(factories rv1.Factories) Context {
	return Context{
		nullable.UnwrapF(defservices.DefaultServiceFactory(factories).Build()),
		defservices.NewPhotoRegistry(),
	}
}

func (self *Context) Close() {
	if nil != self.Factory {
		self.Factory.Clear()
	}
}

func (self *Context) SetUp(t provider.T, factories rv1.Factories) {
	t.WithNewStep("Create service factory", func(sCtx provider.StepCtx) {
		var err error
		self.Factory, err = defservices.DefaultServiceFactory(factories).Build()

		if nil != err {
			t.Breakf("Unable to create service: %s", err)
		}
	})

	t.WithNewStep("Create photo helper", func(sCtx provider.StepCtx) {
		self.PhotoRegistry = defservices.NewPhotoRegistry()
	})
}

func (self *Context) TearDown(t provider.T) {
	t.WithNewStep("Clear service", func(sCtx provider.StepCtx) {
		if nil != self.Factory {
			self.Factory.Clear()
		}
	})
}

