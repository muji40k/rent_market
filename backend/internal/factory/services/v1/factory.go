package v1

import "rent_service/internal/logic/context/v1"

type IFactory interface {
	ToFactories() v1.Factories
	Clear()
}

