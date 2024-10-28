package access

import (
	"errors"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/misc/authorizer"
	"rent_service/internal/repository/context/providers/pickuppoint"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"

	"github.com/google/uuid"
)

type PickUpPoint struct {
	pickUpPointRepo pickuppoint.IProvider
	authorizer      *authorizer.Authorizer
}

func NewPickUpPoint(
	pickUpPoint pickuppoint.IProvider,
	authorizer *authorizer.Authorizer,
) *PickUpPoint {
	return &PickUpPoint{pickUpPoint, authorizer}
}

func (self *PickUpPoint) Access(
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

