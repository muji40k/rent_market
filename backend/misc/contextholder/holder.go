package contextholder

import (
	"context"
	"errors"
	"fmt"
)

type Holder struct {
	ctxs []context.Context
}

func New() *Holder {
	return &Holder{
		make([]context.Context, 0),
	}
}

func (self *Holder) Start(ctx context.Context) error {
	if 0 != len(self.ctxs) {
		return ErrorAlreadyRunning
	}

	self.ctxs = append(self.ctxs, ctx)

	return nil
}

func (self *Holder) Push(ctor func(context.Context) (context.Context, error)) error {
	if 0 == len(self.ctxs) {
		return ErrorNotStarted
	}

	next, err := ctor(self.ctxs[len(self.ctxs)-1])

	if nil == err {
		self.ctxs = append(self.ctxs, next)
	} else {
		err = ErrorCreator{err}
	}

	return err
}

func (self *Holder) Pop() (context.Context, error) {
	if 0 == len(self.ctxs) {
		return nil, ErrorNotStarted
	}

	i := len(self.ctxs) - 1
	last := self.ctxs[i]
	self.ctxs = self.ctxs[:i]

	return last, nil
}

func (self *Holder) Parallel(ctor func(context.Context) error) error {
	if 0 == len(self.ctxs) {
		return ErrorNotStarted
	}

	if err := ctor(self.ctxs[len(self.ctxs)-1]); nil == err {
		return nil
	} else {
		return ErrorCreator{err}
	}
}

var ErrorAlreadyRunning = errors.New("Holder already running")
var ErrorNotStarted = errors.New("Holder not initialised running")

type ErrorCreator struct {
	Err error
}

func (self ErrorCreator) Error() string {
	return fmt.Sprintf("Holder creator error: %s", self.Err)
}

func (self ErrorCreator) Unwrap() error {
	return self.Err
}

