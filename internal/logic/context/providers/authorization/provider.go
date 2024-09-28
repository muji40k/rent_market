package authorization

import "rent_service/internal/logic/services/interfaces/authorization"

type IProvider interface {
	GetAuthorizationService() authorization.IService
}

