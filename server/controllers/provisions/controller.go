package provisions

import (
	"errors"
	"fmt"
	"net/http"
	"rent_service/internal/logic/context/providers/provide"
	"rent_service/internal/logic/services/errors/cmnerrors"
	provide_service "rent_service/internal/logic/services/interfaces/provide"
	"rent_service/internal/logic/services/types/token"
	"rent_service/internal/misc/types/collection"
	"rent_service/server"
	"rent_service/server/authenticator"
	"rent_service/server/errstructs"
	"rent_service/server/lister"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type providers struct {
	provide provide.IProvider
}

type controller struct {
	providers     providers
	authenticator authenticator.IAuthenticator
}

const (
	PARAM_INSTANCE string = "instanceId"
	PARAM_USER     string = "userId"
)

func New(
	authenticator authenticator.IAuthenticator,
	provide provide.IProvider,
) server.IController {
	return &controller{providers{provide}, authenticator}
}

func CorsFiller(config *cors.Config) {
	config.AddAllowMethods("get")
}

func (self *controller) Register(engine *gin.Engine) {
	engine.GET("/provisions", self.get)
}

func (self *controller) getGetter(ctx *gin.Context) (lister.Method[provide_service.Provision], error) {
	var err error
	var method lister.Method[provide_service.Provision]

	service := self.providers.provide.GetProvisionService()
	query := ctx.Request.URL.Query()
	user := query.Get(PARAM_USER)
	instance := query[PARAM_INSTANCE]

	if "" != user && nil == instance {
		method, err = lister.ListSingle(user, service.ListProvisionsByUser)
	} else if "" == user && nil != instance {
		method, err = lister.ListMultiple(instance, service.GetProvisionByInstance)
	} else {
		err = fmt.Errorf(
			"Request must use only one query at a time: '%v' or '%v'",
			PARAM_USER, PARAM_INSTANCE,
		)
	}

	return method, err
}

func (self *controller) get(ctx *gin.Context) {
	var provisions collection.Collection[provide_service.Provision]
	var token token.Token
	method, err := self.getGetter(ctx)

	if nil == err {
		token, err = authenticator.ExchangeToken(ctx, self.authenticator)
	}

	if nil == err {
		provisions, err = method(token)
	}

	if nil == err {
		ctx.JSON(
			http.StatusOK,
			collection.Marshaler(provisions.Iter()),
		)
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

