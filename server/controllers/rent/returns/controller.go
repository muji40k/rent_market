package returns

import (
	"errors"
	"fmt"
	"net/http"
	"rent_service/internal/logic/context/providers/rent"
	"rent_service/internal/logic/services/errors/cmnerrors"
	rent_service "rent_service/internal/logic/services/interfaces/rent"
	"rent_service/internal/logic/services/types/token"
	"rent_service/internal/misc/types/collection"
	"rent_service/server"
	"rent_service/server/authenticator"
	"rent_service/server/errstructs"
	getter_uuid "rent_service/server/getters/uuid"
	"rent_service/server/lister"
	"rent_service/server/pagination"
	"rent_service/server/rqactions"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type providers struct {
	rent rent.IProvider
	ret  rent.IReturnProvider
}

type controller struct {
	providers     providers
	authenticator authenticator.IAuthenticator
}

const (
	PARAM_ID            string = "id"
	PARAM_INSTANCE      string = "instanceId"
	PARAM_USER          string = "userId"
	PARAM_PICK_UP_POINT string = "pickUpPointId"
)

func New(
	authenticator authenticator.IAuthenticator,
	rent rent.IProvider,
	ret rent.IReturnProvider,
) server.IController {
	return &controller{providers{rent, ret}, authenticator}
}

func (self *controller) Register(engine *gin.Engine) {
	engine.GET("/rent-returns", self.get)
	engine.POST("/rent-returns", self.create)
	engine.PUT(fmt.Sprintf("/rent-returns/:%v", PARAM_ID), self.update)
	engine.DELETE(fmt.Sprintf("/rent-returns/:%v", PARAM_ID), self.delete)
}

func (self *controller) getGetter(ctx *gin.Context) (lister.Method[rent_service.ReturnRequest], error) {
	var err error
	var method lister.Method[rent_service.ReturnRequest]

	service := self.providers.ret.GetRentReturnService()
	query := ctx.Request.URL.Query()
	user := query.Get(PARAM_USER)
	instance := query[PARAM_INSTANCE]
	pickUpPointId := query.Get(PARAM_PICK_UP_POINT)

	if "" != user && "" == pickUpPointId && nil == instance {
		method, err = lister.ListSingle(user, service.ListRentReturnsByUser)
	} else if "" == user && "" != pickUpPointId && nil == instance {
		method, err = lister.ListSingle(pickUpPointId, service.ListRentReturnsByPickUpPoint)
	} else if "" == user && "" == pickUpPointId && nil != instance {
		method, err = lister.ListMultiple(instance, service.GetRentReturnByInstance)
	} else {
		err = fmt.Errorf(
			"Request must use only one query at a time: '%v', %v, '%v'",
			PARAM_USER, PARAM_INSTANCE, PARAM_PICK_UP_POINT,
		)
	}

	return method, err
}

func (self *controller) get(ctx *gin.Context) {
	var requests collection.Collection[rent_service.ReturnRequest]
	var iter collection.Iterator[rent_service.ReturnRequest]
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
	var form rent_service.ReturnCreateForm
	var token token.Token
	var request rent_service.ReturnRequest
	err := ctx.ShouldBindJSON(&form)

	if nil == err {
		token, err = authenticator.ExchangeToken(ctx, self.authenticator)
	}

	if nil == err {
		service := self.providers.ret.GetRentReturnService()
		request, err = service.CreateRentReturn(token, form)
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
		ctx.JSON(http.StatusNotFound, errstructs.NewNotFound(cerr))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

const (
	ACTION_SATISFY string = "satisfy"
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

func getSatisfyAction(
	self *controller,
	requestId uuid.UUID,
	form *form,
) (rqactions.Action, error) {
	sform := rent_service.StopForm{
		ReturnId:         requestId,
		VerificationCode: form.VerificationCode,
		TempPhotos:       form.StatePhotos,
	}

	if "" != form.Description {
		sform.Comment = new(string)
		*sform.Comment = form.Description
	}

	service := self.providers.rent.GetRentService()

	return func(token token.Token) error {
		return service.StopRent(token, sform)
	}, nil
}

var actions = rqactions.New(
	map[string]rqactions.Getter[controller, *form]{
		ACTION_SATISFY: getSatisfyAction,
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

func (self *controller) delete(ctx *gin.Context) {
	var token token.Token
	requestId, err := getter_uuid.Parse(ctx.Params.ByName(PARAM_ID))

	if nil == err {
		token, err = authenticator.ExchangeToken(ctx, self.authenticator)
	}

	if nil == err {
		service := self.providers.ret.GetRentReturnService()
		err = service.CancelRentReturn(token, requestId)
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

