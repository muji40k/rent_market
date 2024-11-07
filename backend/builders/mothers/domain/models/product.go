package models

import (
	"fmt"
	modelsb "rent_service/builders/domain/models"
	"rent_service/builders/misc/uuidgen"
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

func ProductRandomId() *modelsb.ProductBuilder {
	return modelsb.NewProduct().
		WithId(uuidgen.Generate())
}

func ProductExmaple(prefix string, categoryId uuid.UUID) *modelsb.ProductBuilder {
	return ProductRandomId().
		WithName("Example " + prefix).
		WithCategoryId(categoryId).
		WithDescription("Product description for tests")
}

func CharacteristicRandomId() *modelsb.CharachteristicBuilder {
	return modelsb.NewCharachteristic().
		WithId(uuidgen.Generate())
}

func CharacteristicExample(name string, value string) *modelsb.CharachteristicBuilder {
	return CharacteristicRandomId().
		WithName(name).
		WithValue(value)
}

func CharacteristicExampleNumeric(name string, value float64) *modelsb.CharachteristicBuilder {
	return CharacteristicRandomId().
		WithName(name).
		WithValue(fmt.Sprint(value))
}

func CharacteristicCollect(builders ...*modelsb.CharachteristicBuilder) []models.Charachteristic {
	return collection.Collect(collection.MapIterator(
		func(builder **modelsb.CharachteristicBuilder) models.Charachteristic {
			return (*builder).Build()
		},
		collection.SliceIterator(builders),
	))
}

func ProductCharacteristics(productId uuid.UUID, chars ...models.Charachteristic) *modelsb.ProductCharacteristicsBuilder {
	return modelsb.NewProductCharacteristics().
		WithProductId(productId).
		WithCharacteristics(chars...)
}

