package generator

import "github.com/google/uuid"

type callGroup struct {
	call   func() uuid.UUID
	amount uint
}

type GeneratorGroup struct {
	generateCalls []callGroup
	finishCalls   []func()
}

func NewGeneratorGroup() *GeneratorGroup {
	return &GeneratorGroup{}
}

func (self *GeneratorGroup) Add(gen IGenerator, gtimes uint) *GeneratorGroup {
	return self.AddGenerate(gen, gtimes).AddFinish(gen)
}

func (self *GeneratorGroup) AddGenerate(gen IGenerator, times uint) *GeneratorGroup {
	self.generateCalls = append(self.generateCalls, callGroup{gen.Generate, times})

	return self
}

func (self *GeneratorGroup) AddFinish(gen IGenerator) *GeneratorGroup {
	self.finishCalls = append(self.finishCalls, gen.Finish)

	return self
}

func (self *GeneratorGroup) Generate() {
	for _, call := range self.generateCalls {
		for range call.amount {
			call.call()
		}
	}
}

func (self *GeneratorGroup) Finish() {
	for _, call := range self.finishCalls {
		call()
	}
}

