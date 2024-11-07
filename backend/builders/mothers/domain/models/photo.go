package models

import (
	modelsb "rent_service/builders/domain/models"
	"rent_service/builders/misc/uuidgen"
	"rent_service/misc/nullable"
	"time"
)

func PhotoRandomId() *modelsb.PhotoBuilder {
	return modelsb.NewPhoto().
		WithId(uuidgen.Generate())
}

func PhotoExample(
	prefix string,
	date *nullable.Nullable[time.Time],
) *modelsb.PhotoBuilder {
	return PhotoRandomId().
		WithPath("/path/to/" + prefix).
		WithMime("image/png").
		WithPlaceholder(prefix).
		WithDescription("Description for " + prefix).
		WithDate(nullable.GetOrFunc(date, time.Now))
}

func TempPhotoRandomId() *modelsb.TempPhotoBuilder {
	return modelsb.NewTempPhoto().
		WithId(uuidgen.Generate())
}

func TempPhotoExampleEmpty(
	prefix string,
	date *nullable.Nullable[time.Time],
) *modelsb.TempPhotoBuilder {
	return TempPhotoRandomId().
		WithMime("image/png").
		WithPlaceholder(prefix).
		WithDescription("Description for " + prefix).
		WithCreate(nullable.GetOrFunc(date, time.Now))
}

func TempPhotoExampleUploaded(
	prefix string,
	date *nullable.Nullable[time.Time],
) *modelsb.TempPhotoBuilder {
	return TempPhotoExampleEmpty(prefix, date).
		WithPath(nullable.Some("/path/to/" + prefix))
}

