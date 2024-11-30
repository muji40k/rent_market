package generator

import "github.com/google/uuid"

type BFunc func() uuid.UUID
type FFunc func()

type FuncGenerator struct {
	base   BFunc
	finish FFunc
}

type wrap struct {
	base   BFunc
	finish FFunc
}

func NewFuncWrapped(w wrap) IGenerator {
	return NewFunc(w.base, w.finish)
}

func NewFunc(base BFunc, finish FFunc) IGenerator {
	return &FuncGenerator{base, finish}
}

func (self *FuncGenerator) Generate() uuid.UUID {
	return self.base()
}

func (self *FuncGenerator) Finish() {
	self.finish()
}

func FuncListWrap[T any](dst *[]T, base func() (T, uuid.UUID), finish func([]T)) wrap {
	return wrap{
		func() uuid.UUID {
			c, id := base()
			*dst = append(*dst, c)
			return id
		},
		func() { finish(*dst) },
	}
}

