package collection

type Collection[T any] interface {
	Iter() Iterator[T]
}

type Iterator[T any] interface {
	Next() (T, bool)
	Skip() bool
}

