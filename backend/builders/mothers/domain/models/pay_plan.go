package models

import (
	"math/rand/v2"
	modelsb "rent_service/builders/domain/models"
	"rent_service/builders/misc/uuidgen"
	mcurrency "rent_service/builders/mothers/currency"
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/misc/types/currency"
	"sort"

	"github.com/google/uuid"
)

func PayPlanRandomId() *modelsb.PayPlanBuilder {
	return modelsb.NewPayPlan().
		WithId(uuidgen.Generate())
}

func PayPlanWithPeriodIdAndPrice(
	periodId uuid.UUID,
	price currency.Currency,
) *modelsb.PayPlanBuilder {
	return PayPlanRandomId().
		WithPeriodId(periodId).
		WithPrice(price)
}

func PayPlanCollect(builders ...*modelsb.PayPlanBuilder) []models.PayPlan {
	return collection.Collect(
		collection.MapIterator(
			func(builder **modelsb.PayPlanBuilder) models.PayPlan {
				return (*builder).Build()
			},
			collection.SliceIterator(builders),
		),
	)
}

func PayPlansWithPeriods(periods ...models.Period) []*modelsb.PayPlanBuilder {
	scale := 0.3 + 4.7*rand.Float64()
	sort.Slice(periods, func(i, j int) bool {
		return periods[i].Duration < periods[j].Duration
	})

	var price float64 = 1000
	return collection.Collect(
		collection.MapIterator(
			func(period *models.Period) *modelsb.PayPlanBuilder {
				p := scale * price
				price += 1000
				return PayPlanWithPeriodIdAndPrice(
					period.Id,
					mcurrency.RUB(p).Build(),
				)
			},
			collection.SliceIterator(periods),
		),
	)
}

