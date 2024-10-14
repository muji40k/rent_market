package user

import (
	"rent_service/internal/domain/models"
	// "rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

type IRepository interface {
	Create(user models.User) (models.User, error)

	Update(user models.User) error

	GetById(userId uuid.UUID) (models.User, error)
	GetByEmail(email string) (models.User, error)
	GetByToken(token models.Token) (models.User, error)
}

type IProfileRepository interface {
	Create(profile models.UserProfile) (models.UserProfile, error)

	Update(profile models.UserProfile) error

	GetByUserId(userId uuid.UUID) (models.UserProfile, error)
}

type IFavoriteRepository interface {
	Create(
		profile models.UserFavoritePickUpPoint,
	) (models.UserFavoritePickUpPoint, error)

	Update(profile models.UserFavoritePickUpPoint) error

	GetByUserId(userId uuid.UUID) (models.UserFavoritePickUpPoint, error)
}

type IPayMethodsRepository interface {
	CreatePayMethod(
		userId uuid.UUID,
		payMethod models.UserPayMethod,
	) (models.UserPayMethods, error)

	Update(payMethods models.UserPayMethods) error

	GetByUserId(userId uuid.UUID) (models.UserPayMethods, error)
}

