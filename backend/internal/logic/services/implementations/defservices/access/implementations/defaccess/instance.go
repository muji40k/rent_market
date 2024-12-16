package defaccess

import (
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/records"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/access"
	"rent_service/internal/logic/services/implementations/defservices/authorizer"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/context/providers/instance"
	"rent_service/internal/repository/context/providers/provision"
	"rent_service/internal/repository/context/providers/storage"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"

	"github.com/google/uuid"
)

type instanceProviders struct {
	provision provision.IProvider
	instance  instance.IProvider
	storage   storage.IProvider
}

type iaccess struct {
	repos      instanceProviders
	authorizer authorizer.IAuthorizer
}

func NewInstance(
	authorizer authorizer.IAuthorizer,
	provision provision.IProvider,
	instance instance.IProvider,
	storage storage.IProvider,
) access.IInstance {
	return &iaccess{
		instanceProviders{provision, instance, storage},
		authorizer,
	}
}

func check[T any](
	userId uuid.UUID,
	instanceId uuid.UUID,
	auth func(uuid.UUID) (T, error),
	access func(T, uuid.UUID) error,
) (bool, error) {
	if v, verr := auth(userId); nil == verr {
		if err := access(v, instanceId); nil == err {
			return false, nil
		} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(err, &cerr) {
			return false, err
		}
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(verr, &cerr) {
		return false, verr
	}

	return true, nil
}

func (self *iaccess) Access(userId uuid.UUID, instanceId uuid.UUID) error {
	repo := self.repos.instance.GetInstanceRepository()
	if _, err := repo.GetById(instanceId); nil != err {
		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			return cmnerrors.NotFound(cerr.What...)
		}

		return cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	c, err := check(userId, instanceId,
		self.authorizer.IsAdministrator,
		self.accessAdministrator,
	)

	if c {
		c, err = check(userId, instanceId,
			self.authorizer.IsRenter,
			self.accessRenter,
		)
	}

	if c {
		c, err = check(userId, instanceId,
			self.authorizer.IsStorekeeper,
			self.accessStorekeeper,
		)
	}

	if c {
		err = cmnerrors.Authorization(cmnerrors.NoAccess("instance"))
	}

	return err
}

func (self *iaccess) accessAdministrator(
	_ models.Administrator,
	_ uuid.UUID,
) error {
	return nil // Can access any instance
}

func (self *iaccess) accessRenter(
	renter models.Renter,
	instanceId uuid.UUID,
) error {
	repo := self.repos.provision.GetProvisionRepository()
	provisions, err := repo.GetByInstanceId(instanceId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NoAccess()
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	if nil == err && !collection.Any(provisions.Iter(), func(p *records.Provision) bool {
		return p.RenterId == renter.Id
	}) {
		err = cmnerrors.NoAccess()
	}

	return err
}

func (self *iaccess) accessStorekeeper(
	sk models.Storekeeper,
	instanceId uuid.UUID,
) error {
	repo := self.repos.storage.GetStorageRepository()
	storage, err := repo.GetActiveByInstanceId(instanceId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NoAccess()
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	if nil == err && storage.PickUpPointId != sk.PickUpPointId {
		err = cmnerrors.NoAccess()
	}

	return err
}

