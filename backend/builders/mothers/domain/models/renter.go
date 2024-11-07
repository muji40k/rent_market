package models

import (
	modelsb "rent_service/builders/domain/models"
	"rent_service/builders/misc/uuidgen"

	"github.com/google/uuid"
)

func RenterRandomId() *modelsb.RenterBuilder {
	return modelsb.NewRenter().
		WithId(uuidgen.Generate())
}

func RenterWithUserId(userId uuid.UUID) *modelsb.RenterBuilder {
	return RenterRandomId().
		WithUserId(userId)
}

