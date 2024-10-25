package collection

import "cmp"

func Collect[T any](iterator Iterator[T]) []T {
	if nil == iterator {
		return nil
	}

	var out []T

	for value, next := iterator.Next(); next; value, next = iterator.Next() {
		out = append(out, value)
	}

	return out
}

func Find[T any](iterator Iterator[T], f func(*T) bool) (T, bool) {
	var empty T
	if nil == iterator || nil == f {
		return empty, false
	}

	for value, next := iterator.Next(); next; value, next = iterator.Next() {
		if f(&value) {
			return value, true
		}
	}

	return empty, false
}

func Count[T any](iterator Iterator[T]) uint {
	if nil == iterator {
		return 0
	}

	var count uint

	for _, next := iterator.Next(); next; next = iterator.Skip() {
		count++
	}

	return count
}

func Max[T cmp.Ordered](iterator Iterator[T]) (T, bool) {
	return Reduce(iterator, func(a, b *T) T {
		if *a > *b {
			return *a
		} else {
			return *b
		}
	})
}

func Min[T cmp.Ordered](iterator Iterator[T]) (T, bool) {
	return Reduce(iterator, func(a, b *T) T {
		if *a < *b {
			return *a
		} else {
			return *b
		}
	})
}

func All[T any](iterator Iterator[T], f func(*T) bool) bool {
	if nil == iterator || nil == f {
		return false
	}

	for v, next := iterator.Next(); next; v, next = iterator.Next() {
		if !f(&v) {
			return false
		}
	}

	return true
}

func Any[T any](iterator Iterator[T], f func(*T) bool) bool {
	if nil == iterator || nil == f {
		return false
	}

	for v, next := iterator.Next(); next; v, next = iterator.Next() {
		if f(&v) {
			return true
		}
	}

	return false
}

func ForEach[T any](iterator Iterator[T], f func(*T)) {
	if nil == iterator || nil == f {
		return
	}

	for v, next := iterator.Next(); next; v, next = iterator.Next() {
		f(&v)
	}
}

func Reduce[T any](iterator Iterator[T], f func(*T, *T) T) (T, bool) {
	var empty T

	if nil == iterator || nil == f {
		return empty, false
	}

	acc, any := iterator.Next()

	if !any {
		return empty, false
	}

	for value, next := iterator.Next(); next; value, next = iterator.Next() {
		acc = f(&acc, &value)
	}

	return acc, true
}

func Fold[T any, F any](iterator Iterator[T], acc F, f func(*F, *T) F) F {
	if nil == iterator || nil == f {
		return acc
	}

	for value, next := iterator.Next(); next; value, next = iterator.Next() {
		acc = f(&acc, &value)
	}

	return acc
}

