package models

import (
	"rent_service/internal/domain/models"

	"github.com/google/uuid"
)

type ProductBuilder struct {
	id          uuid.UUID
	name        string
	categoryId  uuid.UUID
	description string
}

func NewProduct() *ProductBuilder {
	return &ProductBuilder{}
}

func (self *ProductBuilder) WithId(id uuid.UUID) *ProductBuilder {
	self.id = id
	return self
}

func (self *ProductBuilder) WithName(name string) *ProductBuilder {
	self.name = name
	return self
}

func (self *ProductBuilder) WithCategoryId(categoryId uuid.UUID) *ProductBuilder {
	self.categoryId = categoryId
	return self
}

func (self *ProductBuilder) WithDescription(description string) *ProductBuilder {
	self.description = description
	return self
}

func (self *ProductBuilder) Build() models.Product {
	return models.Product{
		Id:          self.id,
		Name:        self.name,
		CategoryId:  self.categoryId,
		Description: self.description,
	}
}

type CharachteristicBuilder struct {
	id    uuid.UUID
	name  string
	value string
}

func NewCharachteristic() *CharachteristicBuilder {
	return &CharachteristicBuilder{}
}

func (self *CharachteristicBuilder) WithId(id uuid.UUID) *CharachteristicBuilder {
	self.id = id
	return self
}

func (self *CharachteristicBuilder) WithName(name string) *CharachteristicBuilder {
	self.name = name
	return self
}

func (self *CharachteristicBuilder) WithValue(value string) *CharachteristicBuilder {
	self.value = value
	return self
}

func (self *CharachteristicBuilder) Build() models.Charachteristic {
	return models.Charachteristic{
		Id:    self.id,
		Name:  self.name,
		Value: self.value,
	}
}

type ProductCharacteristicsBuilder struct {
	productId uuid.UUID
	mmap      map[string]models.Charachteristic
}

func NewProductCharacteristics() *ProductCharacteristicsBuilder {
	return &ProductCharacteristicsBuilder{}
}

func (self *ProductCharacteristicsBuilder) WithProductId(productId uuid.UUID) *ProductCharacteristicsBuilder {
	self.productId = productId
	return self
}

func (self *ProductCharacteristicsBuilder) WithMmap(chars ...models.Charachteristic) *ProductCharacteristicsBuilder {
	self.mmap = make(map[string]models.Charachteristic)

	for _, c := range chars {
		self.mmap[c.Name] = c
	}

	return self
}

func (self *ProductCharacteristicsBuilder) Build() models.ProductCharacteristics {
	return models.ProductCharacteristics{
		ProductId: self.productId,
		Map:       self.mmap,
	}
}

