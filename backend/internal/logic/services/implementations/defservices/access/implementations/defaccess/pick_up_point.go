package defaccess

import (
	"errors"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/access"
	"rent_service/internal/logic/services/implementations/defservices/authorizer"
	"rent_service/internal/repository/context/providers/pickuppoint"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"

	"github.com/google/uuid"
)

type pupaccess struct {
	pickUpPointRepo pickuppoint.IProvider
	authorizer      authorizer.IAuthorizer
}

func NewPickUpPoint(
	pickUpPoint pickuppoint.IProvider,
	authorizer authorizer.IAuthorizer,
) access.IPickUpPoint {
	return &pupaccess{pickUpPoint, authorizer}
}

func (self *pupaccess) Access(
	userId uuid.UUID,
	pickUpPointId uuid.UUID,
) error {
	repo := self.pickUpPointRepo.GetPickUpPointRepository()
	_, err := repo.GetById(pickUpPointId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		return cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		return cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	if _, aerr := self.authorizer.IsAdministrator(userId); nil == aerr {
		return nil
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(aerr, &cerr) {
		return aerr
	}

	if sk, skerr := self.authorizer.IsStorekeeper(userId); nil == skerr {
		if sk.PickUpPointId == pickUpPointId {
			return nil
		}
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(skerr, &cerr) {
		return skerr
	}

	return cmnerrors.Authorization(cmnerrors.NoAccess("pick_up_point"))
}

