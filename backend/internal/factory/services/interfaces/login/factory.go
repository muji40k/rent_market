package login

import "rent_service/internal/logic/services/interfaces/login"

type IFactory interface {
	CreateLoginService() login.IService
}

