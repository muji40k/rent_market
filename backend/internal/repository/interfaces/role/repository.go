package role

import (
	"rent_service/internal/domain/models"

	"github.com/google/uuid"
)

type IAdministratorRepository interface {
	GetByUserId(userId uuid.UUID) (models.Administrator, error)
}

type IRenterRepository interface {
	Create(userId uuid.UUID) (models.Renter, error)

	GetById(renterId uuid.UUID) (models.Renter, error)
	GetByUserId(userId uuid.UUID) (models.Renter, error)
}

type IStorekeeperRepository interface {
	GetByUserId(userId uuid.UUID) (models.Storekeeper, error)
}
