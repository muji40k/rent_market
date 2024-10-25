package v1

import v1 "rent_service/internal/repository/context/v1"

type IBuilder interface {
	Build() (v1.Factories, error)
}

