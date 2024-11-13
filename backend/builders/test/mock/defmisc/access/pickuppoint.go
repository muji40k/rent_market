package access

import (
	"rent_service/internal/logic/services/implementations/defservices/access"
	"rent_service/internal/logic/services/implementations/defservices/access/implementations/defaccess"
	"rent_service/internal/logic/services/implementations/defservices/authorizer"
	"rent_service/internal/repository/implementation/mock/pickuppoint"

	pickuppoint_pmock "rent_service/internal/repository/context/mock/pickuppoint"

	"go.uber.org/mock/gomock"
)

type PickUpPointBuilder struct {
	authorizer  authorizer.IAuthorizer
	pickUpPoint *mock_pickuppoint.MockIRepository
}

func NewPickUpPoint(ctrl *gomock.Controller) *PickUpPointBuilder {
	return &PickUpPointBuilder{
		pickUpPoint: mock_pickuppoint.NewMockIRepository(ctrl),
	}
}

func (self *PickUpPointBuilder) WithAuthorizer(authorizer authorizer.IAuthorizer) *PickUpPointBuilder {
	self.authorizer = authorizer
	return self
}

func (self *PickUpPointBuilder) WithPickUpPointRepository(f func(*mock_pickuppoint.MockIRepository)) *PickUpPointBuilder {
	f(self.pickUpPoint)
	return self
}

func (self *PickUpPointBuilder) Build() access.IPickUpPoint {
	return defaccess.NewPickUpPoint(
		pickuppoint_pmock.New(self.pickUpPoint),
		self.authorizer,
	)
}

