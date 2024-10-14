package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id       uuid.UUID
	Name     string
	Email    string
	Password string
	Token    Token
}

type UserProfile struct {
	Id         uuid.UUID
	UserId     uuid.UUID
	Name       *string
	Surname    *string
	Patronymic *string
	BirthDate  *time.Time
	PhotoId    *uuid.UUID
}

type UserFavoritePickUpPoint struct {
	Id            uuid.UUID
	UserId        uuid.UUID
	PickUpPointId *uuid.UUID
}

type UserPayMethod struct {
	Name     string
	MethodId uuid.UUID
	PayerId  string
	Priority uint
}

type UserPayMethods struct {
	UserId uuid.UUID
	Map    map[uuid.UUID]UserPayMethod
}

// func NewUserPayMethods() UserPayMethods {
//     out := UserPayMethods{}
//     out.Map = make(map[uuid.UUID]struct {
//         Name     string
//         Method   PayMethod
//         PayerId  string
//         Priority uint
//     })
//
//     return out
// }

