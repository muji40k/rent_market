package models

import (
	"math/rand/v2"
	modelsb "rent_service/builders/domain/models"
	"rent_service/builders/misc/uuidgen"
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	"time"

	"github.com/google/uuid"
)

func PickUpPointRandomId() *modelsb.PickUpPointBuilder {
	return modelsb.NewPickUpPoint().
		WithId(uuidgen.Generate())
}

func PickUpPointExample(prefix string) *modelsb.PickUpPointBuilder {
	return PickUpPointRandomId().
		WithCapacity(rand.Uint64()).
		WithAddress(AddressExmapleWithoutFlat(prefix).Build())
}

func WorkingHoursRandomId() *modelsb.WorkingHoursBuilder {
	return modelsb.NewWorkingHours().
		WithId(uuidgen.Generate())
}

func WorkingHoursExample(day time.Weekday, start time.Duration, end time.Duration) *modelsb.WorkingHoursBuilder {
	return WorkingHoursRandomId().
		WithDay(day).
		WithBegin(start).
		WithEnd(end)
}

func WorkingHoursWeek(start time.Duration, end time.Duration) []*modelsb.WorkingHoursBuilder {
	base := int(time.Monday)
	return collection.Collect(
		collection.MapIterator(
			func(i *int) *modelsb.WorkingHoursBuilder {
				return WorkingHoursExample(time.Weekday(base+*i), start, end)
			},
			collection.RangeIterator(collection.RangeEnd(5)),
		),
	)
}

func PickUpPointWorkingHours(pickUpPointId uuid.UUID, wh ...models.WorkingHours) *modelsb.PickUpPointWorkingHoursBuilder {
	return modelsb.NewPickUpPointWorkingHours().
		WithPickUpPointId(pickUpPointId).
		WithWorkingHours(wh...)
}

