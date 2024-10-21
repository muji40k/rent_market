package rent

import "rent_service/internal/repository/interfaces/rent"

type IFactory interface {
	CreateRentRepository() rent.IRepository
}

type IRequestFactory interface {
	CreateRentRequestRepository() rent.IRequestRepository
}

type IReturnFactory interface {
	CreateRentReturnRepository() rent.IReturnRepository
}

