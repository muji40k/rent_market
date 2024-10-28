package constructors

type Cleaner struct {
	stages []func()
}

func NewCleaner() *Cleaner {
	return &Cleaner{nil}
}

func (self *Cleaner) AddStage(f func()) {
	if nil != f {
		self.stages = append(self.stages, f)
	}
}

func (self *Cleaner) Clean() {
	for _, f := range self.stages {
		f()
	}
}

