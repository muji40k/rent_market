package generator

import "github.com/google/uuid"

type BFunc func(map[string]GFunc) uuid.UUID
type FFunc func()
type GFunc func() uuid.UUID

type FuncGenerator struct {
	base       BFunc
	finish     FFunc
	generators map[string]GFunc
}

type Gpair struct {
	name      string
	generator GFunc
}

func Pair(name string, generator GFunc) Gpair {
	return Gpair{name, generator}
}

type wrap struct {
	base   BFunc
	finish FFunc
}

func NewFuncWrapped(w wrap, gens ...Gpair) IGenerator {
	return NewFunc(w.base, w.finish, gens...)
}

func NewFunc(base BFunc, finish FFunc, gens ...Gpair) IGenerator {
	m := make(map[string]GFunc, len(gens))

	for _, p := range gens {
		m[p.name] = p.generator
	}

	return &FuncGenerator{base, finish, m}
}

func (self *FuncGenerator) Generate() uuid.UUID {
	return self.base(self.generators)
}

func (self *FuncGenerator) Finish() {
	self.finish()
}

func FuncListWrap[T any](dst *[]T, base func(map[string]GFunc) (T, uuid.UUID), finish func([]T)) wrap {
	return wrap{
		func(gens map[string]GFunc) uuid.UUID {
			c, id := base(gens)
			*dst = append(*dst, c)
			return id
		},
		func() { finish(*dst) },
	}
}

