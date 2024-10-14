package login

import (
	"errors"
	"fmt"
	"regexp"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/misc/emptymathcer"
	"rent_service/internal/logic/services/interfaces/login"
	"rent_service/internal/logic/services/types/token"
	"rent_service/internal/repository/context/providers/user"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
)

type repoproviders struct {
	user user.IProvider
}

type service struct {
	repos repoproviders
}

func New(user user.IProvider) login.IService {
	return &service{repoproviders{user}}
}

const (
	ALLOWED_LOCAL  string = "[\\w!#$%&'*+\\-\\/=?^_`{|}~]"
	ALLOWED_DOMAIN string = "\\w"
)

var emailRegex = regexp.MustCompile(fmt.Sprintf(
	"^(%v\\.{0,1})*%v@%v\\.{0,1}([%v\\-]\\.{0,1})*%v$",
	ALLOWED_LOCAL, ALLOWED_LOCAL, ALLOWED_DOMAIN, ALLOWED_DOMAIN,
	ALLOWED_DOMAIN,
))

func (self *service) Register(email string, password string, name string) error {
	user := models.User{
		Email:    email,
		Name:     name,
		Password: password,
	}

	err := emptymathcer.Match(
		emptymathcer.Item("email", email),
		emptymathcer.Item("password", password),
		emptymathcer.Item("name", name),
	)

	if nil == err && !emailRegex.MatchString(email) {
		err = cmnerrors.Incorrect("email")
	}

	if nil == err {
		repo := self.repos.user.GetUserRepository()
		_, err = repo.Create(user)

		if cerr := (repo_errors.ErrorDuplicate{}); errors.As(err, &cerr) {
			err = cmnerrors.AlreadyExists(cerr.What...)
		} else if err != nil {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return nil
}

func (self *service) Login(email string, password string) (token.Token, error) {
	var user models.User
	var err error

	if "" == email {
		err = cmnerrors.Empty("email")
	}

	if nil == err {
		repo := self.repos.user.GetUserRepository()
		user, err = repo.GetByEmail(email)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.Authentication(cmnerrors.NotFound(cerr.What...))
		} else if err != nil {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	if nil == err && user.Password != password {
		err = cmnerrors.Authentication(errors.New("Passwords don't match"))
	}

	return token.Token(user.Token), err
}

