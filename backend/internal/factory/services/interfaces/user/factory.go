package user

import "rent_service/internal/logic/services/interfaces/user"

type IFactory interface {
	CreateUserService() user.IService
}

type IProfileFactory interface {
	CreateUserProfileService() user.IProfileService
}

type IFavoriteFactory interface {
	CreateUserFavoriteService() user.IFavoriteService
}

type IRoleFactory interface {
	CreateRoleService() user.IRoleService
}

