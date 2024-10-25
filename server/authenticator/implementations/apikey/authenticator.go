package apikey

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"rent_service/internal/logic/context/providers/login"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/types/token"
	"rent_service/server/authenticator"
	"time"
)

type TokenHandle struct {
	Value   string
	ValidTo time.Time
}

type ITokenRepository interface {
	WriteToken(
		token token.Token,
		access TokenHandle,
		renew TokenHandle,
	) error
	GetToken(access string) (token.Token, error)
	DeleteToken(access string) error
	RenewToken(token authenticator.ApiToken) (token.Token, error)
}

var ErrorNotFound = errors.New("Can't validate token")
var ErrorDataAccess = errors.New("Error during data access")

type auth struct {
	login       login.IProvider
	repo        ITokenRepository
	validAccess time.Duration
	validRenew  time.Duration
}

func New(
	login login.IProvider,
	repo ITokenRepository,
) authenticator.IAuthenticator {
	return &auth{login, repo, 24 * time.Hour, 7 * 24 * time.Hour}
}

func (self *auth) GetToken(access string) (token.Token, error) {
	var token token.Token
	var err error

	if "" == access {
		err = authenticator.ErrorNoApiKeyHeader
	}

	if nil == err {
		token, err = self.repo.GetToken(access)

		if errors.Is(err, ErrorNotFound) {
			err = cmnerrors.Authentication(ErrorNotFound)
		} else if nil != err {
			err = cmnerrors.Internal(ErrorDataAccess)
		}
	}

	return token, err
}

func (self *auth) getToken(token token.Token) (TokenHandle, TokenHandle) {
	now := time.Now()
	accessValidTo := now.Add(self.validAccess)
	renewValidTo := now.Add(self.validRenew)

	access := getAccess(token, accessValidTo)
	renew := getRenew(token, access, renewValidTo)

	return TokenHandle{access, accessValidTo}, TokenHandle{renew, renewValidTo}
}

func getAccess(token token.Token, validTo time.Time) string {
	return toToken(fmt.Sprintf("a%vb%vc", token, validTo))
}

func getRenew(token token.Token, access string, validTo time.Time) string {
	return toToken(fmt.Sprintf("a%vb%vc%vd", token, access, validTo))
}

func toToken(value string) string {
	res := md5.Sum([]byte(value))
	return hex.EncodeToString(res[:])
}

func (self *auth) Login(email string, password string) (authenticator.ApiToken, error) {
	var apiToken authenticator.ApiToken
	service := self.login.GetLoginService()
	token, err := service.Login(email, password)

	if nil == err {
		handleAccess, handleRenew := self.getToken(token)
		err = self.repo.WriteToken(token, handleAccess, handleRenew)

		if nil == err {
			apiToken.Access = handleAccess.Value
			apiToken.Renew = handleRenew.Value
		} else if errors.Is(err, ErrorNotFound) {
			err = cmnerrors.Authentication(ErrorNotFound)
		} else {
			err = cmnerrors.Internal(ErrorDataAccess)
		}
	}

	return apiToken, err
}

func (self *auth) Logout(access string) error {
	var err error

	if "" == access {
		err = authenticator.ErrorNoApiKeyHeader
	}

	if nil == err {
		err = self.repo.DeleteToken(access)

		if errors.Is(err, ErrorNotFound) {
			err = cmnerrors.Authentication(ErrorNotFound)
		} else if nil != err {
			err = cmnerrors.Internal(ErrorDataAccess)
		}
	}

	return err
}

func (self *auth) RenewKey(apiToken authenticator.ApiToken) (authenticator.ApiToken, error) {
	var err error
	var token token.Token

	if "" == apiToken.Access {
		err = authenticator.ErrorNoApiKeyHeader
	}

	if nil == err && "" == apiToken.Renew {
		err = authenticator.ErrorNoApiKeyHeader
	}

	if nil == err {
		token, err = self.repo.RenewToken(apiToken)
	}

	if nil == err {
		handleAccess, handleRenew := self.getToken(token)
		err = self.repo.WriteToken(token, handleAccess, handleRenew)

		if nil == err {
			apiToken.Access = handleAccess.Value
			apiToken.Renew = handleRenew.Value
		} else if errors.Is(err, ErrorNotFound) {
			err = cmnerrors.Authentication(ErrorNotFound)
		} else {
			err = cmnerrors.Internal(ErrorDataAccess)
		}
	}

	return apiToken, err
}

