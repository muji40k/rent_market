package requests

import (
	requestsb "rent_service/builders/domain/requests"
	"rent_service/builders/misc/dategen"
	"rent_service/builders/misc/uuidgen"
	"rent_service/misc/nullable"
	"time"

	"github.com/google/uuid"
)

func DeliveryRandomId() *requestsb.DeliveryBuilder {
	return requestsb.NewDelivery().
		WithId(uuidgen.Generate())
}

func baseDelivery(
	companyId uuid.UUID,
	instanceId uuid.UUID,
	fromId uuid.UUID,
	toId uuid.UUID,
	deliveryId *nullable.Nullable[string],
	verificationCode *nullable.Nullable[string],
) *requestsb.DeliveryBuilder {
	return DeliveryRandomId().
		WithCompanyId(companyId).
		WithInstanceId(instanceId).
		WithFromId(fromId).
		WithToId(toId).
		WithDeliveryId(nullable.GetOrFunc(deliveryId, func() string {
			return companyId.String() + time.Now().String()
		})).
		WithVerificationCode(nullable.GetOrFunc(verificationCode, GetCode))
}

func DeliveryRandomCreated(
	companyId uuid.UUID,
	instanceId uuid.UUID,
	fromId uuid.UUID,
	toId uuid.UUID,
	deliveryId *nullable.Nullable[string],
	sBeginDate *nullable.Nullable[time.Time],
	sEndDate *nullable.Nullable[time.Time],
	verificationCode *nullable.Nullable[string],
	createDate *nullable.Nullable[time.Time],
) *requestsb.DeliveryBuilder {
	begin := nullable.GetOrFunc(sBeginDate, getStartDate)
	return baseDelivery(companyId, instanceId, fromId, toId, deliveryId, verificationCode).
		WithScheduledBeginDate(begin).
		WithScheduledEndDate(nullable.GetOrFunc(sEndDate, func() time.Time {
			return dategen.GetDate(
				dategen.FromTime(begin),
				dategen.FromTime(begin.Add(14*24*time.Hour)),
			)
		})).
		WithCreateDate(nullable.GetOrFunc(createDate, func() time.Time {
			return dategen.GetDate(
				dategen.FromTime(begin.Add(-14*24*time.Hour)),
				dategen.FromTime(begin),
			)
		}))
}

func DeliveryRandomSent(
	companyId uuid.UUID,
	instanceId uuid.UUID,
	fromId uuid.UUID,
	toId uuid.UUID,
	deliveryId *nullable.Nullable[string],
	sBeginDate *nullable.Nullable[time.Time],
	aBeginDate *nullable.Nullable[time.Time],
	sEndDate *nullable.Nullable[time.Time],
	verificationCode *nullable.Nullable[string],
	createDate *nullable.Nullable[time.Time],
) *requestsb.DeliveryBuilder {
	begin := nullable.GetOrFunc(sBeginDate, getStartDate)
	end := nullable.GetOrFunc(sEndDate, func() time.Time {
		return dategen.GetDate(
			dategen.FromTime(begin),
			dategen.FromTime(begin.Add(14*24*time.Hour)),
		)
	})
	return baseDelivery(companyId, instanceId, fromId, toId, deliveryId, verificationCode).
		WithScheduledBeginDate(begin).
		WithActualBeginDate(nullable.OrFunc(aBeginDate, func() *nullable.Nullable[time.Time] {
			m := end
			now := time.Now()

			if m.After(now) {
				m = now
			}

			return nullable.Some(dategen.GetDate(
				dategen.FromTime(begin),
				dategen.FromTime(m),
			))
		})).
		WithScheduledEndDate(end).
		WithCreateDate(nullable.GetOrFunc(createDate, func() time.Time {
			return dategen.GetDate(
				dategen.FromTime(begin.Add(-14*24*time.Hour)),
				dategen.FromTime(begin),
			)
		}))
}

func DeliveryRandomAccepted(
	companyId uuid.UUID,
	instanceId uuid.UUID,
	fromId uuid.UUID,
	toId uuid.UUID,
	deliveryId *nullable.Nullable[string],
	sBeginDate *nullable.Nullable[time.Time],
	aBeginDate *nullable.Nullable[time.Time],
	sEndDate *nullable.Nullable[time.Time],
	aEndDate *nullable.Nullable[time.Time],
	verificationCode *nullable.Nullable[string],
	createDate *nullable.Nullable[time.Time],
) *requestsb.DeliveryBuilder {
	begin := nullable.GetOrFunc(sBeginDate, getStartDate)
	end := nullable.GetOrFunc(sEndDate, func() time.Time {
		return dategen.GetDate(
			dategen.FromTime(begin),
			dategen.FromTime(begin.Add(14*24*time.Hour)),
		)
	})
	return baseDelivery(companyId, instanceId, fromId, toId, deliveryId, verificationCode).
		WithScheduledBeginDate(begin).
		WithActualBeginDate(nullable.OrFunc(aBeginDate, func() *nullable.Nullable[time.Time] {
			m := end
			now := time.Now()

			if m.After(now) {
				m = now
			}

			return nullable.Some(dategen.GetDate(
				dategen.FromTime(begin),
				dategen.FromTime(m),
			))
		})).
		WithScheduledEndDate(end).
		WithActualEndDate(nullable.OrFunc(aEndDate, func() *nullable.Nullable[time.Time] {
			m := end.Add(14 * 24 * time.Hour)
			now := time.Now()

			if m.After(now) {
				m = now
			}

			return nullable.Some(dategen.GetDate(
				dategen.FromTime(end),
				dategen.FromTime(m),
			))
		})).
		WithCreateDate(nullable.GetOrFunc(createDate, func() time.Time {
			return dategen.GetDate(
				dategen.FromTime(begin.Add(-14*24*time.Hour)),
				dategen.FromTime(begin),
			)
		}))
}

