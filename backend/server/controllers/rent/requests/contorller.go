package requests

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
	rent    rent.IProvider
	request rent.IRequestProvider
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
	request rent.IRequestProvider,
) server.IController {
	return &controller{providers{rent, request}, authenticator}
}

func CorsFiller(config *cors.Config) {
	config.AddAllowMethods("get", "post")
}

func (self *controller) Register(engine *gin.Engine) {
	engine.GET("/rent-requests", self.get)
	engine.POST("/rent-requests", self.create)
	engine.PUT(fmt.Sprintf("/rent-requests/:%v", PARAM_ID), self.update)
}

func (self *controller) getGetter(ctx *gin.Context) (lister.Method[rent_service.RentRequest], error) {
	var err error
	var method lister.Method[rent_service.RentRequest]

	service := self.providers.request.GetRentRequestService()
	query := ctx.Request.URL.Query()
	user := query.Get(PARAM_USER)
	instance := query[PARAM_INSTANCE]
	pickUpPointId := query.Get(PARAM_PICK_UP_POINT)

	if "" != user && "" == pickUpPointId && nil == instance {
		method, err = lister.ListSingle(user, service.ListRentRequstsByUser)
	} else if "" == user && "" != pickUpPointId && nil == instance {
		method, err = lister.ListSingle(pickUpPointId, service.ListRentRequstsByPickUpPoint)
	} else if "" == user && "" == pickUpPointId && nil != instance {
		method, err = lister.ListMultiple(instance, service.GetRentRequestByInstance)
	} else {
		err = fmt.Errorf(
			"Request must use only one query at a time: '%v', %v, '%v'",
			PARAM_USER, PARAM_INSTANCE, PARAM_PICK_UP_POINT,
		)
	}

	return method, err
}

func (self *controller) get(ctx *gin.Context) {
	var requests collection.Collection[rent_service.RentRequest]
	var iter collection.Iterator[rent_service.RentRequest]
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
	var form rent_service.RequestCreateForm
	var token token.Token
	var request rent_service.RentRequest
	err := ctx.ShouldBindJSON(&form)

	if nil == err {
		token, err = authenticator.ExchangeToken(ctx, self.authenticator)
	}

	if nil == err {
		service := self.providers.request.GetRentRequestService()
		request, err = service.CreateRentRequest(token, form)
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
	ACTION_REJECT  string = "reject"
)

type form struct {
	ActionName       string      `json:"action" binding:"required"`
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
	cform := rent_service.StartForm{
		RequestId:        requestId,
		VerificationCode: form.VerificationCode,
		TempPhotos:       form.StatePhotos,
	}
	service := self.providers.rent.GetRentService()

	return func(token token.Token) error {
		return service.StartRent(token, cform)
	}, nil
}

func getRejectAction(
	self *controller,
	requestId uuid.UUID,
	_ *form,
) (rqactions.Action, error) {
	service := self.providers.rent.GetRentService()

	return func(token token.Token) error {
		return service.RejectRent(token, requestId)
	}, nil
}

var actions = rqactions.New(
	map[string]rqactions.Getter[controller, *form]{
		ACTION_SATISFY: getSatisfyAction,
		ACTION_REJECT:  getRejectAction,
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

