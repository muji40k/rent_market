package profiles

import (
	"errors"
	"fmt"
	"net/http"
	"rent_service/internal/logic/context/providers/payment"
	"rent_service/internal/logic/context/providers/user"
	"rent_service/internal/logic/services/errors/cmnerrors"
	payment_service "rent_service/internal/logic/services/interfaces/payment"
	user_service "rent_service/internal/logic/services/interfaces/user"
	"rent_service/internal/logic/services/types/token"
	"rent_service/internal/misc/types/collection"
	"rent_service/server"
	"rent_service/server/authenticator"
	"rent_service/server/errstructs"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	PARAM_ID       string = "id"
	PARAM_CATEGORY string = "category"
)

type providers struct {
	profile   user.IProfileProvider
	favorite  user.IFavoriteProvider
	payMethod payment.IUserPayMethodProvider
}

type controller struct {
	providers     providers
	authenticator authenticator.IAuthenticator
}

func New(
	authenticator authenticator.IAuthenticator,
	profile user.IProfileProvider,
	favorite user.IFavoriteProvider,
	payMethod payment.IUserPayMethodProvider,
) server.IController {
	return &controller{providers{profile, favorite, payMethod}, authenticator}
}

func (self *controller) Register(engine *gin.Engine) {
	engine.PATCH("/profiles/self", self.update)
	engine.GET(fmt.Sprintf("/profiles/:%v", PARAM_ID), self.get)
}

func updateGeneral(
	ctx *gin.Context,
	token token.Token,
	service user_service.IProfileService,
) error {
	var profile user_service.UserProfile
	err := ctx.ShouldBindJSON(&profile)

	if nil == err {
		err = service.UpdateSelfUserProfile(token, profile)
	}

	return err
}

func updateFavorite(
	ctx *gin.Context,
	token token.Token,
	service user_service.IFavoriteService,
) error {
	var favorite user_service.UserFavoritePickUpPoint
	err := ctx.ShouldBindJSON(&favorite)

	if nil == err {
		err = service.UpdateSelfUserFavorite(token, favorite)
	}

	return err
}

func updatePayMethods(
	ctx *gin.Context,
	token token.Token,
	service payment_service.IUserPayMethodService,
) error {
	var methods struct {
		Delete []uuid.UUID `json:"delete"`
		Keep   []struct {
			Move     *uuid.UUID                                 `json:"move"`
			Register *payment_service.PayMethodRegistrationForm `json:"register"`
		} `json:"items"`
	}

	err := ctx.ShouldBindJSON(&methods)

	// Check form correctness first
	for i := 0; nil == err && len(methods.Keep) > i; i++ {
		item := methods.Keep[i]

		if (nil == item.Move && nil == item.Register) ||
			(nil != item.Move && nil != item.Register) {
			err = errors.New(
				"Item must specify only one of methods: 'move' or 'register'",
			)
		} else if nil == item.Register &&
			slices.Contains(methods.Delete, *item.Move) {
			err = errors.New("Move of deleted method")
		}
	}

	// Delete
	for i := 0; nil == err && len(methods.Delete) > i; i++ {
		err = service.RemovePayMethod(token, methods.Delete[i])
	}

	// Compose new order of methods and register new ones
	var order []uuid.UUID

	if nil == err {
		order = make([]uuid.UUID, len(methods.Keep))
	}

	for i := 0; nil == err && len(methods.Keep) > i; i++ {
		item := methods.Keep[i]

		if nil != item.Move {
			order[i] = *item.Move
		} else {
			var id uuid.UUID
			id, err = service.RegisterPayMethod(token, *item.Register)

			if nil == err {
				order[i] = id
			}
		}
	}

	// Apply new order
	if nil == err {
		err = service.UpdatePayMethodsPriority(token, order)
	}

	return err
}

func (self *controller) update(ctx *gin.Context) {
	token, err := authenticator.ExchangeToken(ctx, self.authenticator)

	if nil == err {
		switch category := ctx.Request.URL.Query().Get(PARAM_CATEGORY); category {
		case "general":
			err = updateGeneral(
				ctx,
				token,
				self.providers.profile.GetUserProfileService(),
			)
		case "favorite":
			err = updateFavorite(
				ctx,
				token,
				self.providers.favorite.GetUserFavoriteService(),
			)
		case "pay-methods":
			err = updatePayMethods(
				ctx,
				token,
				self.providers.payMethod.GetUserPayMethodService(),
			)
		case "":
			err = errors.New("No category provided")
		default:
			err = fmt.Errorf("Unsupported category: '%v'", category)
		}
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

type profile struct {
	General    user_service.UserProfile             `json:"general"`
	Favorite   user_service.UserFavoritePickUpPoint `json:"favorite"`
	PayMethods *[]payment_service.UserPayMethod     `json:"pay_methods,omitempty"`
}

func (self *controller) getSelf(ctx *gin.Context) (profile, error) {
	var profile profile
	token, err := authenticator.ExchangeToken(ctx, self.authenticator)

	if nil == err {
		service := self.providers.profile.GetUserProfileService()
		profile.General, err = service.GetSelfUserProfile(token)

		if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	if nil == err {
		service := self.providers.favorite.GetUserFavoriteService()
		profile.Favorite, err = service.GetSelfUserFavorite(token)

		if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	if nil == err {
		service := self.providers.payMethod.GetUserPayMethodService()
		methods, cerr := service.GetPayMethods(token)

		if cnferr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cnferr) {
			err = nil
		} else if nil != cerr {
			err = cerr
		} else {
			profile.PayMethods = new([]payment_service.UserPayMethod)
			*profile.PayMethods = collection.Collect(methods.Iter())
		}
	}

	return profile, err
}

func (self *controller) getById(
	uuid uuid.UUID,
) (profile, error) {
	var profile profile
	var err error

	service := self.providers.profile.GetUserProfileService()
	profile.General, err = service.GetUserProfile(uuid)

	if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = nil
	}

	if nil == err {
		service := self.providers.favorite.GetUserFavoriteService()
		profile.Favorite, err = service.GetUserFavorite(uuid)

		if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	return profile, err
}

func (self *controller) get(ctx *gin.Context) {
	var err error
	var profile profile
	rawId := ctx.Params.ByName(PARAM_ID)

	if "" == rawId {
		err = errors.New("No id specified in path...")
	} else if "self" == rawId {
		profile, err = self.getSelf(ctx)
	} else if uuid, cerr := uuid.Parse(rawId); nil == cerr {
		profile, err = self.getById(uuid)
	} else {
		err = fmt.Errorf("Error parsing uuid: %w", cerr)
	}

	if nil == err {
		ctx.JSON(http.StatusOK, profile)
	} else if cerr := (cmnerrors.ErrorAuthentication{}); errors.As(err, &cerr) {
		ctx.Status(http.StatusUnauthorized)
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
		ctx.JSON(http.StatusInternalServerError, errstructs.NewInternalErr(err))
	} else {
		ctx.JSON(http.StatusBadRequest, errstructs.NewBadRequestErr(err))
	}
}

