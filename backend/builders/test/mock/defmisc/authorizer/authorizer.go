package authorizer

import (
	"rent_service/internal/logic/services/implementations/defservices/authorizer"
	"rent_service/internal/logic/services/implementations/defservices/authorizer/implementations/defauthorizer"
	"rent_service/internal/repository/implementation/mock/role"

	role_pmock "rent_service/internal/repository/context/mock/role"

	"go.uber.org/mock/gomock"
)

type AuthorizerBuilder struct {
	administrator *mock_role.MockIAdministratorRepository
	renter        *mock_role.MockIRenterRepository
	storekeeper   *mock_role.MockIStorekeeperRepository
}

func New(ctrl *gomock.Controller) *AuthorizerBuilder {
	return &AuthorizerBuilder{
		mock_role.NewMockIAdministratorRepository(ctrl),
		mock_role.NewMockIRenterRepository(ctrl),
		mock_role.NewMockIStorekeeperRepository(ctrl),
	}
}

func (self *AuthorizerBuilder) WithAdministratorRepository(f func(repo *mock_role.MockIAdministratorRepository)) *AuthorizerBuilder {
	f(self.administrator)
	return self
}

func (self *AuthorizerBuilder) WithRenterRepository(f func(repo *mock_role.MockIRenterRepository)) *AuthorizerBuilder {
	f(self.renter)
	return self
}

func (self *AuthorizerBuilder) WithStorekeeperRepository(f func(repo *mock_role.MockIStorekeeperRepository)) *AuthorizerBuilder {
	f(self.storekeeper)
	return self
}

func (self *AuthorizerBuilder) Build() authorizer.IAuthorizer {
	return defauthorizer.New(
		role_pmock.NewAdministrator(self.administrator),
		role_pmock.NewRenter(self.renter),
		role_pmock.NewStorekeeper(self.storekeeper),
	)
}

