package requests

import (
	"fmt"
	"net/http"
	"rent_service/internal/logic/context/providers/provide"
	provide_service "rent_service/internal/logic/services/interfaces/provide"
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
	provide provide.IProvider
	request provide.IRequestProvider
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
	provide provide.IProvider,
	request provide.IRequestProvider,
) server.IController {
	return &controller{providers{provide, request}, authenticator}
}

func CorsFiller(config *cors.Config) {
	config.AddAllowMethods("get", "post", "put")
}

func (self *controller) Register(engine *gin.Engine) {
	engine.GET("/provision-requests", self.get)
	engine.POST("/provision-requests", self.create)
	engine.PUT(fmt.Sprintf("/provision-requests/:%v", PARAM_ID), self.update)
}

func (self *controller) getGetter(ctx *gin.Context) (lister.Method[provide_service.ProvideRequest], error) {
	var err error
	var method lister.Method[provide_service.ProvideRequest]

	service := self.providers.request.GetProvisionRequestService()
	query := ctx.Request.URL.Query()
	user := query.Get(PARAM_USER)
	instance := query[PARAM_INSTANCE]
	pickUpPointId := query.Get(PARAM_PICK_UP_POINT)

	if "" != user && "" == pickUpPointId && nil == instance {
		method, err = lister.ListSingle(user, service.ListProvisionRequstsByUser)
	} else if "" == user && "" != pickUpPointId && nil == instance {
		method, err = lister.ListSingle(pickUpPointId, service.ListProvisionRequstsByPickUpPoint)
	} else if "" == user && "" == pickUpPointId && nil != instance {
		method, err = lister.ListMultiple(instance, service.GetProvisionRequestByInstance)
	} else {
		err = fmt.Errorf(
			"Request must use only one query at a time: '%v', %v, '%v'",
			PARAM_USER, PARAM_INSTANCE, PARAM_PICK_UP_POINT,
		)
	}

	return method, err
}

func (self *controller) get(ctx *gin.Context) {
	var requests collection.Collection[provide_service.ProvideRequest]
	var iter collection.Iterator[provide_service.ProvideRequest]
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
	var form provide_service.RequestCreateForm
	var token token.Token
	var request provide_service.ProvideRequest
	err := ctx.ShouldBindJSON(&form)

	if nil == err {
		token, err = authenticator.ExchangeToken(ctx, self.authenticator)
	}

	if nil == err {
		service := self.providers.request.GetProvisionRequestService()
		request, err = service.CreateProvisionRequest(token, form)
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
	ActionName       string                    `json:"action" binding:"required"`
	Overrides        provide_service.Overrides `json:"overrides"`
	StatePhotos      []uuid.UUID               `json:"state_photos" binding:"required"`
	VerificationCode string                    `json:"verification_code" binding:"required"`
}

func (self *form) Action() string {
	return self.ActionName
}

func getSatisfyAction(
	self *controller,
	requestId uuid.UUID,
	form *form,
) (rqactions.Action, error) {
	cform := provide_service.StartForm{
		RequestId:        requestId,
		VerificationCode: form.VerificationCode,
		Overrides:        form.Overrides,
		TempPhotos:       form.StatePhotos,
	}
	service := self.providers.provide.GetProvisionService()

	return func(token token.Token) error {
		return service.StartProvision(token, cform)
	}, nil
}

func getRejectAction(
	self *controller,
	requestId uuid.UUID,
	_ *form,
) (rqactions.Action, error) {
	service := self.providers.provide.GetProvisionService()

	return func(token token.Token) error {
		return service.RejectProvision(token, requestId)
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

