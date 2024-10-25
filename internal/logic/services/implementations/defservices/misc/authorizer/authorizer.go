package authorizer

import (
	"errors"
	"fmt"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/repository/context/providers/role"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"

	"github.com/google/uuid"
)

type repoproviders struct {
	administrator role.IAdministratorProvider
	renter        role.IRenterProvider
	storekeeper   role.IStorekeeperProvider
}

type Authorizer struct {
	repos repoproviders
}

func New(
	administrator role.IAdministratorProvider,
	renter role.IRenterProvider,
	storekeeper role.IStorekeeperProvider,
) *Authorizer {
	return &Authorizer{repoproviders{administrator, renter, storekeeper}}
}

func (self *Authorizer) IsAdministrator(userId uuid.UUID) (models.Administrator, error) {
	repo := self.repos.administrator.GetAdministratorRepository()
	admin, err := repo.GetByUserId(userId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.Authorization(ErrorUnauthorized{userId, "administrator"})
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return admin, err
}

func (self *Authorizer) IsRenter(userId uuid.UUID) (models.Renter, error) {
	repo := self.repos.renter.GetRenterRepository()
	renter, err := repo.GetByUserId(userId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.Authorization(ErrorUnauthorized{userId, "renter"})
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return renter, err
}

func (self *Authorizer) IsStorekeeper(userId uuid.UUID) (models.Storekeeper, error) {
	repo := self.repos.storekeeper.GetStorekeeperRepository()
	storekeeper, err := repo.GetByUserId(userId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.Authorization(ErrorUnauthorized{userId, "storekeeper"})
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return storekeeper, err
}

type ErrorUnauthorized struct {
	id   uuid.UUID
	Role string
}

func (e ErrorUnauthorized) Error() string {
	return fmt.Sprintf(
		"User '%v' attempt to authorize to insufficient role '%v'",
		e.id, e.Role,
	)
}

