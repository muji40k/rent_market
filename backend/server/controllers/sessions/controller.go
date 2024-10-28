package sessions

import (
	"errors"
	"net/http"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/server"
	"rent_service/server/authenticator"
	"rent_service/server/errstructs"
	"rent_service/server/headers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type controller struct {
	authenticator authenticator.IAuthenticator
}

func New(authenticator authenticator.IAuthenticator) server.IController {
	return &controller{authenticator}
}

func CorsFiller(config *cors.Config) {
	config.AddAllowMethods("post", "put", "delete")
}

func (self *controller) Register(engine *gin.Engine) {
	engine.POST("/sessions", self.login)
	engine.PUT("/sessions", self.renew)
	engine.DELETE("/sessions", self.logout)
}

func (self *controller) login(ctx *gin.Context) {
	var form struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var token authenticator.ApiToken

	err := ctx.ShouldBindJSON(&form)

	if nil == err {
		token, err = self.authenticator.Login(form.Email, form.Password)
	}

	if nil == err {
		ctx.JSON(http.StatusOK, token)
	} else if cerr := (cmnerrors.ErrorAuthentication{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusUnauthorized)
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(
			http.StatusInternalServerError,
			errstructs.NewInternalErr(err),
		)
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func (self *controller) renew(ctx *gin.Context) {
	var err error
	var token authenticator.ApiToken

	token.Access = ctx.Request.Header.Get(headers.API_KEY)
	token.Renew = ctx.Request.Header.Get(headers.API_RENEW)

	token, err = self.authenticator.RenewKey(token)

	if nil == err {
		ctx.JSON(http.StatusOK, token)
	} else if cerr := (cmnerrors.ErrorAuthentication{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusUnauthorized)
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(
			http.StatusInternalServerError,
			errstructs.NewInternalErr(err),
		)
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func (self *controller) logout(ctx *gin.Context) {
	access := ctx.Request.Header.Get(headers.API_KEY)

	err := self.authenticator.Logout(access)

	if nil == err {
		ctx.Status(http.StatusOK)
	} else if cerr := (cmnerrors.ErrorAuthentication{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusUnauthorized)
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(
			http.StatusInternalServerError,
			errstructs.NewInternalErr(err),
		)
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

