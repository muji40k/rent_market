package models

import (
	modelsb "rent_service/builders/domain/models"
	"rent_service/internal/domain/models"

	"github.com/google/uuid"
)

func CategoryRandomId() *modelsb.CategoryBuilder {
	id, err := uuid.NewRandom()

	if nil != err {
		panic(err)
	}

	return modelsb.NewCategory().
		WithId(id)
}

func CategoryWithParentRandomId(parent models.Category) *modelsb.CategoryBuilder {
	id, err := uuid.NewRandom()

	if nil != err {
		panic(err)
	}

	return modelsb.NewCategory().
		WithId(id).
		WithParentId(&parent.Id)
}

func CategoryToPath(path []*modelsb.CategoryBuilder) []models.Category {
	out := make([]models.Category, len(path))

	for i, c := range path {
		out[i] = c.Build()
	}

	return out
}

func CategoryPath(path ...string) []*modelsb.CategoryBuilder {
	if 0 == len(path) {
		return nil
	}

	out := make([]*modelsb.CategoryBuilder, len(path))
	previous := CategoryRandomId().WithName(path[0])
	out[0] = previous

	for i, name := range path[1:] {
		previous = CategoryWithParentRandomId(previous.Build()).WithName(name)
		out[i+1] = previous
	}

	return out
}

func CategoryDefaultPath() []*modelsb.CategoryBuilder {
	return CategoryPath("root", "level 1", "level 2")
}

