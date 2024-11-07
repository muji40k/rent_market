package records

import (
	recordsb "rent_service/builders/domain/records"
	"rent_service/builders/misc/dategen"
	"rent_service/builders/misc/uuidgen"
	"rent_service/misc/nullable"
	"time"

	"github.com/google/uuid"
)

func RentRandomId() *recordsb.RentBuilder {
	return recordsb.NewRent().
		WithId(uuidgen.Generate())
}

func RentActive(
	userId uuid.UUID,
	instanceId uuid.UUID,
	startDate *nullable.Nullable[time.Time],
	paymentPeriodId uuid.UUID,
) *recordsb.RentBuilder {
	return RentRandomId().
		WithUserId(userId).
		WithInstanceId(instanceId).
		WithStartDate(nullable.GetOrFunc(startDate, getStartDate)).
		WithPaymentPeriodId(paymentPeriodId)
}

func RentReturned(
	userId uuid.UUID,
	instanceId uuid.UUID,
	startDate *nullable.Nullable[time.Time],
	endDate *nullable.Nullable[time.Time],
	paymentPeriodId uuid.UUID,
) *recordsb.RentBuilder {
	start := nullable.GetOrFunc(startDate, getStartDate)
	return RentRandomId().
		WithUserId(userId).
		WithInstanceId(instanceId).
		WithStartDate(start).
		WithEndDate(nullable.OrFunc(endDate, func() *nullable.Nullable[time.Time] {
			return nullable.Some(dategen.GetDate(
				dategen.FromTime(start),
				dategen.FromTime(time.Now()),
			))
		})).
		WithPaymentPeriodId(paymentPeriodId)
}

