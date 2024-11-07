package models

import (
	modelsb "rent_service/builders/domain/models"
	"rent_service/builders/misc/uuidgen"
)

func PayMethodRandomId() *modelsb.PayMethodBuilder {
	return modelsb.NewPayMethod().
		WithId(uuidgen.Generate())
}

func PayMethodExample(prefix string) *modelsb.PayMethodBuilder {
	return PayMethodRandomId().
		WithName("Example Pay Method " + prefix).
		WithDescription("Example pay method for tests")
}

