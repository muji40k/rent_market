package models

import (
	modelsb "rent_service/builders/domain/models"
	"rent_service/builders/misc/uuidgen"
)

func DeliveryCompanyRandomId() *modelsb.DeliveryCompanyBuilder {
	return modelsb.NewDeliveryCompany().
		WithId(uuidgen.Generate())
}

func DeliveryCompanyExample(prefix string) *modelsb.DeliveryCompanyBuilder {
	return DeliveryCompanyRandomId().
		WithName("Example Company " + prefix).
		WithSite("https://example." + prefix + ".delivery.company.org").
		WithPhoneNumber("88005553535").
		WithDescription("Example delivery company for tests")
}

