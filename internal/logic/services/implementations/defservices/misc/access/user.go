package access

import (
	"errors"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/misc/authorizer"
	"rent_service/internal/repository/context/providers/user"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"

	"github.com/google/uuid"
)

type userProviders struct {
	user user.IProvider
}

type User struct {
	repos      userProviders
	authorizer *authorizer.Authorizer
}

func (self *User) Access(rqUserId uuid.UUID, userId uuid.UUID) error {
	repo := self.repos.user.GetUserRepository()
	_, err := repo.GetById(userId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		return cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		return cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	if _, aerr := self.authorizer.IsAdministrator(userId); nil == aerr {
		return nil
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(aerr, &cerr) {
		return aerr
	}

	if rqUserId == userId {
		return nil
	}

	return cmnerrors.Authorization(cmnerrors.NoAccess("user"))
}
