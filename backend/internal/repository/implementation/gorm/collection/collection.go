package collection

import (
	"database/sql"
	"rent_service/internal/misc/types/collection"

	"gorm.io/gorm"
)

type gorm_collection[T any] struct {
	db    *gorm.DB
	query *gorm.DB
}

func New[T any](db *gorm.DB, query *gorm.DB) collection.Collection[T] {
	return &gorm_collection[T]{db, query.Session(&gorm.Session{})}
}

func (self *gorm_collection[T]) Iter() collection.Iterator[T] {
	return newIterator[T](self.db, self.query)
}

type iterator[T any] struct {
	db     *gorm.DB
	query  *gorm.DB
	rows   *sql.Rows
	offset uint
	broken bool
}

func newIterator[T any](db *gorm.DB, query *gorm.DB) collection.Iterator[T] {
	return &iterator[T]{db, query, nil, 0, false}
}

func (self *iterator[T]) getRows() (*sql.Rows, error) {
	if nil != self.rows {
		return self.rows, nil
	}

	var err error
	self.rows, err = self.query.Offset(int(self.offset)).Rows()

	if nil != err {
		self.broken = true
	}

	return self.rows, err
}

func (self *iterator[T]) Next() (T, bool) {
	var next bool
	var value T
	var empty T

	if self.broken {
		return empty, false
	}

	if rows, err := self.getRows(); nil != err {
		next = false
	} else {
		if next = rows.Next(); next {
			if nil != self.db.ScanRows(rows, &value) {
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

