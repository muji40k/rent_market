package collection

import (
	"bytes"
	"encoding/json"
)

type marshaler[T any] struct {
	iter Iterator[T]
}

func Marshaler[T any](iter Iterator[T]) json.Marshaler {
	return &marshaler[T]{iter}
}

func (self *marshaler[T]) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteRune('[')

	if nil != self.iter {
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
	}

	buf.WriteRune(']')

	return buf.Bytes(), nil
}

