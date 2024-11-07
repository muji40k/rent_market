package models

import (
	modelsb "rent_service/builders/domain/models"
	"rent_service/builders/misc/uuidgen"

	"github.com/google/uuid"
)

func StorekeeperRandomId() *modelsb.StorekeeperBuilder {
	return modelsb.NewStorekeeper().
		WithId(uuidgen.Generate())
}

func StorekeeperWithUserId(userId uuid.UUID, pickUpPointId uuid.UUID) *modelsb.StorekeeperBuilder {
	return StorekeeperRandomId().
		WithUserId(userId).
		WithPickUpPointId(pickUpPointId)
}

