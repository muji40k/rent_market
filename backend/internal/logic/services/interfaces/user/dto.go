package user

import (
	"rent_service/internal/logic/services/types/date"

	"github.com/google/uuid"
)

type Info struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type UpdateForm struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type PasswordUpdateRequest struct {
	Required bool      `json:"required"`
	Id       uuid.UUID `json:"id"`
	ValidTo  date.Date `json:"valid_to"`
}

type UserProfile struct {
	Name       *string    `json:"name"`
	Surname    *string    `json:"surname"`
	Patronymic *string    `json:"patronymic"`
	BirthDate  *date.Date `json:"birth_date"`
	PhotoId    *uuid.UUID `json:"photo"`
}

type UserFavoritePickUpPoint struct {
	PickUpPointId *uuid.UUID `json:"pick_up_point"`
}

type StoreKeeper struct {
	PickUpPointId uuid.UUID `json:"pick_up_point"`
}

