package requests

import (
	requestsb "rent_service/builders/domain/requests"
	"rent_service/builders/misc/uuidgen"
	"rent_service/misc/nullable"
	"time"

	"github.com/google/uuid"
)

func RentRandomId() *requestsb.RentBuilder {
	return requestsb.NewRent().
		WithId(uuidgen.Generate())
}

func Rent(
	instanceId uuid.UUID,
	userId uuid.UUID,
	pickUpPointId uuid.UUID,
	paymentPeriodId uuid.UUID,
	verificationCode *nullable.Nullable[string],
	createDate *nullable.Nullable[time.Time],
) *requestsb.RentBuilder {
	return RentRandomId().
		WithInstanceId(instanceId).
		WithUserId(userId).
		WithPickUpPointId(pickUpPointId).
		WithPaymentPeriodId(paymentPeriodId).
		WithVerificationCode(nullable.GetOrFunc(verificationCode, GetCode)).
		WithCreateDate(nullable.GetOrFunc(createDate, time.Now))
}

