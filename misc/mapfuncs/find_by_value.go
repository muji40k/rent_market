package mapfuncs

func FindByValue[K comparable, V comparable](m map[K]V, v V) (K, bool) {
	return FindByValueF(m, func(value *V) bool { return *value == v })
}

func FindByValueF[K comparable, V any](m map[K]V, f func(*V) bool) (K, bool) {
	for key, value := range m {
		if f(&value) {
			return key, true
		}
	}

	var empty K
	return empty, false
}

