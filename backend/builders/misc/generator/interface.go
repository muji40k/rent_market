package generator

import "github.com/google/uuid"

type IGenerator interface {
	Generate() uuid.UUID
	Finish()
}

