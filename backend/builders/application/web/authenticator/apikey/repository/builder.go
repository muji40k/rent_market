package repository

import "rent_service/server/authenticator/implementations/apikey"

type IBuilder interface {
	Build() (apikey.ITokenRepository, error)
}

