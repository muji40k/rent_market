package collect

type IBuilder[T any] interface {
	Build() T
}

func Do[T any, B IBuilder[T]](builders ...B) []T {
	out := make([]T, len(builders))

	for i, v := range builders {
		out[i] = v.Build()
	}

	return out
}

