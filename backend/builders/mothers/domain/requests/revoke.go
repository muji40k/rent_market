package requests

import (
	requestsb "rent_service/builders/domain/requests"
	"rent_service/builders/misc/uuidgen"
	"rent_service/misc/nullable"
	"time"

	"github.com/google/uuid"
)

func RevokeRandomId() *requestsb.RevokeBuilder {
	return requestsb.NewRevoke().
		WithId(uuidgen.Generate())
}

func Revoke(
	instanceId uuid.UUID,
	renterId uuid.UUID,
	pickUpPointId uuid.UUID,
	verificationCode *nullable.Nullable[string],
	createDate *nullable.Nullable[time.Time],
) *requestsb.RevokeBuilder {
	return RevokeRandomId().
		WithInstanceId(instanceId).
		WithRenterId(renterId).
		WithPickUpPointId(pickUpPointId).
		WithVerificationCode(nullable.GetOrFunc(verificationCode, GetCode)).
		WithCreateDate(nullable.GetOrFunc(createDate, time.Now))
}

