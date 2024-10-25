package photos

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"rent_service/internal/logic/context/providers/photo"
	"rent_service/internal/logic/services/errors/cmnerrors"
	photo_service "rent_service/internal/logic/services/interfaces/photo"
	"rent_service/internal/logic/services/types/token"
	"rent_service/server"
	"rent_service/server/authenticator"
	"rent_service/server/errstructs"
	getter_uuid "rent_service/server/getters/uuid"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type providers struct {
	photo photo.IProvider
}

type controller struct {
	providers     providers
	authenticator authenticator.IAuthenticator
}

func New(
	photo photo.IProvider,
	authenticator authenticator.IAuthenticator,
) server.IController {
	return &controller{providers{photo}, authenticator}
}

func CorsFiller(config *cors.Config) {
	config.AddAllowMethods("get", "post")
}

const PARAM_ID string = "id"

func (self *controller) Register(engine *gin.Engine) {
	engine.POST("/photos/temp", self.createTemp)
	engine.POST(fmt.Sprintf("/photos/temp/:%v", PARAM_ID), self.uploadTemp)
	engine.GET(fmt.Sprintf("/photos/:%v", PARAM_ID), self.get)
}

func (self *controller) createTemp(ctx *gin.Context) {
	var form photo_service.Description
	var token token.Token
	var id uuid.UUID
	err := ctx.ShouldBindJSON(&form)

	if nil == err {
		token, err = authenticator.ExchangeToken(ctx, self.authenticator)
	}

	if nil == err {
		service := self.providers.photo.GetPhotoService()
		id, err = service.CreateTempPhoto(token, form)
	}

	if nil == err {
		ctx.JSON(http.StatusCreated, id)
	} else if cerr := (cmnerrors.ErrorAuthentication{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusUnauthorized)
	} else if cerr := (cmnerrors.ErrorAuthorization{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusForbidden)
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else if cerr := (photo_service.ErrorUnsupportedMime{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusUnsupportedMediaType)
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

var ErrorMimeDoesntMatch = errors.New(
	"Provided mime doesn't match with created instance",
)

func (self *controller) uploadTemp(ctx *gin.Context) {
	var photo []byte
	var token token.Token
	var mime string
	service := self.providers.photo.GetPhotoService()
	id, err := getter_uuid.Parse(ctx.Params.ByName(PARAM_ID))

	if nil == err {
		mime = ctx.ContentType()

		if "" == mime {
			err = errors.New("No mime type provided")
		}
	}

	if nil == err {
		token, err = authenticator.ExchangeToken(ctx, self.authenticator)
	}

	if nil == err {
		data, cerr := service.GetTempPhoto(token, id)

		if nil != cerr {
			err = cerr
		} else if data.Mime != mime {
			err = ErrorMimeDoesntMatch
		}
	}

	if nil == err {
		photo, err = io.ReadAll(ctx.Request.Body)

		if nil != err {
			err = cmnerrors.Internal(err)
		}
	}

	if nil == err && 0 == len(photo) {
		err = errors.New("No photo was supplied")
	}

	if nil == err {
		err = service.UploadTempPhoto(token, id, photo)
	}

	if nil == err {
		ctx.Status(http.StatusCreated)
	} else if cerr := (cmnerrors.ErrorAuthentication{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusUnauthorized)
	} else if cerr := (cmnerrors.ErrorAuthorization{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusForbidden)
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusNotFound, errstructs.NewNotFound(cerr))
	} else if errors.Is(err, ErrorMimeDoesntMatch) {
		ctx.Status(http.StatusUnsupportedMediaType)
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func (self *controller) get(ctx *gin.Context) {
	var photo photo_service.Photo
	id, err := getter_uuid.Parse(ctx.Params.ByName(PARAM_ID))

	if nil == err {
		service := self.providers.photo.GetPhotoService()
		photo, err = service.GetPhoto(id)
	}

	if nil == err {
		ctx.JSON(http.StatusOK, photo)
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusNotFound, errstructs.NewNotFound(cerr))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

