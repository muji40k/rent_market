package collection

type Pair[A any, B any] struct {
	A A
	B B
}

func ZipIterator[A any, B any](iter1 Iterator[A], iter2 Iterator[B]) Iterator[Pair[A, B]] {
	if nil == iter1 || nil == iter2 {
		return nil
	}

	return &zipIterator[A, B]{false, iter1, iter2}
}

type zipIterator[A any, B any] struct {
	end   bool
	iter1 Iterator[A]
	iter2 Iterator[B]
}

func (self *zipIterator[A, B]) Next() (Pair[A, B], bool) {
	var pair Pair[A, B]
	var anext, bnext bool

	if self.end {
		return pair, false
	}

	pair.A, anext = self.iter1.Next()
	pair.B, bnext = self.iter2.Next()

	self.end = !anext || !bnext

	return pair, !self.end
}

func (self *zipIterator[A, B]) Skip() bool {
	if self.end {
		return false
	}

	self.end = !self.iter1.Skip() || !self.iter2.Skip()

	return !self.end
}

