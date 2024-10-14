package collection

func Collect[T any](iterator Iterator[T]) []T {
	var out []T

	for value, next := iterator.Next(); next; value, next = iterator.Next() {
		out = append(out, value)
	}

	return out
}

func Find[T any](iterator Iterator[T], f func(*T) bool) (T, bool) {
	for value, next := iterator.Next(); next; value, next = iterator.Next() {
		if f(&value) {
			return value, true
		}
	}

	var empty T
	return empty, false
}

func Count[T any](iterator Iterator[T]) uint {
	var count uint

	for next := iterator.Skip(); next; next = iterator.Skip() {
		count++
	}

	return count
}

func Reduce[T any](iterator Iterator[T], f func(*T, *T) T) T {
	var value T
	acc, next := iterator.Next()

	for next {
		value, next = iterator.Next()
		acc = f(&acc, &value)
	}

	return acc
}

func Fold[T any, F any](iterator Iterator[T], acc F, f func(*F, *T) F) F {
	for value, next := iterator.Next(); next; value, next = iterator.Next() {
		acc = f(&acc, &value)
	}

	return acc
}

