package models

import (
	"rent_service/internal/domain/models"

	"github.com/google/uuid"
)

type Product struct {
	ID          string   `gorm:"column:id;type:uuid"`
	Name        string   `gorm:"column:name"`
	CategoryId  string   `gorm:"column:category_id;type:uuid"`
	Category    Category `gorm:"references:ID;foreignKey:CategoryId"`
	Description string   `gorm:"column:description"`
	// TSVector    []string `gorm:"column:ts;type:tsvector"`
	Meta
}

func (Product) TableName() string {
	return "products.products"
}

func MapProduct(value *Product) models.Product {
	return models.Product{
		Id:          uuid.MustParse(value.ID),
		Name:        value.Name,
		CategoryId:  uuid.MustParse(value.CategoryId),
		Description: value.Description,
	}
}

type ProductCharacteristic struct {
	ID        string  `gorm:"column:id;type:uuid"`
	ProductId string  `gorm:"column:product_id;type:uuid"`
	Product   Product `gorm:"references:ID;foreignKey:ProductId"`
	Name      string  `gorm:"column:name"`
	Value     string  `gorm:"column:value"`
	Meta
}

func (ProductCharacteristic) TableName() string {
	return "products.characteristics"
}

