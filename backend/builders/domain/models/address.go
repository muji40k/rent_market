package models

import (
	"rent_service/builders/misc/nullcommon"
	"rent_service/internal/domain/models"
	"rent_service/misc/nullable"

	"github.com/google/uuid"
)

type AddressBuilder struct {
	id      uuid.UUID
	country string
	city    string
	street  string
	house   string
	flat    *string
}

func NewAddress() *AddressBuilder {
	return &AddressBuilder{}
}

func (self *AddressBuilder) WithId(id uuid.UUID) *AddressBuilder {
	self.id = id
	return self
}

func (self *AddressBuilder) WithCountry(country string) *AddressBuilder {
	self.country = country
	return self
}

func (self *AddressBuilder) WithCity(city string) *AddressBuilder {
	self.city = city
	return self
}

func (self *AddressBuilder) WithStreet(street string) *AddressBuilder {
	self.street = street
	return self
}

func (self *AddressBuilder) WithHouse(house string) *AddressBuilder {
	self.house = house
	return self
}

func (self *AddressBuilder) WithFlat(flat *nullable.Nullable[string]) *AddressBuilder {
	self.flat = nullcommon.CopyPtrIfSome(flat)
	return self
}

func (self *AddressBuilder) Build() models.Address {
	return models.Address{
		Id:      self.id,
		Country: self.country,
		City:    self.city,
		Street:  self.street,
		House:   self.house,
		Flat:    self.flat,
	}
}

