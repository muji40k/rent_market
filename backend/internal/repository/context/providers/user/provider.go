package user

import "rent_service/internal/repository/interfaces/user"

type IProvider interface {
	GetUserRepository() user.IRepository
}

type IPasswordUpdateProvider interface {
	GetUserPasswordUpdateRepository() user.IPasswordUpdateRepository
}

type IProfileProvider interface {
	GetUserProfileRepository() user.IProfileRepository
}

type IFavoriteProvider interface {
	GetUserFavoriteRepository() user.IFavoriteRepository
}

type IPayMethodsProvider interface {
	GetUserPayMethodsRepository() user.IPayMethodsRepository
}

