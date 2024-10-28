package collection

func ChainIterator[T any](iter1 Iterator[T], iter2 Iterator[T]) Iterator[T] {
	if nil == iter1 && nil == iter2 {
		return nil
	}

	if nil == iter1 {
		return iter2
	}

	if nil == iter2 {
		return iter1
	}

	return &chainIterator[T]{true, iter1, iter2}
}

type chainIterator[T any] struct {
	first bool
	iter1 Iterator[T]
	iter2 Iterator[T]
}

func (self *chainIterator[T]) Next() (T, bool) {
	if self.first {
		if v, next := self.iter1.Next(); next {
			return v, true
		}

		self.first = false
	}

	return self.iter2.Next()
}

func (self *chainIterator[T]) Skip() bool {
	if self.first {
		if next := self.iter1.Skip(); next {
			return true
		}

		self.first = false
	}

	return self.iter2.Skip()
}

