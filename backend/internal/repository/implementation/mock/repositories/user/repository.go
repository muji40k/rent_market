package user

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/repository/implementation/mock/cmnerrors"

	"github.com/google/uuid"
)

type MockRepository struct {
	create     func(user models.User) (models.User, error)
	update     func(user models.User) error
	getById    func(userId uuid.UUID) (models.User, error)
	getByEmail func(email string) (models.User, error)
	getByToken func(token models.Token) (models.User, error)
}

func New() *MockRepository {
	return &MockRepository{
		func(user models.User) (models.User, error) {
			return models.User{}, cmnerrors.ErrorNotSet
		},
		func(user models.User) error {
			return cmnerrors.ErrorNotSet
		},
		func(userId uuid.UUID) (models.User, error) {
			return models.User{}, cmnerrors.ErrorNotSet
		},
		func(email string) (models.User, error) {
			return models.User{}, cmnerrors.ErrorNotSet
		},
		func(token models.Token) (models.User, error) {
			return models.User{}, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockRepository) Create(user models.User) (models.User, error) {
	return self.create(user)
}

func (self *MockRepository) WithCreate(f func(user models.User) (models.User, error)) *MockRepository {
	self.create = f
	return self
}

func (self *MockRepository) Update(user models.User) error {
	return self.update(user)
}

func (self *MockRepository) WithUpdate(f func(user models.User) error) *MockRepository {
	self.update = f
	return self
}

func (self *MockRepository) GetById(userId uuid.UUID) (models.User, error) {
	return self.getById(userId)
}

func (self *MockRepository) WithGetById(f func(userId uuid.UUID) (models.User, error)) *MockRepository {
	self.getById = f
	return self
}

func (self *MockRepository) GetByEmail(email string) (models.User, error) {
	return self.getByEmail(email)
}

func (self *MockRepository) WithGetByEmail(f func(email string) (models.User, error)) *MockRepository {
	self.getByEmail = f
	return self
}

func (self *MockRepository) GetByToken(token models.Token) (models.User, error) {
	return self.getByToken(token)
}

func (self *MockRepository) WithGetByToken(f func(token models.Token) (models.User, error)) *MockRepository {
	self.getByToken = f
	return self
}

type MockProfileRepository struct {
	create      func(profile models.UserProfile) (models.UserProfile, error)
	update      func(profile models.UserProfile) error
	getByUserId func(userId uuid.UUID) (models.UserProfile, error)
}

func NewProfile() *MockProfileRepository {
	return &MockProfileRepository{
		func(profile models.UserProfile) (models.UserProfile, error) {
			return models.UserProfile{}, cmnerrors.ErrorNotSet
		},
		func(profile models.UserProfile) error {
			return cmnerrors.ErrorNotSet
		},
		func(userId uuid.UUID) (models.UserProfile, error) {
			return models.UserProfile{}, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockProfileRepository) Create(profile models.UserProfile) (models.UserProfile, error) {
	return self.create(profile)
}

func (self *MockProfileRepository) WithCreate(f func(profile models.UserProfile) (models.UserProfile, error)) *MockProfileRepository {
	self.create = f
	return self
}

func (self *MockProfileRepository) Update(profile models.UserProfile) error {
	return self.update(profile)
}

func (self *MockProfileRepository) WithUpdate(f func(profile models.UserProfile) error) *MockProfileRepository {
	self.update = f
	return self
}

func (self *MockProfileRepository) GetByUserId(userId uuid.UUID) (models.UserProfile, error) {
	return self.getByUserId(userId)
}

func (self *MockProfileRepository) WithGetByUserId(f func(userId uuid.UUID) (models.UserProfile, error)) *MockProfileRepository {
	self.getByUserId = f
	return self
}

type MockFavoriteRepository struct {
	create      func(profile models.UserFavoritePickUpPoint) (models.UserFavoritePickUpPoint, error)
	update      func(profile models.UserFavoritePickUpPoint) error
	getByUserId func(userId uuid.UUID) (models.UserFavoritePickUpPoint, error)
}

func NewFavorite() *MockFavoriteRepository {
	return &MockFavoriteRepository{
		func(profile models.UserFavoritePickUpPoint) (models.UserFavoritePickUpPoint, error) {
			return models.UserFavoritePickUpPoint{}, cmnerrors.ErrorNotSet
		},
		func(profile models.UserFavoritePickUpPoint) error {
			return cmnerrors.ErrorNotSet
		},
		func(userId uuid.UUID) (models.UserFavoritePickUpPoint, error) {
			return models.UserFavoritePickUpPoint{}, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockFavoriteRepository) Create(profile models.UserFavoritePickUpPoint) (models.UserFavoritePickUpPoint, error) {
	return self.create(profile)
}

func (self *MockFavoriteRepository) WithCreate(f func(profile models.UserFavoritePickUpPoint) (models.UserFavoritePickUpPoint, error)) *MockFavoriteRepository {
	self.create = f
	return self
}

func (self *MockFavoriteRepository) Update(profile models.UserFavoritePickUpPoint) error {
	return self.update(profile)
}

func (self *MockFavoriteRepository) WithUpdate(f func(profile models.UserFavoritePickUpPoint) error) *MockFavoriteRepository {
	self.update = f
	return self
}

func (self *MockFavoriteRepository) GetByUserId(userId uuid.UUID) (models.UserFavoritePickUpPoint, error) {
	return self.getByUserId(userId)
}

func (self *MockFavoriteRepository) WithGetByUserId(f func(userId uuid.UUID) (models.UserFavoritePickUpPoint, error)) *MockFavoriteRepository {
	self.getByUserId = f
	return self
}

type MockPayMethodsRepository struct {
	createPayMethod func(userId uuid.UUID, payMethod models.UserPayMethod) (models.UserPayMethods, error)
	update          func(payMethods models.UserPayMethods) error
	getByUserId     func(userId uuid.UUID) (models.UserPayMethods, error)
}

func NewPayMethods() *MockPayMethodsRepository {
	return &MockPayMethodsRepository{
		func(userId uuid.UUID, payMethod models.UserPayMethod) (models.UserPayMethods, error) {
			return models.UserPayMethods{}, cmnerrors.ErrorNotSet
		},
		func(payMethods models.UserPayMethods) error {
			return cmnerrors.ErrorNotSet
		},
		func(userId uuid.UUID) (models.UserPayMethods, error) {
			return models.UserPayMethods{}, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockPayMethodsRepository) CreatePayMethod(userId uuid.UUID, payMethod models.UserPayMethod) (models.UserPayMethods, error) {
	return self.createPayMethod(userId, payMethod)
}

func (self *MockPayMethodsRepository) WithCreatePayMethod(f func(userId uuid.UUID, payMethod models.UserPayMethod) (models.UserPayMethods, error)) *MockPayMethodsRepository {
	self.createPayMethod = f
	return self
}

func (self *MockPayMethodsRepository) Update(payMethods models.UserPayMethods) error {
	return self.update(payMethods)
}

func (self *MockPayMethodsRepository) WithUpdate(f func(payMethods models.UserPayMethods) error) *MockPayMethodsRepository {
	self.update = f
	return self
}

func (self *MockPayMethodsRepository) GetByUserId(userId uuid.UUID) (models.UserPayMethods, error) {
	return self.getByUserId(userId)
}

func (self *MockPayMethodsRepository) WithGetByUserId(f func(userId uuid.UUID) (models.UserPayMethods, error)) *MockPayMethodsRepository {
	self.getByUserId = f
	return self
}

