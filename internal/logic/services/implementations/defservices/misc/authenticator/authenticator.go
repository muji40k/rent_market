package authenticator

import (
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/types/token"
	"rent_service/internal/repository/context/providers/user"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
)

type repoproviders struct {
	user user.IProvider
}

type Authenticator struct {
	repos repoproviders
}

func New(user user.IProvider) *Authenticator {
	return &Authenticator{repoproviders{user}}
}

func (self *Authenticator) LoginWithToken(token token.Token) (models.User, error) {
	var user models.User
	var err error

	if "" == token {
		err = cmnerrors.Empty("token")
	}

	if nil == err {
		repo := self.repos.user.GetUserRepository()
		user, err = repo.GetByToken(models.Token(token))

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.Authentication(cmnerrors.NotFound(cerr.What...))
		} else if err != nil {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return user, err
}

