package defstates

import (
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/records"
	"rent_service/internal/logic/services/errors/cmnerrors"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
	"time"

	"github.com/google/uuid"
)

func (self *InstanceStateMachine) actionCreateStorage(
	instance models.Instance,
	pickUpPointId uuid.UUID,
) (records.Storage, error) {
	repo := self.repos.storage.GetStorageRepository()
	storage, err := repo.Create(records.Storage{
		PickUpPointId: pickUpPointId,
		InstanceId:    instance.Id,
		InDate:        time.Now(),
	})

	if cerr := (repo_errors.ErrorDuplicate{}); errors.As(err, &cerr) {
		err = cmnerrors.Internal(cmnerrors.AlreadyExists(cerr.What...))
	} else if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return storage, err
}

func (self *InstanceStateMachine) actionStopStorage(
	storage records.Storage,
) error {
	storage.OutDate = new(time.Time)
	*storage.OutDate = time.Now()

	repo := self.repos.storage.GetStorageRepository()
	err := repo.Update(storage)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.Internal(cmnerrors.NotFound(cerr.What...))
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return err
}

