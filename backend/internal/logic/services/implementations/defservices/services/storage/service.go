package storage

import (
	"errors"
	"rent_service/internal/domain/records"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/access"
	"rent_service/internal/logic/services/implementations/defservices/authenticator"
	"rent_service/internal/logic/services/interfaces/storage"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/logic/services/types/token"
	. "rent_service/internal/misc/types/collection"
	storage_provider "rent_service/internal/repository/context/providers/storage"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"

	"github.com/google/uuid"
)

type repoproviders struct {
	storage storage_provider.IProvider
}

type accessors struct {
	pickUpPoint access.IPickUpPoint
}

type service struct {
	repos         repoproviders
	access        accessors
	authenticator authenticator.IAuthenticator
}

func New(
	authenticator authenticator.IAuthenticator,
	storage storage_provider.IProvider,
	pickUpPoint access.IPickUpPoint,
) storage.IService {
	return &service{
		repoproviders{storage},
		accessors{pickUpPoint},
		authenticator,
	}
}

func mapf(value *records.Storage) storage.Storage {
	var out = storage.Storage{
		Id:            value.Id,
		PickUpPointId: value.PickUpPointId,
		InstanceId:    value.InstanceId,
		InDate:        date.New(value.InDate),
	}

	if nil != value.OutDate {
		out.OutDate = new(date.Date)
		*out.OutDate = date.New(*value.OutDate)
	}

	return out
}

func (self *service) ListStoragesByPickUpPoint(
	token token.Token,
	pickUpPointId uuid.UUID,
) (Collection[storage.Storage], error) {
	var storage Collection[records.Storage]
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.pickUpPoint.Access(user.Id, pickUpPointId)
	}

	if nil == err {
		repo := self.repos.storage.GetStorageRepository()
		storage, err = repo.GetActiveByPickUpPointId(pickUpPointId)

		if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return MapCollection(mapf, storage), err
}

func (self *service) GetStorageByInstance(instanceId uuid.UUID) (storage.Storage, error) {
	repo := self.repos.storage.GetStorageRepository()
	storage, err := repo.GetActiveByInstanceId(instanceId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return mapf(&storage), err
}

