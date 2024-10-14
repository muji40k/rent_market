package collection

type sliceCollection[T any] []T

func SliceCollection[T any](slice []T) Collection[T] {
	return sliceCollection[T](slice)
}

type sliceIterator[T any] struct {
	data  []T
	index uint
}

func (self sliceCollection[T]) Iter() Iterator[T] {
	return &sliceIterator[T]{
		data:  self,
		index: 0,
	}
}

func (self *sliceIterator[T]) Next() (T, bool) {
	if uint(len(self.data)) <= self.index {
		var out T
		return out, false
	}

	index := self.index
	self.index++

	return self.data[index], true
}

func (self *sliceIterator[T]) Skip() bool {
	if uint(len(self.data)) <= self.index {
		return false
	}

	self.index++

	return true
}

