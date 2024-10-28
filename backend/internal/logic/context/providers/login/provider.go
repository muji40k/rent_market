package login

import "rent_service/internal/logic/services/interfaces/login"

type IProvider interface {
	GetLoginService() login.IService
}

