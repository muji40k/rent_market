package rent

import "rent_service/internal/repository/interfaces/rent"

type IProvider interface {
	GetRentRepository() rent.IRepository
}

type IRequestProvider interface {
	GetRentRequestRepository() rent.IRequestRepository
}

type IReturnProvider interface {
	GetRentReturnRepository() rent.IReturnRepository
}

