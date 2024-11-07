package uuidgen

import "github.com/google/uuid"

func Generate() uuid.UUID {
	id, err := uuid.NewRandom()

	if nil != err {
		panic(err)
	}

	return id
}

