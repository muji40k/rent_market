package models

import (
	"math/rand/v2"
	modelsb "rent_service/builders/domain/models"
	"rent_service/builders/misc/uuidgen"
	"rent_service/misc/nullable"
	"time"

	"github.com/google/uuid"
)

func ReviewRandomId() *modelsb.ReviewBuilder {
	return modelsb.NewReview().
		WithId(uuidgen.Generate())
}

func ReviewExample(
	prefix string,
	instanceId uuid.UUID,
	userId uuid.UUID,
	rating *nullable.Nullable[float64],
	date *nullable.Nullable[time.Time],
) *modelsb.ReviewBuilder {
	return ReviewRandomId().
		WithInstanceId(instanceId).
		WithUserId(userId).
		WithContent("Review " + prefix + " content for tests").
		WithRating(nullable.GetOrFunc(rating, func() float64 {
			return 1 + 4*rand.Float64()
		})).
		WithDate(nullable.GetOrFunc(date, time.Now))
}

