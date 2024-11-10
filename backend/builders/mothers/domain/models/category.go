package models

import (
	modelsb "rent_service/builders/domain/models"
	"rent_service/builders/misc/uuidgen"
	"rent_service/misc/nullable"

	"github.com/google/uuid"
)

func CategoryRandomId() *modelsb.CategoryBuilder {
	return modelsb.NewCategory().
		WithId(uuidgen.Generate())
}

func CategoryWithParentId(parentId uuid.UUID) *modelsb.CategoryBuilder {
	return CategoryRandomId().
		WithParentId(nullable.Some(parentId))
}

func CategoryPath(path ...string) []*modelsb.CategoryBuilder {
	if 0 == len(path) {
		return nil
	}

	out := make([]*modelsb.CategoryBuilder, len(path))
	previous := CategoryRandomId().WithName(path[0])
	out[0] = previous

	for i, name := range path[1:] {
		previous = CategoryWithParentId(previous.Build().Id).WithName(name)
		out[i+1] = previous
	}

	return out
}

func CategoryDefaultPath() []*modelsb.CategoryBuilder {
	return CategoryPath("root", "level 1", "level 2")
}

