package collection

type Collection[T any] interface {
	Iter() Iterator[T]
}

type Iterator[T any] interface {
	Next() (T, bool) // Each call after reaching (T, false) must return false
	Skip() bool
}

