package collection

type map_collection[T any, F any] struct {
	col  Collection[T]
	mapf func(*T) F
}

type map_iterator[T any, F any] struct {
	iter Iterator[T]
	mapf func(*T) F
}

func MapCollection[T any, F any](mapf func(*T) F, col Collection[T]) Collection[F] {
	if nil == col {
		return nil
	}

	return &map_collection[T, F]{col, mapf}
}

func (self *map_collection[T, F]) Iter() Iterator[F] {
	return MapIterator(self.mapf, self.col.Iter())
}

func MapIterator[T any, F any](mapf func(*T) F, iter Iterator[T]) Iterator[F] {
	if nil == iter {
		return nil
	}

	return &map_iterator[T, F]{iter, mapf}
}

func (self *map_iterator[T, F]) Next() (F, bool) {
	if v, next := self.iter.Next(); next {
		return self.mapf(&v), true
	} else {
		var out F
		return out, false
	}
}

func (self *map_iterator[T, F]) Skip() bool {
	return self.iter.Skip()
}

