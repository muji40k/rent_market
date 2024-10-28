package payments

import (
	"errors"
	"fmt"
	"net/http"
	"rent_service/internal/logic/context/providers/payment"
	"rent_service/internal/logic/services/errors/cmnerrors"
	service_user "rent_service/internal/logic/services/interfaces/payment"
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
	payment payment.IRentPaymentProvider
}

type controller struct {
	providers     providers
	authenticator authenticator.IAuthenticator
}

const (
	PARAM_INSTANCE string = "instanceId"
	PARAM_RENT     string = "rentId"
)

func New(
	authenticator authenticator.IAuthenticator,
	payment payment.IRentPaymentProvider,
) server.IController {
	return &controller{providers{payment}, authenticator}
}

func CorsFiller(config *cors.Config) {
	config.AddAllowMethods("get")
}

func (self *controller) Register(engine *gin.Engine) {
	engine.GET("/payments", self.get)
}

func (self *controller) getGetter(ctx *gin.Context) (lister.Method[service_user.Payment], error) {
	var err error
	var method lister.Method[service_user.Payment]

	service := self.providers.payment.GetRentPaymentService()
	query := ctx.Request.URL.Query()
	instance := query.Get(PARAM_INSTANCE)
	rent := query.Get(PARAM_RENT)

	if "" != instance && "" == rent {
		method, err = lister.ListSingle(instance, service.GetPaymentsByInstance)
	} else if "" == instance && "" != rent {
		method, err = lister.ListSingle(rent, service.GetPaymentsByRentId)
	} else {
		err = fmt.Errorf(
			"Request must use only one query at a time: '%v' or '%v'",
			PARAM_RENT, PARAM_INSTANCE,
		)
	}

	return method, err
}

func (self *controller) get(ctx *gin.Context) {
	var payments collection.Collection[service_user.Payment]
	var token token.Token
	method, err := self.getGetter(ctx)

	if nil == err {
		token, err = authenticator.ExchangeToken(ctx, self.authenticator)
	}

	if nil == err {
		payments, err = method(token)
	}

	if nil == err {
		ctx.JSON(http.StatusOK, collection.Marshaler(payments.Iter()))
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

