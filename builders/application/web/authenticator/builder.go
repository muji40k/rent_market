package authenticator

import "rent_service/server/authenticator"

type IBuilder interface {
	Build() (authenticator.IAuthenticator, error)
}

