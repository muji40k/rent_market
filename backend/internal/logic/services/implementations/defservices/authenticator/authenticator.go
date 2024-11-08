package authenticator

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/types/token"
)

type IAuthenticator interface {
	LoginWithToken(token token.Token) (models.User, error)
}

