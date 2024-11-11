package photoregistry

import (
	"rent_service/internal/logic/services/implementations/defservices/photoregistry"
	"rent_service/internal/logic/services/implementations/defservices/photoregistry/implementations/defregistry"
	"rent_service/internal/logic/services/implementations/defservices/photoregistry/implementations/defregistry/storages/mock"
	"rent_service/internal/repository/implementation/mock/photo"

	photo_pmock "rent_service/internal/repository/context/mock/photo"

	"go.uber.org/mock/gomock"
)

type PhotoRegistryBuilder struct {
	photo   *mock_photo.MockIRepository
	temp    *mock_photo.MockITempRepository
	storage *mock_defregistry.MockIStorage
}

func New(ctrl *gomock.Controller) *PhotoRegistryBuilder {
	return &PhotoRegistryBuilder{
		mock_photo.NewMockIRepository(ctrl),
		mock_photo.NewMockITempRepository(ctrl),
		mock_defregistry.NewMockIStorage(ctrl),
	}
}

func (self *PhotoRegistryBuilder) WithPhotoRepository(f func(repo *mock_photo.MockIRepository)) *PhotoRegistryBuilder {
	f(self.photo)
	return self
}

func (self *PhotoRegistryBuilder) WithPhotoTempRepository(f func(repo *mock_photo.MockITempRepository)) *PhotoRegistryBuilder {
	f(self.temp)
	return self
}

func (self *PhotoRegistryBuilder) WithStorage(f func(storage *mock_defregistry.MockIStorage)) *PhotoRegistryBuilder {
	f(self.storage)
	return self
}

func (self *PhotoRegistryBuilder) Build() photoregistry.IRegistry {
	return defregistry.New(
		photo_pmock.New(self.photo),
		photo_pmock.NewTemp(self.temp),
		self.storage,
	)
}

