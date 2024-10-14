package emptymathcer

import "rent_service/internal/logic/services/errors/cmnerrors"

type Pair struct {
	key   string
	value string
}

func Item(key string, value string) Pair {
	return Pair{key, value}
}

func Match(pairs ...Pair) error {
	empty := make([]string, 0, len(pairs))

	for _, pair := range pairs {
		if "" == pair.value {
			empty = append(empty, pair.key)
		}
	}

	if 0 == len(empty) {
		return nil
	} else {
		return cmnerrors.Empty(empty...)
	}
}

