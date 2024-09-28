package authorization

import "rent_service/internal/logic/services/interfaces/authorization"

type IFactory interface {
	CreateAuthorizationService() authorization.IService
}

