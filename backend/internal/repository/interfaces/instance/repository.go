package instance

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

//go:generate mockgen -source=repository.go -destination=../../implementation/mock/instance/repository.go

type Filter struct {
	ProductId uuid.UUID
}

type Sort uint

const (
	SORT_NONE Sort = iota // Must be consistent
	SORT_RATING_ASC
	SORT_RATING_DSC
	SORT_DATE_ASC
	SORT_DATE_DSC
	SORT_PRICE_ASC
	SORT_PRICE_DSC
	SORT_USAGE_ASC
	SORT_USAGE_DSC
)

type IRepository interface {
	Create(instance models.Instance) (models.Instance, error)

	Update(instance models.Instance) error

	GetById(instanceId uuid.UUID) (models.Instance, error)
	GetWithFilter(
		filter Filter,
		sort Sort,
	) (collection.Collection[models.Instance], error)
}

type IPayPlansRepository interface {
	Create(payPlans models.InstancePayPlans) (models.InstancePayPlans, error)
	AddPayPlan(
		instanceId uuid.UUID,
		plan models.PayPlan,
	) (models.InstancePayPlans, error)

	Update(models.InstancePayPlans) error

	GetByInstanceId(instanceId uuid.UUID) (models.InstancePayPlans, error)
}

type IPhotoRepository interface {
	Create(instanceId uuid.UUID, photoId uuid.UUID) error
	GetByInstanceId(instanceId uuid.UUID) (collection.Collection[uuid.UUID], error)
}

