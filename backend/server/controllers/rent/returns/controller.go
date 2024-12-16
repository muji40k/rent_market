package returns

import (
	"fmt"
	"net/http"
	"rent_service/internal/logic/context/providers/rent"
	rent_service "rent_service/internal/logic/services/interfaces/rent"
	"rent_service/internal/logic/services/types/token"
	"rent_service/internal/misc/types/collection"
	"rent_service/server"
	"rent_service/server/authenticator"
	"rent_service/server/errmap"
	getter_uuid "rent_service/server/getters/uuid"
	"rent_service/server/lister"
	"rent_service/server/pagination"
	"rent_service/server/rqactions"

	"github.com/gin-contrib/cors"
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

func CorsFiller(config *cors.Config) {
	config.AddAllowMethods("get", "post", "put", "delete")
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

	errmap.MapValue(ctx,
		func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, collection.Marshaler(iter))
		},
		err,
	)
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

	errmap.MapValue(ctx,
		func(ctx *gin.Context) {
			ctx.JSON(http.StatusCreated, request)
		},
		err,
	)
}

const (
	ACTION_SATISFY string = "satisfy"
)

type form struct {
	ActionName       string      `json:"action" binding:"required"`
	Description      string      `json:"description"`
	StatePhotos      []uuid.UUID `json:"state_photos" binding:"required"`
	VerificationCode string      `json:"verification_code" binding:"required"`
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

	errmap.Map(ctx, http.StatusOK, err)
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

	errmap.Map(ctx, http.StatusOK, err)
}

