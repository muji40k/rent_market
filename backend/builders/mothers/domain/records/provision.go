package records

import (
	recordsb "rent_service/builders/domain/records"
	"rent_service/builders/misc/dategen"
	"rent_service/builders/misc/uuidgen"
	"rent_service/misc/nullable"
	"time"

	"github.com/google/uuid"
)

func ProvisionRandomId() *recordsb.ProvisionBuilder {
	return recordsb.NewProvision().
		WithId(uuidgen.Generate())
}

func ProvisionActive(renterId uuid.UUID, instanceId uuid.UUID, startDate *nullable.Nullable[time.Time]) *recordsb.ProvisionBuilder {
	return ProvisionRandomId().
		WithRenterId(renterId).
		WithInstanceId(instanceId).
		WithStartDate(nullable.GetOrFunc(startDate, getStartDate))
}

func ProvisionRevoked(renterId uuid.UUID, instanceId uuid.UUID, startDate *nullable.Nullable[time.Time], endDate *nullable.Nullable[time.Time]) *recordsb.ProvisionBuilder {
	start := nullable.GetOrFunc(startDate, getStartDate)
	return ProvisionRandomId().
		WithRenterId(renterId).
		WithInstanceId(instanceId).
		WithStartDate(start).
		WithEndDate(nullable.OrFunc(endDate, func() *nullable.Nullable[time.Time] {
			return nullable.Some(dategen.GetDate(
				dategen.FromTime(start),
				dategen.FromTime(time.Now()),
			))
		}))
}

