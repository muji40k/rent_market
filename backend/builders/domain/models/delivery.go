package models

import (
	"rent_service/internal/domain/models"

	"github.com/google/uuid"
)

type DeliveryCompanyBuilder struct {
	id          uuid.UUID
	name        string
	site        string
	phoneNumber string
	description string
}

func NewDeliveryCompany() *DeliveryCompanyBuilder {
	return &DeliveryCompanyBuilder{}
}

func (self *DeliveryCompanyBuilder) WithId(id uuid.UUID) *DeliveryCompanyBuilder {
	self.id = id
	return self
}

func (self *DeliveryCompanyBuilder) WithName(name string) *DeliveryCompanyBuilder {
	self.name = name
	return self
}

func (self *DeliveryCompanyBuilder) WithSite(site string) *DeliveryCompanyBuilder {
	self.site = site
	return self
}

func (self *DeliveryCompanyBuilder) WithPhoneNumber(phoneNumber string) *DeliveryCompanyBuilder {
	self.phoneNumber = phoneNumber
	return self
}

func (self *DeliveryCompanyBuilder) WithDescription(description string) *DeliveryCompanyBuilder {
	self.description = description
	return self
}

func (self *DeliveryCompanyBuilder) Build() models.DeliveryCompany {
	return models.DeliveryCompany{
		Id:          self.id,
		Name:        self.name,
		Site:        self.site,
		PhoneNumber: self.phoneNumber,
		Description: self.description,
	}
}

