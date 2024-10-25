package collection

type singleIterator[T any] struct {
	value T
	read  bool
}

func SingleIterator[T any](value T) Iterator[T] {
	return &singleIterator[T]{value, false}
}

func (self *singleIterator[T]) Next() (T, bool) {
	if self.read {
		return self.value, false
	}

	self.read = true

	return self.value, true
}

func (self *singleIterator[T]) Skip() bool {
	if self.read {
		return false
	}

	self.read = true

	return true
}

