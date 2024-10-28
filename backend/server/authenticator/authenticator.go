package authenticator

import (
	"errors"
	"github.com/gin-gonic/gin"
	"rent_service/internal/logic/services/types/token"
	"rent_service/server/headers"
)

type ApiToken struct {
	Access string `json:"token"`
	Renew  string `json:"renew_token"`
}

// Methods must return cmnerrors ErrorAuthentication, ErrorInternal or
// errors from this file
type IAuthenticator interface {
	Login(email string, password string) (ApiToken, error)
	GetToken(access string) (token.Token, error)
	RenewKey(token ApiToken) (ApiToken, error)
	Logout(access string) error
	Clear()
}

func ExchangeToken(
	ctx *gin.Context,
	authenticator IAuthenticator,
) (token.Token, error) {
	access := ctx.Request.Header.Get(headers.API_KEY)

	return authenticator.GetToken(access)
}

var ErrorNoApiKeyHeader = errors.New("No api key was provided")
var ErrorNoRenewHeader = errors.New("No renew key was provided")

