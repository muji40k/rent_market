package models

import (
	modelsb "rent_service/builders/domain/models"
	"rent_service/builders/misc/uuidgen"
	"time"
)

func PeriodRandomId() *modelsb.PeriodBuilder {
	return modelsb.NewPeriod().
		WithId(uuidgen.Generate())
}

func PeriodDay() *modelsb.PeriodBuilder {
	return PeriodRandomId().
		WithName("day").
		WithDuration(24 * time.Hour)
}

func PeriodWeek() *modelsb.PeriodBuilder {
	return PeriodRandomId().
		WithName("week").
		WithDuration(7 * 24 * time.Hour)
}

func PeriodMonth() *modelsb.PeriodBuilder {
	return PeriodRandomId().
		WithName("month").
		WithDuration(30 * 24 * time.Hour)
}

func PeriodQuarter() *modelsb.PeriodBuilder {
	return PeriodRandomId().
		WithName("quarter").
		WithDuration(90 * 24 * time.Hour)
}

func PeriodHalf() *modelsb.PeriodBuilder {
	return PeriodRandomId().
		WithName("half").
		WithDuration(180 * 24 * time.Hour)
}

func PeriodYear() *modelsb.PeriodBuilder {
	return PeriodRandomId().
		WithName("year").
		WithDuration(360 * 24 * time.Hour)
}

