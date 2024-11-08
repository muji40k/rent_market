package defauthorizer

import (
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/authorizer"
	"rent_service/internal/repository/context/providers/role"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"

	"github.com/google/uuid"
)

type repoproviders struct {
	administrator role.IAdministratorProvider
	renter        role.IRenterProvider
	storekeeper   role.IStorekeeperProvider
}

type auth struct {
	repos repoproviders
}

func New(
	administrator role.IAdministratorProvider,
	renter role.IRenterProvider,
	storekeeper role.IStorekeeperProvider,
) authorizer.IAuthorizer {
	return &auth{repoproviders{administrator, renter, storekeeper}}
}

func (self *auth) IsAdministrator(userId uuid.UUID) (models.Administrator, error) {
	repo := self.repos.administrator.GetAdministratorRepository()
	admin, err := repo.GetByUserId(userId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.Authorization(authorizer.Unauthorized(userId, "administrator"))
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return admin, err
}

func (self *auth) IsRenter(userId uuid.UUID) (models.Renter, error) {
	repo := self.repos.renter.GetRenterRepository()
	renter, err := repo.GetByUserId(userId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.Authorization(authorizer.Unauthorized(userId, "renter"))
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return renter, err
}

func (self *auth) IsStorekeeper(userId uuid.UUID) (models.Storekeeper, error) {
	repo := self.repos.storekeeper.GetStorekeeperRepository()
	storekeeper, err := repo.GetByUserId(userId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.Authorization(authorizer.Unauthorized(userId, "storekeeper"))
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return storekeeper, err
}

