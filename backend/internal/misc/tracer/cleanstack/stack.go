package cleanstack

import (
	"context"
	"log"
)

type ICleaner interface {
	Clean(context.Context)
}

type CleanerF func(context.Context)
type CleanerFE[E error] func(context.Context) E

func (self CleanerF) Clean(ctx context.Context) {
	self(ctx)
}

func (self CleanerFE[E]) Clean(ctx context.Context) {
	err := self(ctx)

	if nil != error(err) {
		log.Println(err)
	}
}

type Cleaner struct {
	cleaners []ICleaner
}

func New() *Cleaner {
	return &Cleaner{make([]ICleaner, 0)}
}

func (self *Cleaner) Push(c ICleaner) *Cleaner {
	if nil != c {
		self.cleaners = append(self.cleaners, c)
	}

	return self
}

func (self *Cleaner) Clean(ctx context.Context) {
	l := len(self.cleaners)
	for i := range l {
		self.cleaners[l-i-1].Clean(ctx)
	}
}

func PushF(self *Cleaner, f func(context.Context)) *Cleaner {
	if nil == f {
		return self
	}

	return self.Push(CleanerF(f))
}

func PushFE[E error](self *Cleaner, f func(context.Context) E) *Cleaner {
	if nil == f {
		return self
	}

	return self.Push(CleanerFE[E](f))
}

