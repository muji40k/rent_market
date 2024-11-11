package access

import (
	"rent_service/internal/logic/services/implementations/defservices/access"
	"rent_service/internal/logic/services/implementations/defservices/access/implementations/defaccess"
	"rent_service/internal/logic/services/implementations/defservices/authorizer"
	"rent_service/internal/repository/implementation/mock/instance"
	"rent_service/internal/repository/implementation/mock/provision"
	"rent_service/internal/repository/implementation/mock/storage"

	instance_pmock "rent_service/internal/repository/context/mock/instance"
	provision_pmock "rent_service/internal/repository/context/mock/provision"
	storage_pmock "rent_service/internal/repository/context/mock/storage"

	"go.uber.org/mock/gomock"
)

type InstanceBuilder struct {
	authorizer authorizer.IAuthorizer
	provision  *mock_provision.MockIRepository
	instance   *mock_instance.MockIRepository
	storage    *mock_storage.MockIRepository
}

func NewInstance(ctrl *gomock.Controller) *InstanceBuilder {
	return &InstanceBuilder{
		provision: mock_provision.NewMockIRepository(ctrl),
		instance:  mock_instance.NewMockIRepository(ctrl),
		storage:   mock_storage.NewMockIRepository(ctrl),
	}
}

func (self *InstanceBuilder) WithAuthorizer(authorizer authorizer.IAuthorizer) *InstanceBuilder {
	self.authorizer = authorizer
	return self
}

func (self *InstanceBuilder) WithProvisionRepository(f func(*mock_provision.MockIRepository)) *InstanceBuilder {
	f(self.provision)
	return self
}

func (self *InstanceBuilder) WithInstanceRepository(f func(*mock_instance.MockIRepository)) *InstanceBuilder {
	f(self.instance)
	return self
}

func (self *InstanceBuilder) WithStorageRepository(f func(*mock_storage.MockIRepository)) *InstanceBuilder {
	f(self.storage)
	return self
}

func (self *InstanceBuilder) Build() access.IInstance {
	return defaccess.NewInstance(
		self.authorizer,
		provision_pmock.New(self.provision),
		instance_pmock.New(self.instance),
		storage_pmock.New(self.storage),
	)
}

