package v1

import v1 "rent_service/internal/factory/services/v1"

type IBuilder interface {
	Build() (v1.IFactory, error)
}

