package collection

type filterIterator[T any] struct {
	filter func(*T) bool
	iter   Iterator[T]
}

func FilterIterator[T any](
	filter func(*T) bool,
	iterator Iterator[T],
) Iterator[T] {
	return &filterIterator[T]{filter, iterator}
}

func (self *filterIterator[T]) Next() (T, bool) {
	for value, next := self.iter.Next(); next; value, next = self.iter.Next() {
		if self.filter(&value) {
			return value, true
		}
	}

	var empty T
	return empty, false
}

func (self *filterIterator[T]) Skip() bool {
	for value, next := self.iter.Next(); next; value, next = self.iter.Next() {
		if self.filter(&value) {
			return true
		}
	}

	return false
}

