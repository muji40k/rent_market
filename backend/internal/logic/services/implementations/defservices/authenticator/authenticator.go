package authenticator

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/types/token"
)

//go:generate mockgen -source=authenticator.go -destination=implementations/mock/authenticator.go

type IAuthenticator interface {
	LoginWithToken(token token.Token) (models.User, error)
}

