package collection

import (
	"bytes"
	"encoding/json"
)

type wrap[T any] struct {
	iter Iterator[T]
}

func Marshaler[T any](iter Iterator[T]) json.Marshaler {
	return &wrap[T]{iter}
}

func (self *wrap[T]) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteRune('[')

	for item, next := self.iter.Next(); next; {
		content, err := json.Marshal(item)

		if err != nil {
			return []byte{}, err
		}

		buf.Write(content)
		item, next = self.iter.Next()

		if next {
			buf.WriteRune(',')
		}
	}

	buf.WriteRune(']')

	return buf.Bytes(), nil
}

