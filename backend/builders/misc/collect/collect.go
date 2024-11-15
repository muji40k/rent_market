package collect

import "fmt"

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

func FmtWrap[T any, B IBuilder[T]](buildersf func(string) B) func(uint) B {
	return func(i uint) B {
		return buildersf(fmt.Sprint(i))
	}
}

func DoN[T any, B IBuilder[T]](rng uint, buildersf func(uint) B) []T {
	out := make([]T, rng)

	for i := range rng {
		out[i] = buildersf(i).Build()
	}

	return out
}

