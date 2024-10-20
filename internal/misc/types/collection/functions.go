package collection

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
	if nil == iterator {
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

	for next := iterator.Skip(); next; next = iterator.Skip() {
		count++
	}

	return count
}

func Reduce[T any](iterator Iterator[T], f func(*T, *T) T) T {
	if nil == iterator {
		var empty T
		return empty
	}

	acc, _ := iterator.Next()

	for value, next := iterator.Next(); next; value, next = iterator.Next() {
		acc = f(&acc, &value)
	}

	return acc
}

func Fold[T any, F any](iterator Iterator[T], acc F, f func(*F, *T) F) F {
	if nil == iterator {
		return acc
	}

	for value, next := iterator.Next(); next; value, next = iterator.Next() {
		acc = f(&acc, &value)
	}

	return acc
}

