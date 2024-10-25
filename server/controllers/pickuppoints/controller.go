package pickuppoints

import (
	"errors"
	"fmt"
	"net/http"
	"rent_service/internal/logic/context/providers/pickuppoint"
	"rent_service/internal/logic/services/errors/cmnerrors"
	pickuppoint_service "rent_service/internal/logic/services/interfaces/pickuppoint"
	"rent_service/internal/misc/types/collection"
	"rent_service/server"
	"rent_service/server/errstructs"
	getter_uuid "rent_service/server/getters/uuid"
	"rent_service/server/pagination"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type providers struct {
	pickUpPoint pickuppoint.IProvider
	photo       pickuppoint.IPhotoProvider
	wh          pickuppoint.IWorkingHoursProvider
}

type controller struct {
	providers providers
}

func New(
	pickUpPoint pickuppoint.IProvider,
	photo pickuppoint.IPhotoProvider,
	wh pickuppoint.IWorkingHoursProvider,
) server.IController {
	return &controller{providers{pickUpPoint, photo, wh}}
}

func CorsFiller(config *cors.Config) {
	config.AddAllowMethods("get")
}

const PARAM_ID string = "id"

func (self *controller) Register(engine *gin.Engine) {
	engine.GET("/pick-up-points", self.get)
	engine.GET(fmt.Sprintf("/pick-up-points/:%v", PARAM_ID), self.getById)
	engine.GET(fmt.Sprintf("/pick-up-points/:%v/photos", PARAM_ID), self.getPhotos)
	engine.GET(fmt.Sprintf("/pick-up-points/:%v/working-hours", PARAM_ID), self.getWH)
}

func (self *controller) get(ctx *gin.Context) {
	var iter collection.Iterator[pickuppoint_service.PickUpPoint]
	service := self.providers.pickUpPoint.GetPickUpPointService()
	col, err := service.ListPickUpPoints()

	if nil == err {
		iter, err = pagination.Apply(ctx, col.Iter())
	}

	if nil == err {
		ctx.JSON(http.StatusOK, collection.Marshaler(iter))
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func (self *controller) getById(ctx *gin.Context) {
	var pickUpPoint pickuppoint_service.PickUpPoint
	id, err := getter_uuid.Parse(ctx.Params.ByName(PARAM_ID))

	if nil == err {
		service := self.providers.pickUpPoint.GetPickUpPointService()
		pickUpPoint, err = service.GetPickUpPointById(id)
	}

	if nil == err {
		ctx.JSON(http.StatusOK, pickUpPoint)
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusNotFound, errstructs.NewNotFound(cerr))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func (self *controller) getPhotos(ctx *gin.Context) {
	var photos collection.Collection[uuid.UUID]
	id, err := getter_uuid.Parse(ctx.Params.ByName(PARAM_ID))

	if nil == err {
		service := self.providers.photo.GetPickUpPointPhotoService()
		photos, err = service.ListPickUpPointPhotos(id)
	}

	if nil == err {
		ctx.JSON(http.StatusOK, collection.Marshaler(photos.Iter()))
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusNotFound, errstructs.NewNotFound(cerr))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func (self *controller) getWH(ctx *gin.Context) {
	var whs collection.Collection[pickuppoint_service.WorkingHours]
	id, err := getter_uuid.Parse(ctx.Params.ByName(PARAM_ID))

	if nil == err {
		service := self.providers.wh.GetPickUpPointWorkingHoursService()
		whs, err = service.ListPickUpPointWorkingHours(id)
	}

	if nil == err {
		ctx.JSON(http.StatusOK, collection.Marshaler(whs.Iter()))
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusNotFound, errstructs.NewNotFound(cerr))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

