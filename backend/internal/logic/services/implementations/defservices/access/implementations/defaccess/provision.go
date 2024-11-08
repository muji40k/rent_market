package defaccess

import (
	"errors"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/access"
	"rent_service/internal/logic/services/implementations/defservices/authorizer"
	"rent_service/internal/repository/context/providers/provision"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"

	"github.com/google/uuid"
)

type paccess struct {
	provisionRepo provision.IProvider
	authorizer    authorizer.IAuthorizer
}

func NewProvision(provision provision.IProvider, authorizer authorizer.IAuthorizer) access.IProvision {
	return &paccess{provision, authorizer}
}

func (self *paccess) Access(userId uuid.UUID, provisionId uuid.UUID) error {
	repo := self.provisionRepo.GetProvisionRepository()
	provision, err := repo.GetById(provisionId)

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

	if renter, rerr := self.authorizer.IsRenter(userId); nil == rerr {
		if provision.RenterId == renter.Id {
			return nil
		}
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(rerr, &cerr) {
		return rerr
	}

	return cmnerrors.Authorization(cmnerrors.NoAccess("provision"))
}

type praccess struct {
	requestRepo provision.IRequestProvider
	authorizer  authorizer.IAuthorizer
}

func NewProvisionRequest(
	requestRepo provision.IRequestProvider,
	authorizer authorizer.IAuthorizer,
) access.IProvisionRequest {
	return &praccess{requestRepo, authorizer}
}

func (self *praccess) Access(userId uuid.UUID, requestId uuid.UUID) error {
	repo := self.requestRepo.GetProvisionRequestRepository()
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

	return cmnerrors.Authorization(cmnerrors.NoAccess("provision_request"))
}

type preaccess struct {
	revokeRepo provision.IRevokeProvider
	authorizer authorizer.IAuthorizer
}

func NewProvisionRevoke(
	revokeRepo provision.IRevokeProvider,
	authorizer authorizer.IAuthorizer,
) access.IProvisionRevoke {
	return &preaccess{revokeRepo, authorizer}
}

func (self *preaccess) Access(userId uuid.UUID, revokeId uuid.UUID) error {
	repo := self.revokeRepo.GetRevokeProvisionRepository()
	revoke, err := repo.GetById(revokeId)

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
		if sk.PickUpPointId == revoke.PickUpPointId {
			return nil
		}
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(skerr, &cerr) {
		return skerr
	}

	return cmnerrors.Authorization(cmnerrors.NoAccess("provision_revoke"))
}

