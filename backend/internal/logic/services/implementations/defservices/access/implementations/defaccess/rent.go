package defaccess

import (
	"errors"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/access"
	"rent_service/internal/logic/services/implementations/defservices/authorizer"
	"rent_service/internal/repository/context/providers/rent"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"

	"github.com/google/uuid"
)

type raccess struct {
	rentRepo   rent.IProvider
	authorizer authorizer.IAuthorizer
}

func NewRent(rent rent.IProvider, authorizer authorizer.IAuthorizer) access.IRent {
	return &raccess{rent, authorizer}
}

func (self *raccess) Access(userId uuid.UUID, rentId uuid.UUID) error {
	repo := self.rentRepo.GetRentRepository()
	rent, err := repo.GetById(rentId)

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

	if rent.UserId == userId {
		return nil
	}

	return cmnerrors.Authorization(cmnerrors.NoAccess("rent"))
}

type rreqaccess struct {
	requestRepo rent.IRequestProvider
	authorizer  authorizer.IAuthorizer
}

func NewRentRequest(
	requestRepo rent.IRequestProvider,
	authorizer authorizer.IAuthorizer,
) access.IRentRequest {
	return &rreqaccess{requestRepo, authorizer}
}

func (self *rreqaccess) Access(userId uuid.UUID, requestId uuid.UUID) error {
	repo := self.requestRepo.GetRentRequestRepository()
	request, err := repo.GetById(requestId)

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
		if sk.PickUpPointId == request.PickUpPointId {
			return nil
		}
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(skerr, &cerr) {
		return skerr
	}

	return cmnerrors.Authorization(cmnerrors.NoAccess("rent_request"))
}

type rretaccess struct {
	returnRepo rent.IReturnProvider
	authorizer authorizer.IAuthorizer
}

func NewRentReturn(
	returnRepo rent.IReturnProvider,
	authorizer authorizer.IAuthorizer,
) access.IRentReturn {
	return &rretaccess{returnRepo, authorizer}
}

func (self *rretaccess) Access(userId uuid.UUID, requestId uuid.UUID) error {
	repo := self.returnRepo.GetRentReturnRepository()
	request, err := repo.GetById(requestId)

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
		if sk.PickUpPointId == request.PickUpPointId {
			return nil
		}
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(skerr, &cerr) {
		return skerr
	}

	return cmnerrors.Authorization(cmnerrors.NoAccess("rent_return"))
}

