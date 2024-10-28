package pagination

import (
	"rent_service/internal/misc/types/collection"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	PARAM_OFFSET string = "offset"
	PARAM_SIZE   string = "size"
)

func parse(value string) (uint, error) {
	if out, err := strconv.ParseUint(value, 10, 32); nil == err {
		return uint(out), nil
	} else {
		return 0, err
	}
}

func Apply[T any](
	ctx *gin.Context,
	iter collection.Iterator[T],
) (collection.Iterator[T], error) {
	var err error
	query := ctx.Request.URL.Query()
	offsetRaw := query.Get(PARAM_OFFSET)
	sizeRaw := query.Get(PARAM_SIZE)

	if "" != offsetRaw {
		if offset, cerr := parse(offsetRaw); nil == cerr && 0 != offset {
			iter = collection.SkipIterator(offset, iter)
		} else {
			err = cerr
		}
	}

	if nil == err && "" != sizeRaw {
		if size, cerr := parse(sizeRaw); nil == cerr {
			iter = collection.LimitIterator(size, iter)
		} else {
			err = cerr
		}
	}

	return iter, err
}

