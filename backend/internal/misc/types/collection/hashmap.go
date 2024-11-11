package collection

type KV[K comparable, V any] struct {
	Key   K
	Value V
}

func getKeys[K comparable, V any](m map[K]V) []K {
	out := make([]K, len(m))
	i := 0

	for k := range m {
		out[i] = k
		i++
	}

	return out
}

func HashMapIterator[K comparable, V any](m map[K]V) Iterator[KV[K, V]] {
	if nil == m {
		return nil
	}

	return &hmapIterator[K, V]{SliceIterator(getKeys(m)), m}
}

type hmapIterator[K comparable, V any] struct {
	keys Iterator[K]
	hmap map[K]V
}

func (self *hmapIterator[K, V]) Next() (KV[K, V], bool) {
	if key, next := self.keys.Next(); next {
		return KV[K, V]{key, self.hmap[key]}, true
	} else {
		return KV[K, V]{}, false
	}
}

func (self *hmapIterator[K, V]) Skip() bool {
	return self.keys.Skip()
}

