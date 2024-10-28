package user

import "rent_service/internal/repository/interfaces/user"

type IFactory interface {
	CreateUserRepository() user.IRepository
}

type IProfileFactory interface {
	CreateUserProfileRepository() user.IProfileRepository
}

type IFavoriteFactory interface {
	CreateUserFavoriteRepository() user.IFavoriteRepository
}

type IPayMethodsFactory interface {
	CreateUserPayMethodsRepository() user.IPayMethodsRepository
}

