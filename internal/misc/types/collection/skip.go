package collection

type skipIterator[T any] struct {
	iter Iterator[T]
	size uint
	last bool
}

func SkipIterator[T any](skip uint, iterator Iterator[T]) Iterator[T] {
	return &skipIterator[T]{iterator, skip, true}
}

func (self *skipIterator[T]) initial_skip() bool {
	for ; 0 != self.size; self.size-- {
		if !self.iter.Skip() {
			return false
		}
	}

	return true
}

func (self *skipIterator[T]) Next() (T, bool) {
	var out T

	if !self.last {
		return out, false
	}

	out, self.last = self.next()

	return out, self.last
}

func (self *skipIterator[T]) next() (T, bool) {
	if 0 != self.size && !self.initial_skip() {
		var out T
		return out, false
	}

	return self.iter.Next()
}

func (self *skipIterator[T]) Skip() bool {
	if !self.last {
		return false
	}

	self.last = self.skip()

	return self.last
}

func (self *skipIterator[T]) skip() bool {
	if 0 != self.size && !self.initial_skip() {
		return false
	}

	return self.iter.Skip()
}

