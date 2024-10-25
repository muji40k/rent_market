package categories

import (
	"errors"
	"fmt"
	"net/http"
	"rent_service/internal/logic/context/providers/category"
	"rent_service/internal/logic/services/errors/cmnerrors"
	category_service "rent_service/internal/logic/services/interfaces/category"
	"rent_service/internal/misc/types/collection"
	"rent_service/server"
	"rent_service/server/errstructs"
	getter_uuid "rent_service/server/getters/uuid"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type providers struct {
	category category.IProvider
}

type controller struct {
	providers providers
}

func New(category category.IProvider) server.IController {
	return &controller{providers{category}}
}

func CorsFiller(config *cors.Config) {
	config.AddAllowMethods("get")
}

const PARAM_ID = "id"

func (self *controller) Register(engine *gin.Engine) {
	engine.GET("/categories", self.get)
	engine.GET(fmt.Sprintf("/categories/:%v", PARAM_ID), self.getById)
}

func (self *controller) get(ctx *gin.Context) {
	service := self.providers.category.GetCategoryService()
	col, err := service.ListCategories()

	if nil == err {
		ctx.JSON(http.StatusOK, collection.Marshaler(col.Iter()))
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

func (self *controller) getById(ctx *gin.Context) {
	var col collection.Collection[category_service.Category]
	id, err := getter_uuid.Parse(ctx.Params.ByName(PARAM_ID))

	if nil == err {
		service := self.providers.category.GetCategoryService()
		col, err = service.GetFullCategory(id)
	}

	if nil == err {
		ctx.JSON(http.StatusOK, collection.Marshaler(col.Iter()))
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

