package access

import (
	"errors"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/misc/authorizer"

	"github.com/google/uuid"
)

type Renter struct {
	authorizer *authorizer.Authorizer
}

func NewRenter(authorizer *authorizer.Authorizer) *Renter {
	return &Renter{authorizer}
}

func (self *Renter) Access(userId uuid.UUID, renterUserId uuid.UUID) error {
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

