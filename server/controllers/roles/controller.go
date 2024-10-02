package roles

import (
	"errors"
	"net/http"
	"rent_service/internal/logic/context/providers/user"
	"rent_service/internal/logic/services/errors/cmnerrors"
	service_user "rent_service/internal/logic/services/interfaces/user"
	"rent_service/server"
	"rent_service/server/authenticator"
	"rent_service/server/errstructs"

	"github.com/gin-gonic/gin"
)

type providers struct {
	role user.IRoleProvider
}

type controller struct {
	providers     providers
	authenticator authenticator.IAuthenticator
}

func New(
	authenticator authenticator.IAuthenticator,
	role user.IRoleProvider,
) server.IController {
	return &controller{providers{role}, authenticator}
}

func (self *controller) Register(engine *gin.Engine) {
	engine.GET("/administrators/self", self.isAdministrator)
	engine.GET("/renters/self", self.isRenter)
	engine.POST("/renters/self", self.registerRenter)
	engine.GET("/storekeeper/self", self.isStorekeeper)
}

func (self *controller) isAdministrator(ctx *gin.Context) {
	token, err := authenticator.ExchangeToken(ctx, self.authenticator)

	if nil == err {
		service := self.providers.role.GetRoleService()
		err = service.IsAdministrator(token)
	}

	if nil == err {
		ctx.Status(http.StatusOK)
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

func (self *controller) isRenter(ctx *gin.Context) {
	token, err := authenticator.ExchangeToken(ctx, self.authenticator)

	if nil == err {
		service := self.providers.role.GetRoleService()
		err = service.IsRenter(token)
	}

	if nil == err {
		ctx.Status(http.StatusOK)
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

func (self *controller) registerRenter(ctx *gin.Context) {
	token, err := authenticator.ExchangeToken(ctx, self.authenticator)

	if nil == err {
		service := self.providers.role.GetRoleService()
		err = service.RegisterAsRenter(token)
	}

	if nil == err {
		ctx.Status(http.StatusCreated)
	} else if cerr := (service_user.ErrorAlreadyRenter{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusOK)
	} else if cerr := (cmnerrors.ErrorAuthentication{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusUnauthorized)
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func (self *controller) isStorekeeper(ctx *gin.Context) {
	var sk service_user.StoreKeeper
	token, err := authenticator.ExchangeToken(ctx, self.authenticator)

	if nil == err {
		service := self.providers.role.GetRoleService()
		sk, err = service.IsStoreKeeper(token)
	}

	if nil == err {
		ctx.JSON(http.StatusOK, sk)
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

