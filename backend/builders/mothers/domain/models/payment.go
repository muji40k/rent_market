package models

import (
	modelsb "rent_service/builders/domain/models"
	"rent_service/builders/misc/uuidgen"
	"rent_service/internal/misc/types/currency"
	"rent_service/misc/nullable"
	"time"

	"github.com/google/uuid"
)

func PaymentRandomId() *modelsb.PaymentBuilder {
	return modelsb.NewPayment().
		WithId(uuidgen.Generate())
}

func PaymentExampleIncomplete(
	rentId uuid.UUID,
	periodStart time.Time,
	periodEnd time.Time,
	value currency.Currency,
	createDate *nullable.Nullable[time.Time],
) *modelsb.PaymentBuilder {
	return PaymentRandomId().
		WithRentId(rentId).
		WithPeriodStart(periodStart).
		WithPeriodEnd(periodEnd).
		WithValue(value).
		WithStatus("created").
		WithCreateDate(nullable.GetOrFunc(createDate, time.Now))
}

func PaymentExampleCompleted(
	rentId uuid.UUID,
	payMethodId uuid.UUID,
	paymentId string,
	periodStart time.Time,
	periodEnd time.Time,
	value currency.Currency,
	paymentDate *nullable.Nullable[time.Time],
	createDate *nullable.Nullable[time.Time],
) *modelsb.PaymentBuilder {
	return PaymentExampleIncomplete(rentId, periodStart, periodEnd, value, createDate).
		WithPayMethodId(nullable.Some(payMethodId)).
		WithPaymentId(nullable.Some(paymentId)).
		WithStatus("completed").
		WithPaymentDate(
			nullable.OrFunc(paymentDate, func() *nullable.Nullable[time.Time] {
				return nullable.Some(time.Now())
			}),
		)
}

