package requests

import (
	requestsb "rent_service/builders/domain/requests"
	"rent_service/builders/misc/uuidgen"
	"rent_service/internal/domain/models"
	"rent_service/misc/nullable"
	"time"

	"github.com/google/uuid"
)

func ProvideRandomId() *requestsb.ProvideBuilder {
	return requestsb.NewProvide().
		WithId(uuidgen.Generate())
}

func ProvideExample(
	prefix string,
	productId uuid.UUID,
	renterId uuid.UUID,
	pickUpPointId uuid.UUID,
	verificationCode *nullable.Nullable[string],
	createDate *nullable.Nullable[time.Time],
	plans ...models.PayPlan,
) *requestsb.ProvideBuilder {
	return ProvideRandomId().
		WithProductId(productId).
		WithRenterId(renterId).
		WithPickUpPointId(pickUpPointId).
		WithPayPlans(plans...).
		WithName("Example provision " + prefix).
		WithDescription("Example provision for tests").
		WithCondition("Condition " + prefix).
		WithVerificationCode(nullable.GetOrFunc(verificationCode, GetCode)).
		WithCreateDate(nullable.GetOrFunc(createDate, time.Now))
}

