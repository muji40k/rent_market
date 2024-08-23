
package models

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    Id uuid.UUID
    Name string
    Email string
    Password string
}

type UserProfile struct {
    Id uuid.UUID
    UserId uuid.UUID
    Name *string
    Surname *string
    Patronymic *string
    BirthDate *time.Time
    PhotoId *uuid.UUID
}

type UserFavoritePickUpPoint struct {
    Id uuid.UUID
    UserId uuid.UUID
    PickUpPointId *uuid.UUID
}

type UserPayMethods struct {
    UserId uuid.UUID
    Map map[uuid.UUID]struct {
        Method PayMethod
        Priority uint
    }
}

func NewUserPayMethods() UserPayMethods {
    out := UserPayMethods{}
    out.Map = make(map[uuid.UUID]struct {
        Method PayMethod
        Priority uint
    })

    return out
}

