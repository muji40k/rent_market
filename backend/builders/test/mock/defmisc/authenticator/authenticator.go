package authenticator

import (
	"rent_service/internal/logic/services/implementations/defservices/authenticator"
	"rent_service/internal/logic/services/implementations/defservices/authenticator/implementations/defauth"
	"rent_service/internal/repository/context/mock/user"
	"rent_service/internal/repository/implementation/mock/user"

	"go.uber.org/mock/gomock"
)

type AuthenticatorBuilder struct {
	user *mock_user.MockIRepository
}

func New(ctrl *gomock.Controller) *AuthenticatorBuilder {
	return &AuthenticatorBuilder{
		mock_user.NewMockIRepository(ctrl),
	}
}

func (self *AuthenticatorBuilder) WithUserRepository(f func(repo *mock_user.MockIRepository)) *AuthenticatorBuilder {
	f(self.user)
	return self
}

func (self *AuthenticatorBuilder) Build() authenticator.IAuthenticator {
	return defauth.New(user.New(self.user))
}

