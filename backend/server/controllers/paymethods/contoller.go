package paymethods

import (
	"errors"
	"net/http"
	"rent_service/internal/logic/context/providers/payment"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/misc/types/collection"
	"rent_service/server"
	"rent_service/server/errstructs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type providers struct {
	payment payment.IPayMethodProvider
}

type controller struct {
	providers providers
}

func New(payment payment.IPayMethodProvider) server.IController {
	return &controller{providers{payment}}
}

func CorsFiller(config *cors.Config) {
	config.AddAllowMethods("get")
}

func (self *controller) Register(engine *gin.Engine) {
	engine.GET("/pay-methods", self.get)
}

func (self *controller) get(ctx *gin.Context) {
	service := self.providers.payment.GetPayMethodService()
	col, err := service.GetPayMethods()

	if nil == err {
		ctx.JSON(http.StatusOK, collection.Marshaler(col.Iter()))
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

