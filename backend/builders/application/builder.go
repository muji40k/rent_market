package application

import "rent_service/application"

type IBuilder interface {
	Build() (application.IApplication, error)
}

