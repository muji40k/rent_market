package user

import (
	"fmt"
	"rent_service/internal/domain/models"

	"github.com/google/uuid"
)

type Info struct {
	Id    uuid.UUID
	Name  string
	Email string
}

type IService interface {
	GetSelfUserInfo(token models.Token) (Info, error)
	UpdateUserInfo(token models.Token, info Info) error
	UpdateUserPassword(
		token models.Token,
		old_password string,
		new_password string,
	) error
}

type IProfileService interface {
	GetUserProfile(userId uuid.UUID) (models.UserProfile, error)
	GetSelfUserProfile(token models.Token) (models.UserProfile, error)
	UpdateUserProfile(token models.Token, data models.UserProfile) error
}

type IFavoriteService interface {
	GetUserFavorite(
		userId uuid.UUID,
	) (models.UserFavoritePickUpPoint, error)
	GetSelfUserFavorite(
		token models.Token,
	) (models.UserFavoritePickUpPoint, error)
	UpdateUserFavorite(
		token models.Token,
		data models.UserFavoritePickUpPoint,
	) error
}

type StoreKeeper struct {
	PickUpPointId uuid.UUID
}

type IRoleService interface {
	IsRenter(token models.Token) error
	RegisterAsRenter(token models.Token) error
	IsAdministrator(token models.Token) error
	IsStoreKeeper(token models.Token) (StoreKeeper, error)
}

type ErrorAlreadyRenter struct{ email string }

func (e ErrorAlreadyRenter) Error() string {
	return fmt.Sprintf("User with email '%v' is already renter", e.email)
}

