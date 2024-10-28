package dummy

import (
	// "errors"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"rent_service/internal/logic/services/interfaces/product"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/logic/services/types/daytime"
	"rent_service/internal/misc/types/collection"
	"rent_service/server"

	"rent_service/server/errstructs"
	getter_uuid "rent_service/server/getters/uuid"
	"rent_service/server/pagination"
	"rent_service/server/queryparser"

	"github.com/gin-gonic/gin"
)

type controller struct{}

type handle struct {
	Id    int           `json:"id"`
	Value string        `json:"value"`
	Date  *date.Date    `json:"date,omitempty"`
	Time  *daytime.Time `json:"time,omitempty"`
}

var values []handle

func New() server.IController {
	if nil == values {
		values = make([]handle, 26)

		for i := 0; 26 > i; i++ {
			values[i] = handle{i, fmt.Sprintf("text%v", i), nil, nil}

			if 0 == i%2 {
				values[i].Date = new(date.Date)
				*values[i].Date = date.Date{Time: time.Now()}
			} else {
				values[i].Time = new(daytime.Time)
				*values[i].Time = daytime.Time{Time: time.Now()}
			}
		}
	}

	return &controller{}
}

func (self *controller) Register(engine *gin.Engine) {
	engine.GET("/dummy", self.get)
	engine.GET("/dummy/:id", self.getId)
	engine.GET("/dummy/monster", self.parse)
	engine.GET("/dummy/check", self.check)
}

func (self *controller) get(ctx *gin.Context) {
	col := collection.SliceCollection(values)
	iter, err := pagination.Apply(ctx, col.Iter())

	if nil == err {
		ctx.JSON(http.StatusOK, collection.Marshaler(iter))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func (self *controller) getId(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Params.ByName("id"))

	if nil != err {
		ctx.String(http.StatusBadRequest, err.Error())
	} else if 0 > id || len(values) <= id {
		ctx.Status(http.StatusNotFound)
	} else {
		ctx.JSON(http.StatusOK, values[id])
	}
}

type Person struct {
	Name   string `json:"name"`
	Values struct {
		Num   int      `json:"number"`
		Str   string   `json:"quote"`
		Words []string `json:"words"`
	} `json:"values"`
}

func CollectPerson(params []queryparser.Param) ([]Person, error) {
	return queryparser.Collect(
		params,
		func(param *queryparser.Param) bool { return "name" == param.Key },
		func(person *Person) error {
			if person.Name == "" {
				return errors.New("No name...")
			}

			return nil
		},
		map[string]func(string, *Person) error{
			"name": func(value string, p *Person) error {
				p.Name = value
				return nil
			},
			"values[num]": func(value string, p *Person) error {
				var err error
				p.Values.Num, err = strconv.Atoi(value)
				return err
			},
			"values[str]": func(value string, p *Person) error {
				p.Values.Str = value
				return nil
			},
			"values[say]": func(value string, p *Person) error {
				p.Values.Words = append(p.Values.Words, value)
				return nil
			},
		},
		true,
	)
}

func (self *controller) parse(ctx *gin.Context) {
	var person []Person
	var iter collection.Iterator[Person]
	items, err := queryparser.Parse(ctx.Request.URL.RawQuery)

	if nil == err {
		person, err = CollectPerson(items)
	}

	if nil == err {
		iter, err = pagination.Apply(
			ctx,
			collection.SliceCollection(person).Iter(),
		)
	}

	if nil == err {
		ctx.JSON(http.StatusOK, collection.Marshaler(iter))
	} else {
		ctx.String(http.StatusBadRequest, err.Error())
	}
}

const (
	PARAM_ID                        string = "id"
	PARAM_SORT_BY                   string = "sortBy"
	PARAM_SEARCH                    string = "search"
	PARAM_CATEGORY                  string = "category"
	PARAM_CHARACTERISTICS_KEY       string = "characteristics[key]"
	PARAM_CHARACTERISTICS_VALUES    string = "characteristics[values]"
	PARAM_CHARACTERISTICS_RANGE_MIN string = "characteristics[range][min]"
	PARAM_CHARACTERISTICS_RANGE_MAX string = "characteristics[range][max]"
)

func getSortBy(raw string) (product.Sort, error) {
	switch raw {
	case "":
		return product.SORT_NONE, nil
	case "offersAsc":
		return product.SORT_OFFERS_ASC, nil
	case "offersDesc":
		return product.SORT_OFFERS_DSC, nil
	default:
		return product.SORT_NONE, fmt.Errorf("Unknown sort by value '%v'", raw)
	}
}

func next(param *queryparser.Param) bool {
	return PARAM_CHARACTERISTICS_KEY == param.Key
}

func checker(char *product.FilterCharachteristic) error {
	if char.Key == "" {
		return errors.New("Characteristic key wasn't provided")
	}

	if nil == char.Range && nil == char.Values ||
		nil != char.Range && nil != char.Values {
		return errors.New(
			"Only one type for characteristic is allowed",
		)
	}

	return nil
}

var setters = map[string]func(string, *product.FilterCharachteristic) error{
	PARAM_CHARACTERISTICS_KEY: func(
		value string,
		char *product.FilterCharachteristic,
	) error {
		char.Key = value
		return nil
	},
	PARAM_CHARACTERISTICS_VALUES: func(
		value string,
		char *product.FilterCharachteristic,
	) error {
		char.Values = append(char.Values, value)
		return nil
	},
	PARAM_CHARACTERISTICS_RANGE_MIN: func(
		value string,
		char *product.FilterCharachteristic,
	) error {
		if nil == char.Range {
			char.Range = new(struct {
				Min float64
				Max float64
			})
		}

		var err error
		char.Range.Min, err = strconv.ParseFloat(value, 64)

		return err
	},
	PARAM_CHARACTERISTICS_RANGE_MAX: func(
		value string,
		char *product.FilterCharachteristic,
	) error {
		if nil == char.Range {
			char.Range = new(struct {
				Min float64
				Max float64
			})
		}

		var err error
		char.Range.Max, err = strconv.ParseFloat(value, 64)

		return err
	},
}

func collectCharacteristics(
	query string,
) ([]product.FilterCharachteristic, error) {
	var out []product.FilterCharachteristic
	params, err := queryparser.Parse(query)

	if nil == err {
		out, err = queryparser.Collect(params, next, checker, setters, true)
	}

	return out, err
}

func (self *controller) check(ctx *gin.Context) {
	var form product.Filter
	var sort product.Sort
	var err error

	query := ctx.Request.URL.Query()
	form.CategoryId, err = getter_uuid.Parse(query.Get(PARAM_CATEGORY))

	if nil == err && "" != query.Get(PARAM_SEARCH) {
		form.Query = new(string)
		*form.Query = query.Get(PARAM_SEARCH)
	}

	if nil == err {
		sort, err = getSortBy(query.Get(PARAM_SORT_BY))
	}

	if nil == err {
		form.Characteristics, err = collectCharacteristics(ctx.Request.URL.RawQuery)
	}

	fmt.Println(form)
	fmt.Println(sort)

	if nil == err {
		ctx.Status(http.StatusOK)
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

