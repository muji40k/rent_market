package records

import (
	recordsb "rent_service/builders/domain/records"
	"rent_service/builders/misc/dategen"
	"rent_service/builders/misc/uuidgen"
	"rent_service/misc/nullable"
	"time"

	"github.com/google/uuid"
)

func StorageRandomId() *recordsb.StorageBuilder {
	return recordsb.NewStorage().
		WithId(uuidgen.Generate())
}

func StorageActive(
	pickUpPointId uuid.UUID,
	instanceId uuid.UUID,
	inDate *nullable.Nullable[time.Time],
) *recordsb.StorageBuilder {
	return StorageRandomId().
		WithInstanceId(instanceId).
		WithPickUpPointId(pickUpPointId).
		WithInDate(nullable.GetOrFunc(inDate, getStartDate))
}

func StorageFinished(
	pickUpPointId uuid.UUID,
	instanceId uuid.UUID,
	inDate *nullable.Nullable[time.Time],
	outDate *nullable.Nullable[time.Time],
) *recordsb.StorageBuilder {
	in := nullable.GetOrFunc(inDate, getStartDate)
	return StorageRandomId().
		WithInstanceId(instanceId).
		WithPickUpPointId(pickUpPointId).
		WithInDate(in).
		WithOutDate(nullable.OrFunc(outDate, func() *nullable.Nullable[time.Time] {
			return nullable.Some(dategen.GetDate(
				dategen.FromTime(in),
				dategen.FromTime(time.Now()),
			))
		}))
}

