package users

import (
	"errors"
	"fmt"
	"net/http"
	"rent_service/internal/logic/context/providers/login"
	"rent_service/internal/logic/context/providers/user"
	"rent_service/internal/logic/services/errors/cmnerrors"
	service_user "rent_service/internal/logic/services/interfaces/user"
	"rent_service/internal/logic/services/types/token"
	"rent_service/server"
	"rent_service/server/authenticator"
	"rent_service/server/errstructs"
	"slices"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type providers struct {
	user  user.IProvider
	login login.IProvider
}

type controller struct {
	providers     providers
	authenticator authenticator.IAuthenticator
}

func New(
	authenticator authenticator.IAuthenticator,
	user user.IProvider,
	login login.IProvider,
) server.IController {
	return &controller{providers{user, login}, authenticator}
}

func CorsFiller(config *cors.Config) {
	config.AddAllowMethods("get", "post", "put")
}

func (self *controller) Register(engine *gin.Engine) {
	engine.POST("/users", self.register)
	engine.GET("/users/self", self.get)
	engine.PUT("/users/self", self.update)
}

func (self *controller) register(ctx *gin.Context) {
	var form struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}

	err := ctx.ShouldBindJSON(&form)

	if nil == err {
		service := self.providers.login.GetLoginService()
		err = service.Register(form.Email, form.Password, form.Name)
	}

	if nil == err {
		ctx.Status(http.StatusCreated)
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else if cerr := (cmnerrors.ErrorAlreadyExists{}); errors.As(err, &cerr) {
		if slices.Contains(cerr.What, "email") {
			ctx.Status(http.StatusUnauthorized)
		} else {
			ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(
				fmt.Errorf("Caught unsupported already exists error: %w", err),
			))
		}
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func (self *controller) get(ctx *gin.Context) {
	token, err := authenticator.ExchangeToken(ctx, self.authenticator)

	var info service_user.Info
	if nil == err {
		service := self.providers.user.GetUserService()
		info, err = service.GetSelfUserInfo(token)
	}

	if nil == err {
		ctx.JSON(http.StatusOK, info)
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

func (self *controller) update(ctx *gin.Context) {
	var token token.Token
	var form struct {
		Data     *service_user.UpdateForm `json:"data"`
		Password *struct {
			Old string `json:"old" binding:"required"`
			New string `json:"new" binding:"required"`
		} `json:"password"`
	}

	err := ctx.ShouldBindJSON(&form)

	if nil == err &&
		((nil == form.Data && nil == form.Password) ||
			(nil != form.Data && nil != form.Password)) {
		err = errors.New("User can update only one thing at at time")
	}

	if nil == err {
		token, err = authenticator.ExchangeToken(ctx, self.authenticator)
	}

	if nil == err && nil != form.Data {
		service := self.providers.user.GetUserService()
		err = service.UpdateSelfUserInfo(token, *form.Data)
	}

	if nil == err && nil != form.Password {
		service := self.providers.user.GetUserService()
		err = service.UpdateSelfUserPassword(
			token, form.Password.Old, form.Password.New,
		)
	}

	if nil == err {
		ctx.Status(http.StatusOK)
	} else if cerr := (cmnerrors.ErrorAuthentication{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusUnauthorized)
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

