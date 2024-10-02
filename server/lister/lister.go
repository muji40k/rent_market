package lister

import (
	"github.com/google/uuid"
	"rent_service/internal/logic/services/types/token"
	. "rent_service/internal/misc/types/collection"
)

type Method[T any] func(token.Token) (Collection[T], error)
type BaseListMethod[T any] func(token.Token, uuid.UUID) (Collection[T], error)

func ListSingle[T any](rawId string, base BaseListMethod[T]) (Method[T], error) {
	var method Method[T]
	id, err := uuid.Parse(rawId)

	if nil == err {
		method = func(token token.Token) (Collection[T], error) {
			col, err := base(token, id)

			if nil == err {
				return col, nil
			} else {
				return nil, err
			}
		}
	}

	return method, err
}

type BaseGetMethod[T any] func(token.Token, uuid.UUID) (T, error)

func ListMultiple[T any](rawIds []string, base BaseGetMethod[T]) (Method[T], error) {
	var err error
	var method Method[T]
	ids := make([]uuid.UUID, len(rawIds))

	for i := 0; nil == err && len(rawIds) > i; i++ {
		ids[i], err = uuid.Parse(rawIds[i])
	}

	if nil == err {
		method = func(token token.Token) (Collection[T], error) {
			var err error
			items := make([]T, len(ids))

			for i := 0; nil == err && len(ids) > i; i++ {
				items[i], err = base(token, ids[i])
			}

			if nil == err {
				return SliceCollection[T](items), nil
			} else {
				return nil, err
			}
		}
	}

	return method, err
}

type MethodNA[T any] func() (Collection[T], error)
type BaseListMethodNA[T any] func(uuid.UUID) (Collection[T], error)

func ListSingleNA[T any](rawId string, base BaseListMethodNA[T]) (MethodNA[T], error) {
	var method MethodNA[T]
	id, err := uuid.Parse(rawId)

	if nil == err {
		method = func() (Collection[T], error) {
			col, err := base(id)

			if nil == err {
				return col, nil
			} else {
				return nil, err
			}
		}
	}

	return method, err
}

type BaseGetMethodNA[T any] func(uuid.UUID) (T, error)

func ListMultipleNA[T any](rawIds []string, base BaseGetMethodNA[T]) (MethodNA[T], error) {
	var err error
	var method MethodNA[T]
	ids := make([]uuid.UUID, len(rawIds))

	for i := 0; nil == err && len(rawIds) > i; i++ {
		ids[i], err = uuid.Parse(rawIds[i])
	}

	if nil == err {
		method = func() (Collection[T], error) {
			var err error
			items := make([]T, len(ids))

			for i := 0; nil == err && len(ids) > i; i++ {
				items[i], err = base(ids[i])
			}

			if nil == err {
				return SliceCollection[T](items), nil
			} else {
				return nil, err
			}
		}
	}

	return method, err
}

