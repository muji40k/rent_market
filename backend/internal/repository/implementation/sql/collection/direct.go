package collection

import (
	"github.com/jmoiron/sqlx"
	"rent_service/internal/misc/types/collection"
)

type sql_direct_collection[T any] struct {
	query Query
}

func NewDirect[T any](query Query) collection.Collection[T] {
	return &sql_direct_collection[T]{query}
}

func (self *sql_direct_collection[T]) Iter() collection.Iterator[T] {
	return newDirectIterator[T](self.query)
}

type direct_iterator[T any] struct {
	query  Query
	rows   *sqlx.Rows
	offset uint
	broken bool
}

func newDirectIterator[T any](query Query) collection.Iterator[T] {
	return &direct_iterator[T]{query, nil, 0, false}
}

func (self *direct_iterator[T]) getRows() (*sqlx.Rows, error) {
	if nil != self.rows {
		return self.rows, nil
	}

	var err error
	self.rows, err = self.query(self.offset)

	if nil != err {
		self.broken = true
	}

	return self.rows, err
}

func (self *direct_iterator[T]) Next() (T, bool) {
	var value T
	var empty T
	var next bool

	if self.broken {
		return value, false
	}

	if rows, err := self.getRows(); nil != err {
		next = false
	} else {
		if next = rows.Next(); next {
			if nil != rows.Scan(&value) {
				value = empty
			}
		}
	}

	return value, next
}

func (self *direct_iterator[T]) Skip() bool {
	if self.broken {
		return false
	}

	next := false

	if nil != self.rows {
		next = self.rows.Next()
	} else {
		self.offset++
		next = true
	}

	return next
}

