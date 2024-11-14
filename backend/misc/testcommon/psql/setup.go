package psql

import (
	"rent_service/builders/mothers/test/repository/psql"
	psqlfactory "rent_service/internal/factory/repositories/v1/psql"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type Context struct {
	Inserter *psql.Inserter
	Factory  *psqlfactory.Factory
}

func (self *Context) SetUp(t provider.T) {
	t.WithNewStep("Create insert helper", func(sCtx provider.StepCtx) {
		self.Inserter = psql.NewInserter()
	})

	t.WithNewStep("Create factory", func(sCtx provider.StepCtx) {
		var err error
		self.Factory, err = psql.PSQLRepositoryFactory().Build()

		if nil != err {
			t.Breakf("Unable to create repository: %s", err)
		}
	})
}

func (self *Context) TearDown(t provider.T) {
	t.WithNewStep("Close connections", func(sCtx provider.StepCtx) {
		self.Inserter.Close()
		self.Factory.Clear()
	})
}

