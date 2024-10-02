package login

import "rent_service/internal/logic/services/types/token"

type IService interface {
	Register(email string, password string, name string) error
	Login(email string, password string) (token.Token, error)
}

