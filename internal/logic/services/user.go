package services

import (
	"fmt"
	"rent_service/internal/domain/models"

	"github.com/google/uuid"
)

type UserInfo struct {
	Id    uuid.UUID
	Name  string
	Email string
}

type IUserService interface {
	GetSelfUserInfo(token Token) (UserInfo, error)
	UpdateUserInfo(token Token, info UserInfo) error
	UpdateUserPassword(
		token Token,
		old_password string,
		new_password string,
	) error
}

type IUserProfileService interface {
	GetUserProfile(userId uuid.UUID) (models.UserProfile, error)
	GetSelfUserProfile(token Token) (models.UserProfile, error)
	UpdateUserProfile(token Token, data models.UserProfile) error
}

type IUserFavoriteService interface {
	GetUserFavorite(userId uuid.UUID) (models.UserFavoritePickUpPoint, error)
	GetSelfUserFavorite(token Token) (models.UserFavoritePickUpPoint, error)
	UpdateUserFavorite(token Token, data models.UserFavoritePickUpPoint) error
}

type UserStoreKeeper struct {
	PickUpPointId uuid.UUID
}

type IRoleService interface {
	IsRenter(token Token) error
	RegisterAsRenter(token Token) error
	IsAdministrator(token Token) error
	IsStoreKeeper(token Token) (UserStoreKeeper, error)
}

type ErrorAlreadyRenter struct{ email string }

func (e ErrorAlreadyRenter) Error() string {
	return fmt.Sprintf("User with email '%v' is already renter", e.email)
}

