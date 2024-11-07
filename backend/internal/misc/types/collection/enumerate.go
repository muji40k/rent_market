package collection

func EnumerateIterator[T any](iter Iterator[T]) Iterator[Pair[uint, T]] {
	return &enumerateIterator[T]{iter, 0}
}

type enumerateIterator[T any] struct {
	iter  Iterator[T]
	count uint
}

func (self *enumerateIterator[T]) Next() (Pair[uint, T], bool) {
	i := self.count
	v, next := self.iter.Next()

	self.count++

	return Pair[uint, T]{i, v}, next
}

func (self *enumerateIterator[T]) Skip() bool {
	return self.iter.Skip()
}

