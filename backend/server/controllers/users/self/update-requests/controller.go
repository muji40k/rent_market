package updaterequests

import (
	"errors"
	"fmt"
	"net/http"
	"rent_service/internal/logic/context/providers/user"
	"rent_service/internal/logic/services/errors/cmnerrors"
	service_user "rent_service/internal/logic/services/interfaces/user"
	"rent_service/internal/logic/services/types/token"
	"rent_service/server"
	"rent_service/server/authenticator"
	"rent_service/server/errstructs"
	"rent_service/server/getters/uuid"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type providers struct {
	password user.IPasswordUpdateProvider
}

type controller struct {
	providers     providers
	authenticator authenticator.IAuthenticator
}

func New(
	authenticator authenticator.IAuthenticator,
	password user.IPasswordUpdateProvider,
) server.IController {
	return &controller{providers{password}, authenticator}
}

func CorsFiller(config *cors.Config) {
	config.AddAllowMethods("post")
}

const PARAM_REQUEST_ID string = "request_id"

func (self *controller) Register(engine *gin.Engine) {
	engine.POST("/users/self/update-requests/passwords", self.createUpdateRequest)
	engine.DELETE(
		fmt.Sprintf("/users/self/update-requests/passwords/:%v", PARAM_REQUEST_ID),
		self.authenticateRequest,
	)
}

func (self *controller) createUpdateRequest(ctx *gin.Context) {
	var token token.Token
	var result service_user.PasswordUpdateRequest
	var form struct {
		Old string `json:"old_password" binding:"required"`
		New string `json:"new_password" binding:"required"`
	}

	err := ctx.ShouldBindJSON(&form)

	if nil == err {
		token, err = authenticator.ExchangeToken(ctx, self.authenticator)
	}

	if nil == err {
		service := self.providers.password.GetUserPasswordUpdateService()
		result, err = service.RequestPasswordUpdate(token, form.Old, form.New)
	}

	if nil == err {
		ctx.JSON(http.StatusOK, result)
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

func (self *controller) authenticateRequest(ctx *gin.Context) {
	var token token.Token
	var form struct {
		Code string `json:"code" binding:"required"`
	}

	requestId, err := uuid.Parse(ctx.Params.ByName(PARAM_REQUEST_ID))

	if nil == err {
		err = ctx.ShouldBindJSON(&form)
	}

	if nil == err {
		token, err = authenticator.ExchangeToken(ctx, self.authenticator)
	}

	if nil == err {
		service := self.providers.password.GetUserPasswordUpdateService()
		err = service.AuthenticatePasswordUpdateRequest(
			token,
			requestId,
			form.Code,
		)
	}

	if nil == err {
		ctx.Status(http.StatusOK)
	} else if cerr := (cmnerrors.ErrorAuthentication{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusUnauthorized)
	} else if cerr := (cmnerrors.ErrorAuthorization{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusForbidden)
	} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusNotFound, errstructs.NewNotFound(cerr))
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

