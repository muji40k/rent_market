package user

import "rent_service/internal/logic/services/interfaces/user"

type IProvider interface {
	GetUserService() user.IService
}

type IPasswordUpdateProvider interface {
	GetUserPasswordUpdateService() user.IPasswordUpdateService
}

type IProfileProvider interface {
	GetUserProfileService() user.IProfileService
}

type IFavoriteProvider interface {
	GetUserFavoriteService() user.IFavoriteService
}

type IRoleProvider interface {
	GetRoleService() user.IRoleService
}

