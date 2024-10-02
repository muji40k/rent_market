package storedinstances

import (
	"errors"
	"fmt"
	"net/http"
	"rent_service/internal/logic/context/providers/storage"
	"rent_service/internal/logic/services/errors/cmnerrors"
	storage_service "rent_service/internal/logic/services/interfaces/storage"
	"rent_service/internal/logic/services/types/token"
	"rent_service/internal/misc/types/collection"
	"rent_service/server"
	"rent_service/server/authenticator"
	"rent_service/server/errstructs"
	"rent_service/server/lister"
	"rent_service/server/pagination"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type providers struct {
	storage storage.IProvider
}

type controller struct {
	providers     providers
	authenticator authenticator.IAuthenticator
}

const (
	PARAM_ID            string = "id"
	PARAM_INSTANCE      string = "instanceId"
	PARAM_PICK_UP_POINT string = "pickUpPointId"
)

func New(
	authenticator authenticator.IAuthenticator,
	storage storage.IProvider,
) server.IController {
	return &controller{providers{storage}, authenticator}
}

func (self *controller) Register(engine *gin.Engine) {
	engine.GET("/stored-instances", self.get)
}

func (self *controller) getGetter(ctx *gin.Context) (lister.MethodNA[storage_service.Storage], error) {
	var err error
	var method lister.MethodNA[storage_service.Storage]

	service := self.providers.storage.GetStorageService()
	instance := ctx.Request.URL.Query()[PARAM_INSTANCE]
	pickUpPointId := ctx.Request.URL.Query().Get(PARAM_PICK_UP_POINT)

	if "" != pickUpPointId && nil == instance {
		var token token.Token
		token, err = authenticator.ExchangeToken(ctx, self.authenticator)

		if nil == err {
			method, err = lister.ListSingleNA(
				pickUpPointId,
				func(id uuid.UUID) (collection.Collection[storage_service.Storage], error) {
					return service.ListStoragesByPickUpPoint(token, id)
				},
			)
		}
	} else if "" == pickUpPointId && nil != instance {
		method, err = lister.ListMultipleNA(instance, service.GetStorageByInstance)
	} else {
		err = fmt.Errorf(
			"Request must use only one query at a time: '%v', %v",
			PARAM_INSTANCE, PARAM_PICK_UP_POINT,
		)
	}

	return method, err
}

func (self *controller) get(ctx *gin.Context) {
	var storages collection.Collection[storage_service.Storage]
	var iter collection.Iterator[storage_service.Storage]
	method, err := self.getGetter(ctx)

	if nil == err {
		storages, err = method()
	}

	if nil == err {
		iter, err = pagination.Apply(ctx, storages.Iter())
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

