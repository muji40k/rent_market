package errstructs

import "rent_service/internal/logic/services/errors/cmnerrors"

type BadRequest struct {
	Reason string `json:"error"`
}

func NewBadRequest(reason string) BadRequest {
	return BadRequest{reason}
}

func NewBadRequestErr(err error) BadRequest {
	return BadRequest{err.Error()}
}

type Internal struct {
	Error string `json:"error"`
}

func NewInternal(error string) Internal {
	return Internal{error}
}

func NewInternalErr(err error) Internal {
	return Internal{err.Error()}
}

type NotFound struct {
	What []string `json:"what"`
}

func NewNotFound(err cmnerrors.ErrorNotFound) NotFound {
	return NotFound{err.What}
}

