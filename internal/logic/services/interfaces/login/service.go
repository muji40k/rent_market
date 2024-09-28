package login

import "rent_service/internal/domain/models"

type IService interface {
	Register(email string, password string, name string) error
	Login(email string, password string) (models.Token, error)
}

