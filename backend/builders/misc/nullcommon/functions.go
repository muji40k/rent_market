package nullcommon

import "rent_service/misc/nullable"

func CopyPtrIfSome[T any](value *nullable.Nullable[T]) *T {
	return nullable.GerOrInsert(
		nullable.Map(value, func(ref *T) *T {
			out := new(T)
			*out = *ref
			return out
		}),
		nil,
	)
}

