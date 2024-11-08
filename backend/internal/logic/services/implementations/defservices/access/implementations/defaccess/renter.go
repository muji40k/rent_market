package defaccess

import (
	"errors"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/access"
	"rent_service/internal/logic/services/implementations/defservices/authorizer"

	"github.com/google/uuid"
)

type renaccess struct {
	authorizer authorizer.IAuthorizer
}

func NewRenter(authorizer authorizer.IAuthorizer) access.IRenter {
	return &renaccess{authorizer}
}

func (self *renaccess) Access(userId uuid.UUID, renterUserId uuid.UUID) error {
	if _, rerr := self.authorizer.IsRenter(renterUserId); nil != rerr {
		return rerr
	}

	if _, aerr := self.authorizer.IsAdministrator(userId); nil == aerr {
		return nil
	} else if cerr := (cmnerrors.ErrorInternal{}); errors.As(aerr, &cerr) {
		return aerr
	}

	if userId == renterUserId {
		return nil
	}

	return cmnerrors.Authorization(cmnerrors.NoAccess("renter"))
}

