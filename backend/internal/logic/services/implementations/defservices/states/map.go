package states

import (
	"errors"
	"rent_service/internal/logic/services/errors/cmnerrors"
)

func MapError(err error) error {
	if cerr := (ErrorForbiddenMethod{}); errors.As(err, &cerr) {
		err = cmnerrors.Conflict(cerr.Error())
	} else if errors.Is(err, ErrorBrokenState) {
		err = cmnerrors.Conflict(err.Error())
	} else if errors.Is(err, ErrorUnknownInstance) {
		err = cmnerrors.Incorrect("instance_id")
	}

	return err
}

