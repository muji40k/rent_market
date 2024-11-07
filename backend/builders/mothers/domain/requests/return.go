package requests

import (
	requestsb "rent_service/builders/domain/requests"
	"rent_service/builders/misc/uuidgen"
	"rent_service/misc/nullable"
	"time"

	"github.com/google/uuid"
)

func ReturnRandomId() *requestsb.ReturnBuilder {
	return requestsb.NewReturn().
		WithId(uuidgen.Generate())
}

func Return(
	instanceId uuid.UUID,
	userId uuid.UUID,
	pickUpPointId uuid.UUID,
	rentEndDate *nullable.Nullable[time.Time],
	verificationCode *nullable.Nullable[string],
	createDate *nullable.Nullable[time.Time],
) *requestsb.ReturnBuilder {
	return ReturnRandomId().
		WithInstanceId(instanceId).
		WithUserId(userId).
		WithPickUpPointId(pickUpPointId).
		WithRentEndDate(nullable.GetOrFunc(rentEndDate, time.Now)).
		WithVerificationCode(nullable.GetOrFunc(verificationCode, GetCode)).
		WithCreateDate(nullable.GetOrFunc(createDate, time.Now))
}

