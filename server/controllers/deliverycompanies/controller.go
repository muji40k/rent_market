package deliverycompanies

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
	"rent_service/server/pagination"

	"github.com/gin-gonic/gin"
)

type providers struct {
	company delivery.ICompanyProvider
}

type controller struct {
	providers     providers
	authenticator authenticator.IAuthenticator
}

func New(
	company delivery.ICompanyProvider,
	authenticator authenticator.IAuthenticator,
) server.IController {
	return &controller{providers{company}, authenticator}
}

const PARAM_ID string = "id"

func (self *controller) Register(engine *gin.Engine) {
	engine.GET("/delivery-companies", self.get)
	engine.GET(fmt.Sprintf("/delivery-companies/:%v", PARAM_ID), self.getById)
}

func (self *controller) get(ctx *gin.Context) {
	var companies collection.Collection[delivery_service.DeliveryCompany]
	var iter collection.Iterator[delivery_service.DeliveryCompany]
	token, err := authenticator.ExchangeToken(ctx, self.authenticator)

	if nil == err {
		service := self.providers.company.GetDeliveryCompanyService()
		companies, err = service.ListDeliveryCompanies(token)
	}

	if nil == err {
		iter, err = pagination.Apply(ctx, companies.Iter())
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

func (self *controller) getById(ctx *gin.Context) {
	var token token.Token
	var company delivery_service.DeliveryCompany
	id, err := getter_uuid.Parse(ctx.Params.ByName(PARAM_ID))

	if nil == err {
		token, err = authenticator.ExchangeToken(ctx, self.authenticator)
	}

	if nil == err {
		service := self.providers.company.GetDeliveryCompanyService()
		company, err = service.GetDeliveryCompanyById(token, id)
	}

	if nil == err {
		ctx.JSON(http.StatusOK, company)
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

