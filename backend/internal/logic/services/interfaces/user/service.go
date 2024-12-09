package user

import (
	"fmt"
	"rent_service/internal/logic/services/types/token"

	"github.com/google/uuid"
)

type IService interface {
	GetSelfUserInfo(token token.Token) (Info, error)
	UpdateSelfUserInfo(token token.Token, form UpdateForm) error
	// Deprecated: moved to 2FA. Use IPasswordUpdateService instead.
	UpdateSelfUserPassword(
		token token.Token,
		old_password string,
		new_password string,
	) error
}

type IPasswordUpdateService interface {
	RequestPasswordUpdate(
		token token.Token,
		old_password string,
		new_password string,
	) (PasswordUpdateRequest, error)
	AuthenticatePasswordUpdateRequest(
		token token.Token,
		request uuid.UUID,
		code string,
	) error
}

type IProfileService interface {
	GetUserProfile(userId uuid.UUID) (UserProfile, error)
	GetSelfUserProfile(token token.Token) (UserProfile, error)
	UpdateSelfUserProfile(token token.Token, data UserProfile) error
}

type IFavoriteService interface {
	GetUserFavorite(userId uuid.UUID) (UserFavoritePickUpPoint, error)
	GetSelfUserFavorite(token token.Token) (UserFavoritePickUpPoint, error)
	UpdateSelfUserFavorite(
		token token.Token,
		data UserFavoritePickUpPoint,
	) error
}

type IRoleService interface {
	IsRenter(token token.Token) error
	RegisterAsRenter(token token.Token) error
	IsAdministrator(token token.Token) error
	IsStoreKeeper(token token.Token) (StoreKeeper, error)
}

type ErrorAlreadyRenter struct{ email string }

func AlreadyRenter(email string) ErrorAlreadyRenter {
	return ErrorAlreadyRenter{email}
}

func (e ErrorAlreadyRenter) Error() string {
	return fmt.Sprintf("User with email '%v' is already renter", e.email)
}

