package collection

type limitIterator[T any] struct {
	iter  Iterator[T]
	limit uint
	index uint
}

func LimitIterator[T any](limit uint, iterator Iterator[T]) Iterator[T] {
	if nil == iterator {
		return nil
	}

	return &limitIterator[T]{iterator, limit, 0}
}

func (self *limitIterator[T]) Next() (T, bool) {
	if self.limit <= self.index {
		var out T
		return out, false
	}

	self.index++

	return self.iter.Next()
}

func (self *limitIterator[T]) Skip() bool {
	if self.limit <= self.index {
		return false
	}

	self.index++

	return self.iter.Skip()
}

