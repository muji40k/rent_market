package services

type Token string

type ILoginService interface {
	Register(email string, password string, name string) error
	Login(email string, password string) (Token, error)
}

