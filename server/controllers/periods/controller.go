package periods

import (
	"errors"
	"net/http"
	"rent_service/internal/logic/context/providers/period"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/misc/types/collection"
	"rent_service/server"
	"rent_service/server/errstructs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type providers struct {
	period period.IProvider
}

type controller struct {
	providers providers
}

func New(period period.IProvider) server.IController {
	return &controller{providers{period}}
}

func CorsFiller(config *cors.Config) {
	config.AddAllowMethods("get")
}

func (self *controller) Register(engine *gin.Engine) {
	engine.GET("/periods", self.get)
}

func (self *controller) get(ctx *gin.Context) {
	service := self.providers.period.GetPeriodService()
	col, err := service.GetPeriods()

	if nil == err {
		ctx.JSON(http.StatusOK, collection.Marshaler(col.Iter()))
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

