package rent

import "rent_service/internal/logic/services/interfaces/rent"

type IFactory interface {
	CreateRentService() rent.IService
}

type IRequestFactory interface {
	CreateRentRequestService() rent.IRequestService
}

type IReturnFactory interface {
	CreateRentReturnService() rent.IReturnService
}

