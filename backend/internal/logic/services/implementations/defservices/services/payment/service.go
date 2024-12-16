package payment

import (
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/access"
	"rent_service/internal/logic/services/implementations/defservices/authenticator"
	"rent_service/internal/logic/services/implementations/defservices/authorizer"
	"rent_service/internal/logic/services/implementations/defservices/emptymathcer"
	"rent_service/internal/logic/services/implementations/defservices/paymentcheckers"
	"rent_service/internal/logic/services/interfaces/payment"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/logic/services/types/token"
	. "rent_service/internal/misc/types/collection"
	payment_provider "rent_service/internal/repository/context/providers/payment"
	paymethod_provider "rent_service/internal/repository/context/providers/paymethod"
	user_provider "rent_service/internal/repository/context/providers/user"
	"rent_service/misc/mapfuncs"

	"github.com/google/uuid"
)

type payMethodRepoProviders struct {
	payMethod paymethod_provider.IProvider
}

type payMethodService struct {
	repos payMethodRepoProviders
}

func NewPayMethod(payMethod paymethod_provider.IProvider) payment.IPayMethodService {
	return &payMethodService{payMethodRepoProviders{payMethod}}
}

func mapPayMethod(value *models.PayMethod) payment.PayMethod {
	return payment.PayMethod{
		Id:          value.Id,
		Name:        value.Name,
		Description: value.Description,
	}
}

func (self *payMethodService) GetPayMethods() (Collection[payment.PayMethod], error) {
	var methods Collection[models.PayMethod]

	err := cmnerrors.RepoCallWrap(func() (err error) {
		repo := self.repos.payMethod.GetPayMethodRepository()
		methods, err = repo.GetAll()
		return
	})

	return MapCollection(mapPayMethod, methods), err
}

type userPayMethodRepoProviders struct {
	payMethod user_provider.IPayMethodsProvider
}

type userPayMethodService struct {
	repos         userPayMethodRepoProviders
	checkers      map[uuid.UUID]paymentcheckers.IRegistrationChecker
	authenticator authenticator.IAuthenticator
	authorizer    authorizer.IAuthorizer
}

func NewUserPayMethod(
	authenticator authenticator.IAuthenticator,
	authorizer authorizer.IAuthorizer,
	checkers map[uuid.UUID]paymentcheckers.IRegistrationChecker,
	payMethod user_provider.IPayMethodsProvider,
) payment.IUserPayMethodService {
	return &userPayMethodService{
		userPayMethodRepoProviders{payMethod},
		checkers,
		authenticator,
		authorizer,
	}
}

func mapUserPayMethods(
	value *models.UserPayMethods,
) (Collection[payment.UserPayMethod], error) {
	out := make([]payment.UserPayMethod, len(value.Map))

	for id, method := range value.Map {
		if uint(len(out)) <= method.Priority {
			return nil, cmnerrors.Internal(cmnerrors.Incorrect("priority"))
		}

		out[method.Priority] = payment.UserPayMethod{
			Id:       id,
			MethodId: method.MethodId,
			Name:     method.Name,
		}
	}

	return SliceCollection(out), nil
}

func (self *userPayMethodService) GetPayMethods(
	token token.Token,
) (Collection[payment.UserPayMethod], error) {
	var paymethods Collection[payment.UserPayMethod]
	var mpaymethods models.UserPayMethods
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = cmnerrors.RepoCallWrap(func() (err error) {
			repo := self.repos.payMethod.GetUserPayMethodsRepository()
			mpaymethods, err = repo.GetByUserId(user.Id)
			return
		})
	}

	if nil == err {
		paymethods, err = mapUserPayMethods(&mpaymethods)
	}

	return paymethods, err
}

func (self *userPayMethodService) RegisterPayMethod(
	token token.Token,
	method payment.PayMethodRegistrationForm,
) (uuid.UUID, error) {
	var id uuid.UUID
	var paymethods models.UserPayMethods
	var user models.User

	err := emptymathcer.Match(
		emptymathcer.Item("payer_id", method.PayerId),
		emptymathcer.Item("name", method.Name),
	)

	if nil == err {
		user, err = self.authenticator.LoginWithToken(token)
	}

	if nil == err {
		if checker, ok := self.checkers[method.MethodId]; ok {
			err = checker.CheckPayerId(method.PayerId)
		} else {
			err = cmnerrors.Unknown("method_id")
		}
	}

	if nil == err {
		err = cmnerrors.RepoCallWrap(func() (err error) {
			repo := self.repos.payMethod.GetUserPayMethodsRepository()
			paymethods, err = repo.CreatePayMethod(user.Id, models.UserPayMethod{
				Name:     method.Name,
				MethodId: method.MethodId,
				PayerId:  method.PayerId,
			})
			return
		})
	}

	if nil == err {
		if fid, ok := mapfuncs.FindByValueF(
			paymethods.Map,
			func(m *models.UserPayMethod) bool {
				return m.MethodId == method.MethodId
			},
		); ok {
			id = fid
		} else {
			err = cmnerrors.Internal(errors.New("Method wasn't added..."))
		}
	}

	return id, err
}

func (self *userPayMethodService) UpdatePayMethodsPriority(
	token token.Token,
	methodsOrder []uuid.UUID,
) error {
	var paymethods models.UserPayMethods
	repo := self.repos.payMethod.GetUserPayMethodsRepository()
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = cmnerrors.RepoCallWrap(func() (err error) {
			paymethods, err = repo.GetByUserId(user.Id)
			return
		})
	}

	if nil == err && len(paymethods.Map) != len(methodsOrder) {
		err = payment.ErrorIncompletePayMethodsList
	}

	if nil == err {
		visited := make(map[uuid.UUID]bool)

		for i := uint(0); nil == err && uint(len(methodsOrder)) > i; i++ {
			id := methodsOrder[i]

			if method, ok := paymethods.Map[id]; !ok {
				err = cmnerrors.Unknown("method_id")
			} else if visited[id] {
				err = payment.ErrorIncompletePayMethodsList
			} else {
				method.Priority = i
				paymethods.Map[id] = method
				visited[id] = true
			}
		}
	}

	if nil == err {
		err = cmnerrors.RepoCallWrap(func() error {
			return repo.Update(paymethods)
		})
	}

	return err
}

func (self *userPayMethodService) RemovePayMethod(
	token token.Token,
	methodId uuid.UUID,
) error {
	var paymethods models.UserPayMethods
	repo := self.repos.payMethod.GetUserPayMethodsRepository()
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = cmnerrors.RepoCallWrap(func() (err error) {
			paymethods, err = repo.GetByUserId(user.Id)
			return
		})
	}

	if nil == err {
		if _, ok := paymethods.Map[methodId]; ok {
			delete(paymethods.Map, methodId)
		} else {
			err = cmnerrors.NotFound("id")
		}
	}

	if nil == err {
		err = cmnerrors.RepoCallWrap(func() error {
			return repo.Update(paymethods)
		})
	}

	return err
}

type rentPaymentRepoProviders struct {
	payment payment_provider.IProvider
}

type rentPaymentAccessors struct {
	instance access.IInstance
	rent     access.IRent
}

type rentPaymentService struct {
	repos         rentPaymentRepoProviders
	access        rentPaymentAccessors
	authenticator authenticator.IAuthenticator
	authorizer    authorizer.IAuthorizer
}

func NewRentPayment(
	authenticator authenticator.IAuthenticator,
	authorizer authorizer.IAuthorizer,
	instance access.IInstance,
	rent access.IRent,
	payment payment_provider.IProvider,
) payment.IRentPaymentService {
	return &rentPaymentService{
		rentPaymentRepoProviders{payment},
		rentPaymentAccessors{instance, rent},
		authenticator,
		authorizer,
	}
}

func mapPayment(value *models.Payment) payment.Payment {
	var out = payment.Payment{
		Id:          value.Id,
		RentId:      value.RentId,
		PeriodStart: date.New(value.PeriodStart),
		PeriodEnd:   date.New(value.PeriodEnd),
		Value:       value.Value,
		Status:      value.Status,
		CreateDate:  date.New(value.CreateDate),
	}

	if nil != value.PayMethodId {
		out.PayMethodId = new(uuid.UUID)
		*out.PayMethodId = *value.PayMethodId
	}

	if nil != value.PaymentDate {
		out.PaymentDate = new(date.Date)
		*out.PaymentDate = date.New(*value.PaymentDate)
	}

	return out
}

func (self *rentPaymentService) GetPaymentsByInstance(
	token token.Token,
	instanceId uuid.UUID,
) (Collection[payment.Payment], error) {
	var payments Collection[models.Payment]
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.instance.Access(user.Id, instanceId)
	}

	if nil == err {
		err = cmnerrors.RepoCallWrap(func() (err error) {
			repo := self.repos.payment.GetPaymentRepository()
			payments, err = repo.GetByInstanceId(instanceId)
			return
		})
	}

	return MapCollection(mapPayment, payments), err
}

func (self *rentPaymentService) GetPaymentsByRentId(
	token token.Token,
	rentId uuid.UUID,
) (Collection[payment.Payment], error) {
	var payments Collection[models.Payment]
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.rent.Access(user.Id, rentId)
	}

	if nil == err {
		err = cmnerrors.RepoCallWrap(func() (err error) {
			repo := self.repos.payment.GetPaymentRepository()
			payments, err = repo.GetByRentId(rentId)
			return
		})
	}

	return MapCollection(mapPayment, payments), err
}

