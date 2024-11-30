package generator

import (
	"rent_service/misc/nullable"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type command int

const (
	GENERATE command = iota
	FINISH
)

type IAllureProvider interface {
	WithNewAsyncStep(stepName string, step func(sCtx provider.StepCtx), params ...*allure.Parameter)
}

type Spy struct {
	name  string
	value any
}

func (self *Spy) SniffValue(name string, value any) {
	self.name = name
	self.value = value
}

type AllureGeneratorWrap struct {
	ctrl chan command
	res  chan uuid.UUID
}

func NewAllureWrap(
	t IAllureProvider,
	name string,
	generator IGenerator,
	spy *nullable.Nullable[Spy],
) IGenerator {
	c := make(chan command, 1)
	r := make(chan uuid.UUID, 1)

	t.WithNewAsyncStep(name, func(sCtx provider.StepCtx) {
		for {
			switch <-c {
			case GENERATE:
				r <- generator.Generate()
			case FINISH:
				generator.Finish()
				nullable.IfSome(spy, func(spy *Spy) {
					sCtx.WithParameters(
						allure.NewParameter(spy.name, spy.value),
					)
				})
				r <- uuid.UUID{}
				return
			}
		}
	})

	return &AllureGeneratorWrap{c, r}
}

func (self *AllureGeneratorWrap) Finish() {
	self.ctrl <- FINISH
	<-self.res
	close(self.ctrl)
	close(self.res)
}

func (self *AllureGeneratorWrap) Generate() uuid.UUID {
	self.ctrl <- GENERATE
	return <-self.res
}

