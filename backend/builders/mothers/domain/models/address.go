package models

import (
	modelsb "rent_service/builders/domain/models"
	"rent_service/builders/misc/uuidgen"
	"rent_service/misc/nullable"
)

func AddressRandomId() *modelsb.AddressBuilder {
	return modelsb.NewAddress().
		WithId(uuidgen.Generate())
}

func AddressExmapleWithoutFlat(prefix string) *modelsb.AddressBuilder {
	return AddressRandomId().
		WithCountry("Country " + prefix).
		WithCity("City " + prefix).
		WithStreet("Street " + prefix).
		WithHouse("House " + prefix)
}

func AddressExmapleWithFlat(prefix string) *modelsb.AddressBuilder {
	return AddressExmapleWithoutFlat(prefix).
		WithFlat(nullable.Some("#" + prefix))
}

