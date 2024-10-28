package products

import (
	"errors"
	"fmt"
	"net/http"
	product_provider "rent_service/internal/logic/context/providers/product"
	"rent_service/internal/logic/services/errors/cmnerrors"
	product "rent_service/internal/logic/services/interfaces/product"
	"rent_service/internal/misc/types/collection"
	"rent_service/server"
	"rent_service/server/errstructs"
	getter_uuid "rent_service/server/getters/uuid"
	"rent_service/server/pagination"
	"rent_service/server/queryparser"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type providers struct {
	product         product_provider.IProvider
	characteristics product_provider.ICharacteristicsProvider
	photo           product_provider.IPhotoProvider
}

type controller struct {
	providers providers
}

func New(
	product product_provider.IProvider,
	characteristics product_provider.ICharacteristicsProvider,
	photo product_provider.IPhotoProvider,
) server.IController {
	return &controller{providers{product, characteristics, photo}}
}

func CorsFiller(config *cors.Config) {
	config.AddAllowMethods("get")
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

func (self *controller) Register(engine *gin.Engine) {
	engine.GET("/products", self.get)
	engine.GET(fmt.Sprintf("/products/:%v", PARAM_ID), self.getById)
	engine.GET(fmt.Sprintf("/products/:%v/photos", PARAM_ID), self.getPhotos)
	engine.GET(
		fmt.Sprintf("/products/:%v/characteristics", PARAM_ID),
		self.getCharacteristics,
	)
}

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

func (self *controller) get(ctx *gin.Context) {
	var form product.Filter
	var sort product.Sort
	var products collection.Collection[product.Product]
	var iter collection.Iterator[product.Product]
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

	if nil == err {
		service := self.providers.product.GetProductService()
		products, err = service.ListProducts(form, sort)
	}

	if nil == err {
		iter, err = pagination.Apply(ctx, products.Iter())
	}

	if nil == err {
		ctx.JSON(http.StatusOK, collection.Marshaler(iter))
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func (self *controller) getById(ctx *gin.Context) {
	var product product.Product
	id, err := getter_uuid.Parse(ctx.Params.ByName(PARAM_ID))

	if nil == err {
		service := self.providers.product.GetProductService()
		product, err = service.GetProductById(id)
	}

	if nil == err {
		ctx.JSON(http.StatusOK, product)
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusNotFound, errstructs.NewNotFound(cerr))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func (self *controller) getPhotos(ctx *gin.Context) {
	var photos collection.Collection[uuid.UUID]
	id, err := getter_uuid.Parse(ctx.Params.ByName(PARAM_ID))

	if nil == err {
		service := self.providers.photo.GetProductPhotoService()
		photos, err = service.ListProductPhotos(id)
	}

	if nil == err {
		ctx.JSON(http.StatusOK, collection.Marshaler(photos.Iter()))
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusNotFound, errstructs.NewNotFound(cerr))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func (self *controller) getCharacteristics(ctx *gin.Context) {
	var characteristics collection.Collection[product.Charachteristic]
	id, err := getter_uuid.Parse(ctx.Params.ByName(PARAM_ID))

	if nil == err {
		service := self.providers.characteristics.GetProductCharacteristicsService()
		characteristics, err = service.ListProductCharacteristics(id)
	}

	if nil == err {
		ctx.JSON(http.StatusOK, collection.Marshaler(characteristics.Iter()))
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusNotFound, errstructs.NewNotFound(cerr))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

