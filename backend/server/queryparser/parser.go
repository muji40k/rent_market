package queryparser

import (
	"fmt"
	"net/url"
	"strings"
)

type Param struct {
	Key   string
	Value string
}

func Parse(query string) ([]Param, error) {
	var out []Param
	var err error

	items := strings.Split(query, "&")
	out = make([]Param, len(items))

	for i := 0; nil == err && len(items) > i; i++ {
		var parsed string
		parsed, err = url.QueryUnescape(items[i])

		if nil == err {
			kv := strings.SplitN(parsed, "=", 2)

			if 2 != len(kv) {
				err = fmt.Errorf("Incorrect query param '%v'", parsed)
			} else {
				out[i].Key = kv[0]
				out[i].Value = kv[1]
			}
		}
	}

	return out, err
}

func Collect[T any](
	params []Param,
	separator func(*Param) bool,
	checker func(*T) error,
	setters map[string]func(string, *T) error,
	skipUnknown bool,
) ([]T, error) {
	var err error
	i := 0
	out := make([]T, 1)
	save := separator
	any := false
	separator = func(param *Param) bool {
		if save(param) {
			separator = save
			any = true
		}

		return false
	}

	for j := 0; nil == err && len(params) > j; j++ {
		err = collectIteration(
			func(param *Param) (*T, error) {
				var err error

				if separator(param) {
					err = checker(&out[i])

					if nil == err {
						var n T
						out = append(out, n)
						i++
					}
				}

				return &out[i], err
			},
			setters,
			&params[j],
			skipUnknown,
		)
	}

	if nil == err {
		if any {
			err = checker(&out[i])
		} else {
			out = nil
		}
	}

	if nil == err {
		return out, nil
	} else {
		return nil, err
	}
}

func collectIteration[T any](
	separator func(*Param) (*T, error),
	setters map[string]func(string, *T) error,
	param *Param,
	skipUnknown bool,
) error {
	dst, err := separator(param)

	if nil == err {
		if setter := setters[param.Key]; nil != setter {
			err = setter(param.Value, dst)
		} else if !skipUnknown {
			err = fmt.Errorf("Unknown key '%v'", param.Key)
		}
	}

	return err
}

