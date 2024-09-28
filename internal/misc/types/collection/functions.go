package collection

func Collect[T any](iterator Iterator[T]) []T {
	var out []T

	for value, next := iterator.Next(); next; {
		out = append(out, value)
	}

	return out
}

