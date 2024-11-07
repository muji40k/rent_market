package models

import (
	modelsb "rent_service/builders/domain/models"
	"rent_service/builders/misc/uuidgen"

	"github.com/google/uuid"
)

func AdministratorRandomId() *modelsb.AdministratorBuilder {
	return modelsb.NewAdministrator().
		WithId(uuidgen.Generate())
}

func AdministratorWithUserId(userId uuid.UUID) *modelsb.AdministratorBuilder {
	return AdministratorRandomId().
		WithUserId(userId)
}

