package instances

import (
	"errors"
	"fmt"
	"net/http"
	instance_provider "rent_service/internal/logic/context/providers/instance"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/instance"
	"rent_service/internal/logic/services/types/token"
	"rent_service/internal/misc/types/collection"
	"rent_service/server"
	"rent_service/server/authenticator"
	"rent_service/server/errstructs"
	getter_uuid "rent_service/server/getters/uuid"
	"rent_service/server/pagination"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type providers struct {
	instance instance_provider.IProvider
	payPlans instance_provider.IPayPlansProvider
	photo    instance_provider.IPhotoProvider
	review   instance_provider.IReviewProvider
}

type controller struct {
	providers     providers
	authenticator authenticator.IAuthenticator
}

func New(
	authenticator authenticator.IAuthenticator,
	instance instance_provider.IProvider,
	payPlans instance_provider.IPayPlansProvider,
	photo instance_provider.IPhotoProvider,
	review instance_provider.IReviewProvider,
) server.IController {
	return &controller{
		providers{instance, payPlans, photo, review},
		authenticator,
	}
}

func CorsFiller(config *cors.Config) {
	config.AddAllowMethods("get", "put", "post")
}

const (
	PARAM_ID            string = "id"
	PARAM_PRODUCT       string = "productId"
	PARAM_SORT_BY       string = "sortBy"
	PARAM_FILTER_RATING string = "filterRating"
)

func (self *controller) Register(engine *gin.Engine) {
	engine.GET("/instances", self.get)
	engine.GET(fmt.Sprintf("/instances/:%v", PARAM_ID), self.getById)
	engine.PUT(fmt.Sprintf("/instances/:%v", PARAM_ID), self.updateById)
	engine.GET(fmt.Sprintf("/instances/:%v/pay-plans", PARAM_ID), self.getPayPlans)
	engine.PUT(fmt.Sprintf("/instances/:%v/pay-plans", PARAM_ID), self.updatePayPlans)
	engine.GET(fmt.Sprintf("/instances/:%v/photos", PARAM_ID), self.getPhotos)
	engine.POST(fmt.Sprintf("/instances/:%v/photos", PARAM_ID), self.addPhotos)
	engine.GET(fmt.Sprintf("/instances/:%v/reviews", PARAM_ID), self.getReviews)
	engine.POST(fmt.Sprintf("/instances/:%v/reviews", PARAM_ID), self.addReview)
}

func getSortBy(raw string) (instance.Sort, error) {
	switch raw {
	case "":
		return instance.SORT_NONE, nil
	case "ratingAsc":
		return instance.SORT_RATING_ASC, nil
	case "ratingDesc":
		return instance.SORT_RATING_DSC, nil
	case "dateAsc":
		return instance.SORT_DATE_ASC, nil
	case "dateDesc":
		return instance.SORT_DATE_DSC, nil
	case "priceAsc":
		return instance.SORT_PRICE_ASC, nil
	case "priceDesc":
		return instance.SORT_PRICE_DSC, nil
	case "usageAsc":
		return instance.SORT_USAGE_ASC, nil
	case "usageDsc":
		return instance.SORT_USAGE_DSC, nil
	default:
		return instance.SORT_NONE, fmt.Errorf(
			"Unsupported sort by value: '%v'", raw,
		)
	}
}

func (self *controller) get(ctx *gin.Context) {
	var form instance.Filter
	var sortBy instance.Sort
	var instances collection.Collection[instance.Instance]
	var iter collection.Iterator[instance.Instance]
	var err error
	query := ctx.Request.URL.Query()

	form.ProductId, err = getter_uuid.Parse(query.Get(PARAM_PRODUCT))

	if nil == err {
		sortBy, err = getSortBy(query.Get(PARAM_SORT_BY))
	}

	if nil == err {
		service := self.providers.instance.GetInstanceService()
		instances, err = service.ListInstances(form, sortBy)
	}

	if nil == err {
		iter, err = pagination.Apply(ctx, instances.Iter())
	}

	if nil == err {
		ctx.JSON(http.StatusOK, collection.Marshaler(iter))
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusNotFound, errstructs.NewNotFound(cerr))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func (self *controller) getById(ctx *gin.Context) {
	var item instance.Instance
	id, err := getter_uuid.Parse(ctx.Params.ByName(PARAM_ID))

	if nil == err {
		service := self.providers.instance.GetInstanceService()
		item, err = service.GetInstanceById(id)
	}

	if nil == err {
		ctx.JSON(http.StatusOK, item)
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusNotFound, errstructs.NewNotFound(cerr))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func (self *controller) updateById(ctx *gin.Context) {
	var form struct {
		ProductId   uuid.UUID `json:"product" binding:"required"`
		Name        string    `json:"name" binding:"required"`
		Description string    `json:"description" binding:"required"`
		Condition   string    `json:"condition" binding:"required"`
	}

	var token token.Token
	id, err := getter_uuid.Parse(ctx.Params.ByName(PARAM_ID))

	if nil == err {
		err = ctx.ShouldBindJSON(&form)
	}

	if nil == err {
		token, err = authenticator.ExchangeToken(ctx, self.authenticator)
	}

	if nil == err {
		var update = instance.Instance{
			Id:          id,
			ProductId:   form.ProductId,
			Name:        form.Name,
			Description: form.Description,
			Condition:   form.Condition,
		}
		service := self.providers.instance.GetInstanceService()
		err = service.UpdateInstance(token, update)
	}

	if nil == err {
		ctx.Status(http.StatusOK)
	} else if cerr := (cmnerrors.ErrorAuthentication{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusUnauthorized)
	} else if cerr := (cmnerrors.ErrorAuthorization{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusForbidden)
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusNotFound, errstructs.NewNotFound(cerr))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func (self *controller) getPayPlans(ctx *gin.Context) {
	var plans collection.Collection[instance.PayPlan]
	id, err := getter_uuid.Parse(ctx.Params.ByName(PARAM_ID))

	if nil == err {
		service := self.providers.payPlans.GetInstancePayPlansService()
		plans, err = service.GetInstancePayPlans(id)
	}

	if nil == err {
		ctx.JSON(http.StatusOK, collection.Marshaler(plans.Iter()))
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusNotFound, errstructs.NewNotFound(cerr))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func (self *controller) updatePayPlans(ctx *gin.Context) {
	var form instance.PayPlansUpdateForm
	var token token.Token

	id, err := getter_uuid.Parse(ctx.Params.ByName(PARAM_ID))

	if nil == err {
		err = ctx.ShouldBind(&form)
	}

	if nil == err {
		token, err = authenticator.ExchangeToken(ctx, self.authenticator)
	}

	if nil == err {
		service := self.providers.payPlans.GetInstancePayPlansService()
		err = service.UpdateInstancePayPlans(token, id, form)
	}

	if nil == err {
		ctx.Status(http.StatusOK)
	} else if cerr := (cmnerrors.ErrorAuthentication{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusUnauthorized)
	} else if cerr := (cmnerrors.ErrorAuthorization{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusForbidden)
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusNotFound, errstructs.NewNotFound(cerr))
	} else if cerr := (cmnerrors.ErrorConflict{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusConflict)
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func (self *controller) getPhotos(ctx *gin.Context) {
	var photos collection.Collection[uuid.UUID]
	id, err := getter_uuid.Parse(ctx.Params.ByName(PARAM_ID))

	if nil == err {
		service := self.providers.photo.GetInstancePhotoService()
		photos, err = service.ListInstancePhotos(id)
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

func (self *controller) addPhotos(ctx *gin.Context) {
	var photos []uuid.UUID
	var token token.Token

	id, err := getter_uuid.Parse(ctx.Params.ByName(PARAM_ID))

	if nil == err {
		err = ctx.ShouldBind(&photos)
	}

	if nil == err {
		token, err = authenticator.ExchangeToken(ctx, self.authenticator)
	}

	if nil == err {
		service := self.providers.photo.GetInstancePhotoService()
		err = service.AddInstancePhotos(token, id, photos)
	}

	if nil == err {
		ctx.Status(http.StatusCreated)
	} else if cerr := (cmnerrors.ErrorAuthentication{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusUnauthorized)
	} else if cerr := (cmnerrors.ErrorAuthorization{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusForbidden)
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusNotFound, errstructs.NewNotFound(cerr))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func getReviewSortBy(raw string) (instance.ReviewSort, error) {
	switch raw {
	case "":
		return instance.REVIEW_SORT_NONE, nil
	case "dateAsc":
		return instance.REVIEW_SORT_DATE_ASC, nil
	case "dateDesc":
		return instance.REVIEW_SORT_DATE_DSC, nil
	case "ratingAsc":
		return instance.REVIEW_SORT_RATING_ASC, nil
	case "ratingDesc":
		return instance.REVIEW_SORT_RATING_DSC, nil
	default:
		return instance.REVIEW_SORT_NONE, fmt.Errorf(
			"Unsupported sort by value: '%v'", raw,
		)
	}
}

func (self *controller) getReviews(ctx *gin.Context) {
	var reviews collection.Collection[instance.Review]
	var iter collection.Iterator[instance.Review]
	var filter instance.ReviewFilter
	query := ctx.Request.URL.Query()

	sort, err := getReviewSortBy(query.Get(PARAM_SORT_BY))

	if nil == err {
		filter.InstanceId, err = getter_uuid.Parse(ctx.Params.ByName(PARAM_ID))
	}

	if nil == err {
		ratings := query[PARAM_FILTER_RATING]

		if nil != ratings {
			filter.Ratings = make([]instance.ReviewRating, len(ratings))

			for i := 0; nil == err && len(ratings) > i; i++ {
				v, cerr := strconv.ParseUint(ratings[i], 10, 32)

				if nil == cerr {
					filter.Ratings[i] = instance.ReviewRating(v)
				} else {
					err = cerr
				}
			}
		}
	}

	if nil == err {
		service := self.providers.review.GetInstanceReviewService()
		reviews, err = service.ListInstanceReviews(filter, sort)
	}

	if nil == err {
		iter, err = pagination.Apply(ctx, reviews.Iter())
	}

	if nil == err {
		ctx.JSON(http.StatusOK, collection.Marshaler(iter))
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusNotFound, errstructs.NewNotFound(cerr))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func (self *controller) addReview(ctx *gin.Context) {
	var form instance.ReviewPostForm
	var token token.Token

	id, err := getter_uuid.Parse(ctx.Params.ByName(PARAM_ID))

	if nil == err {
		err = ctx.ShouldBind(&form)
	}

	if nil == err {
		token, err = authenticator.ExchangeToken(ctx, self.authenticator)
	}

	if nil == err {
		service := self.providers.review.GetInstanceReviewService()
		err = service.PostInstanceReview(token, id, form)
	}

	if nil == err {
		ctx.Status(http.StatusCreated)
	} else if cerr := (cmnerrors.ErrorAuthentication{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusUnauthorized)
	} else if cerr := (cmnerrors.ErrorAuthorization{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusForbidden)
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusNotFound, errstructs.NewNotFound(cerr))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

