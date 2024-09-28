package rent

import "rent_service/internal/logic/services/interfaces/rent"

type IProvider interface {
	GetRentService() rent.IService
}

type IRequestProvider interface {
	GetRentRequestService() rent.IRequestService
}

type IReturnProvider interface {
	GetRentReturnService() rent.IReturnService
}

