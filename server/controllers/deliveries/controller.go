package deliveries

import (
	"errors"
	"fmt"
	"net/http"
	"rent_service/internal/logic/context/providers/delivery"
	"rent_service/internal/logic/services/errors/cmnerrors"
	delivery_service "rent_service/internal/logic/services/interfaces/delivery"
	"rent_service/internal/logic/services/types/token"
	"rent_service/internal/misc/types/collection"
	"rent_service/server"
	"rent_service/server/authenticator"
	"rent_service/server/errstructs"
	getter_uuid "rent_service/server/getters/uuid"
	"rent_service/server/lister"
	"rent_service/server/pagination"
	"rent_service/server/rqactions"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type providers struct {
	delivery delivery.IProvider
}

type controller struct {
	providers     providers
	authenticator authenticator.IAuthenticator
}

const (
	PARAM_ID            string = "id"
	PARAM_INSTANCE      string = "instanceId"
	PARAM_PICK_UP_POINT string = "pickUpPointId"
)

func New(
	authenticator authenticator.IAuthenticator,
	delivery delivery.IProvider,
) server.IController {
	return &controller{providers{delivery}, authenticator}
}

func (self *controller) Register(engine *gin.Engine) {
	engine.GET("/deliveries", self.get)
	engine.POST("/deliveries", self.create)
	engine.PUT(fmt.Sprintf("/deliveries/:%v", PARAM_ID), self.update)
}

func (self *controller) getGetter(ctx *gin.Context) (lister.Method[delivery_service.Delivery], error) {
	var err error
	var method lister.Method[delivery_service.Delivery]

	service := self.providers.delivery.GetDeliveryService()
	query := ctx.Request.URL.Query()
	instance := query[PARAM_INSTANCE]
	pickUpPointId := query.Get(PARAM_PICK_UP_POINT)

	if "" != pickUpPointId && nil == instance {
		method, err = lister.ListSingle(pickUpPointId, service.ListDeliveriesByPickUpPoint)
	} else if "" == pickUpPointId && nil != instance {
		method, err = lister.ListMultiple(instance, service.GetDeliveryByInstance)
	} else {
		err = fmt.Errorf(
			"Request must use only one query at a time: '%v', %v",
			PARAM_INSTANCE, PARAM_PICK_UP_POINT,
		)
	}

	return method, err
}

func (self *controller) get(ctx *gin.Context) {
	var requests collection.Collection[delivery_service.Delivery]
	var iter collection.Iterator[delivery_service.Delivery]
	var token token.Token
	method, err := self.getGetter(ctx)

	if nil == err {
		token, err = authenticator.ExchangeToken(ctx, self.authenticator)
	}

	if nil == err {
		requests, err = method(token)
	}

	if nil == err {
		iter, err = pagination.Apply(ctx, requests.Iter())
	}

	if nil == err {
		ctx.JSON(http.StatusOK, collection.Marshaler(iter))
	} else if cerr := (cmnerrors.ErrorAuthentication{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusUnauthorized)
	} else if cerr := (cmnerrors.ErrorAuthorization{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusForbidden)
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func (self *controller) create(ctx *gin.Context) {
	var form delivery_service.CreateForm
	var token token.Token
	var request delivery_service.Delivery
	err := ctx.ShouldBindJSON(&form)

	if nil == err {
		token, err = authenticator.ExchangeToken(ctx, self.authenticator)
	}

	if nil == err {
		service := self.providers.delivery.GetDeliveryService()
		request, err = service.CreateDelivery(token, form)
	}

	if nil == err {
		ctx.JSON(http.StatusCreated, request)
	} else if cerr := (cmnerrors.ErrorAuthentication{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusUnauthorized)
	} else if cerr := (cmnerrors.ErrorAuthorization{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusForbidden)
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		if slices.Contains(cerr.What, "instance in pick_up_point") {
			ctx.Status(http.StatusConflict)
		} else {
			ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
		}
		ctx.JSON(http.StatusNotFound, errstructs.NewNotFound(cerr))
	} else if cerr := (cmnerrors.ErrorAlreadyExists{}); errors.As(err, &cerr) {
		if slices.Contains(cerr.What, "delivery") {
			ctx.Status(http.StatusConflict)
		} else {
			ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
		}
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

const (
	ACTION_SEND   string = "send"
	ACTION_ACCEPT string = "accept"
)

type form struct {
	ActionName       string      `json:"action"`
	Description      string      `json:"description"`
	StatePhotos      []uuid.UUID `json:"state_photos"`
	VerificationCode string      `json:"verification_code"`
}

func (self *form) Action() string {
	return self.ActionName
}

func getSendAction(
	self *controller,
	requestId uuid.UUID,
	form *form,
) (rqactions.Action, error) {
	sform := delivery_service.SendForm{
		DeliveryId:       requestId,
		VerificationCode: form.VerificationCode,
		TempPhotos:       form.StatePhotos,
	}
	service := self.providers.delivery.GetDeliveryService()

	return func(token token.Token) error {
		return service.SendDelivery(token, sform)
	}, nil
}

func getAcceptAction(
	self *controller,
	requestId uuid.UUID,
	form *form,
) (rqactions.Action, error) {
	aform := delivery_service.AcceptForm{
		DeliveryId:       requestId,
		VerificationCode: form.VerificationCode,
		TempPhotos:       form.StatePhotos,
	}
	service := self.providers.delivery.GetDeliveryService()

	if "" != form.Description {
		aform.Comment = new(string)
		*aform.Comment = form.Description
	}

	return func(token token.Token) error {
		return service.AcceptDelivery(token, aform)
	}, nil
}

var actions = rqactions.New(
	map[string]rqactions.Getter[controller, *form]{
		ACTION_ACCEPT: getAcceptAction,
		ACTION_SEND:   getSendAction,
	},
)

func (self *controller) update(ctx *gin.Context) {
	var form form
	var action rqactions.Action
	var token token.Token
	requestId, err := getter_uuid.Parse(ctx.Params.ByName(PARAM_ID))

	if nil == err {
		err = ctx.ShouldBindJSON(&form)
	}

	if nil == err {
		action, err = actions.GetAction(self, requestId, &form)
	}

	if nil == err {
		token, err = authenticator.ExchangeToken(ctx, self.authenticator)
	}

	if nil == err {
		err = action(token)
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

