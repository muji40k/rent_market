package defservices

import (
	"errors"
	v1 "rent_service/internal/factory/services/v1"
	"rent_service/internal/factory/services/v1/deffactory"
	"rent_service/internal/logic/delivery"
	"rent_service/internal/logic/services/implementations/defservices/codegen"
	"rent_service/internal/logic/services/implementations/defservices/photoregistry/implementations/defregistry"
	"rent_service/internal/logic/services/implementations/defservices/services/payment"
	rv1 "rent_service/internal/repository/context/v1"

	"github.com/google/uuid"
)

type Builder struct {
	repositoryContext *rv1.Context
	codegen           codegen.IGenerator
	photoStorage      defregistry.IStorage
	deliveryCreator   delivery.ICreator
	paymentCheckers   map[uuid.UUID]payment.IRegistrationChecker
}

func New() *Builder {
	return &Builder{
		paymentCheckers: make(map[uuid.UUID]payment.IRegistrationChecker),
	}
}

func (self *Builder) WithRepositoryContext(repositoryContext *rv1.Context) *Builder {
	self.repositoryContext = repositoryContext
	return self
}

func (self *Builder) WithCodegen(codegen codegen.IGenerator) *Builder {
	self.codegen = codegen
	return self
}

func (self *Builder) WithPhotoStorage(photoStorage defregistry.IStorage) *Builder {
	self.photoStorage = photoStorage
	return self
}

func (self *Builder) WithDeliveryCreator(deliveryCreator delivery.ICreator) *Builder {
	self.deliveryCreator = deliveryCreator
	return self
}

func (self *Builder) WithPaymentChecker(checker payment.IRegistrationChecker) *Builder {
	self.paymentCheckers[checker.MethodId()] = checker
	return self
}

func (self *Builder) WithPaymentCheckers(checkers map[uuid.UUID]payment.IRegistrationChecker) *Builder {
	if nil != checkers {
		self.paymentCheckers = checkers
	}

	return self
}

func (self *Builder) Build() (v1.IFactory, error) {
	var factory *deffactory.Factory
	var err error

	if nil == self.repositoryContext {
		err = errors.New("DefaultFactoryBuilder: repository context not set")
	}

	if nil == err && nil == self.codegen {
		err = errors.New("DefaultFactoryBuilder: codegen not set")
	}

	if nil == err && nil == self.photoStorage {
		err = errors.New("DefaultFactoryBuilder: photo storage not set")
	}

	if nil == err && nil == self.deliveryCreator {
		err = errors.New("DefaultFactoryBuilder: delivery creator not set")
	}

	if nil == err && 0 == len(self.paymentCheckers) {
		err = errors.New("DefaultFactoryBuilder: payment checkers not set")
	}

	if nil == err {
		factory = deffactory.New(
			self.repositoryContext,
			self.codegen,
			self.photoStorage,
			self.deliveryCreator,
			self.paymentCheckers,
		)
	}

	return factory, err
}

