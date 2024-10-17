package collection

import (
	"github.com/jmoiron/sqlx"
	"rent_service/internal/misc/types/collection"
)

type Query func(offset uint) (*sqlx.Rows, error)

type sql_collection[T any] struct {
	query Query
}

func New[T any](query Query) collection.Collection[T] {
	return &sql_collection[T]{query}
}

func (self *sql_collection[T]) Iter() collection.Iterator[T] {
	return newIterator[T](self.query)
}

type iterator[T any] struct {
	query  Query
	rows   *sqlx.Rows
	offset uint
	broken bool
}

func newIterator[T any](query Query) collection.Iterator[T] {
	return &iterator[T]{query, nil, 0, false}
}

func (self *iterator[T]) getRows() (*sqlx.Rows, error) {
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

func (self *iterator[T]) Next() (T, bool) {
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
			if nil != rows.StructScan(&value) {
				value = empty
			}
		}
	}

	return value, next
}

func (self *iterator[T]) Skip() bool {
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

