package user

import (
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/authenticator"
	"rent_service/internal/logic/services/implementations/defservices/authorizer"
	"rent_service/internal/logic/services/implementations/defservices/codegen"
	"rent_service/internal/logic/services/implementations/defservices/emptymathcer"
	"rent_service/internal/logic/services/implementations/defservices/photoregistry"
	"rent_service/internal/logic/services/interfaces/user"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/logic/services/types/token"
	"rent_service/internal/repository/context/providers/role"
	user_provider "rent_service/internal/repository/context/providers/user"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
	"time"

	"github.com/google/uuid"
)

func gmap[T any](from *T, to **T) {
	if nil != from {
		if nil == *to {
			*to = new(T)
		}

		**to = *from
	}
}

func gmapf[T any, F any](from *T, to **F, mapf func(*T) F) {
	if nil != from {
		if nil == *to {
			*to = new(F)
		}

		**to = mapf(from)
	}
}

func gunmap[T any](from *T, to **T) {
	if nil != from {
		if nil == *to {
			*to = new(T)
		}

		**to = *from
	} else {
		*to = nil
	}
}

func gunmapf[T any, F any](from *T, to **F, mapf func(*T) F) {
	if nil != from {
		if nil == *to {
			*to = new(F)
		}

		**to = mapf(from)
	} else {
		*to = nil
	}
}

type repoproviders struct {
	user user_provider.IProvider
}

type service struct {
	repos         repoproviders
	authenticator authenticator.IAuthenticator
}

func New(
	user user_provider.IProvider,
	authenticator authenticator.IAuthenticator,
) user.IService {
	return &service{repoproviders{user}, authenticator}
}

func (self *service) GetSelfUserInfo(token token.Token) (user.Info, error) {
	var info user.Info
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		info.Id = user.Id
		info.Email = user.Email
		info.Name = user.Name
	}

	return info, err
}

func (self *service) UpdateSelfUserInfo(
	token token.Token,
	form user.UpdateForm,
) error {
	var user models.User
	err := emptymathcer.Match(
		emptymathcer.Item("email", form.Email),
		emptymathcer.Item("name", form.Name),
	)

	if nil == err {
		user, err = self.authenticator.LoginWithToken(token)
	}

	if nil == err {
		user.Email = form.Email
		user.Name = form.Name
		repo := self.repos.user.GetUserRepository()
		err = repo.Update(user)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return err
}

func (self *service) UpdateSelfUserPassword(
	token token.Token,
	old_password string,
	new_password string,
) error {
	var user models.User
	err := emptymathcer.Match(
		emptymathcer.Item("old_password", old_password),
		emptymathcer.Item("new_password", new_password),
	)

	if nil == err {
		user, err = self.authenticator.LoginWithToken(token)
	}

	if nil == err && old_password != user.Password {
		err = cmnerrors.Authentication(errors.New("Passwords don't match"))
	}

	if nil == err {
		user.Password = new_password
		repo := self.repos.user.GetUserRepository()
		err = repo.Update(user)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return err
}

type passwordUpdateProviders struct {
	user     user_provider.IProvider
	password user_provider.IPasswordUpdateProvider
}

type passwordUpdateService struct {
	repos         passwordUpdateProviders
	authenticator authenticator.IAuthenticator
	verificator   I2FA
	generator     codegen.IGenerator
	encoder       ICodeEncoder
	validTime     time.Duration
}

func mapPassowrdUpdateRequest(value *models.UserPasswordUpdateRequest) user.PasswordUpdateRequest {
	return user.PasswordUpdateRequest{
		Required: true,
		Id:       value.Id,
		ValidTo:  value.ValidTo,
	}
}

type ICodeEncoder interface {
	Encode(code string) string
}

type I2FA interface {
	SendCode(user models.User, code string) error
}

func NewPasswordUpdate(
	authenticator authenticator.IAuthenticator,
	verificator I2FA,
	generator codegen.IGenerator,
	encoder ICodeEncoder,
	validTime time.Duration,
	userRepo user_provider.IProvider,
	passwordRepo user_provider.IPasswordUpdateProvider,
) user.IPasswordUpdateService {
	return &passwordUpdateService{
		passwordUpdateProviders{userRepo, passwordRepo},
		authenticator,
		verificator,
		generator,
		encoder,
		validTime,
	}
}

func (self *passwordUpdateService) updatePassword(
	user models.User,
	password string,
) error {
	user.Password = password
	repo := self.repos.user.GetUserRepository()
	err := repo.Update(user)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return err
}

func (self *passwordUpdateService) RequestPasswordUpdate(
	token token.Token,
	old_password string,
	new_password string,
) (user.PasswordUpdateRequest, error) {
	result := user.PasswordUpdateRequest{Required: false}
	var user models.User
	err := emptymathcer.Match(
		emptymathcer.Item("old_password", old_password),
		emptymathcer.Item("new_password", new_password),
	)

	if nil == err {
		user, err = self.authenticator.LoginWithToken(token)
	}

	if nil == err && old_password != user.Password {
		err = cmnerrors.Authentication(errors.New("Passwords don't match"))
	}

	if nil == err && nil == self.verificator {
		err = self.updatePassword(user, new_password)
	} else if nil == err {
		var request models.UserPasswordUpdateRequest
		code := self.generator.Generate()
		err = self.verificator.SendCode(user, code)

		if nil != err {
			err = cmnerrors.Internal(err)
		} else {
			request = models.UserPasswordUpdateRequest{
				UserId:      user.Id,
				NewPassword: new_password,
				Code:        self.encoder.Encode(code),
				ValidTo:     time.Now().Add(self.validTime),
			}
			repo := self.repos.password.GetUserPasswordUpdateRepository()
			request, err = repo.Create(request)

			if cerr := (repo_errors.ErrorDuplicate{}); errors.As(err, &cerr) {
				err = cmnerrors.AlreadyExists("password_update_request")
			} else if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
				err = cmnerrors.NotFound(cerr.What...)
			} else if nil != err {
				err = cmnerrors.Internal(cmnerrors.DataAccess(err))
			}
		}

		if nil == err {
			result = mapPassowrdUpdateRequest(&request)
		}
	}

	return result, err
}

func (self *passwordUpdateService) AuthenticatePasswordUpdateRequest(
	token token.Token,
	requestId uuid.UUID,
	code string,
) error {
	var user models.User
	var request models.UserPasswordUpdateRequest
	err := emptymathcer.Match(
		emptymathcer.Item("code", code),
	)

	if nil == err {
		user, err = self.authenticator.LoginWithToken(token)
	}

	if nil == err {
		repo := self.repos.password.GetUserPasswordUpdateRepository()
		request, err = repo.GetById(requestId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	if after := time.Now().After(request.ValidTo); nil == err &&
		(request.UserId != user.Id || request.Code != code || after) {
		err = cmnerrors.Authorization(
			errors.New("Can't verify request"),
		)

		if after {
			_ = self.repos.password.
				GetUserPasswordUpdateRepository().
				Remove(requestId)
		}
	}

	if nil == err {
		err = self.updatePassword(user, request.NewPassword)
	}

	if nil == err {
		repo := self.repos.password.GetUserPasswordUpdateRepository()
		err = repo.Remove(requestId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return err
}

type repoProfileProviders struct {
	profile user_provider.IProfileProvider
}

type profileService struct {
	repos         repoProfileProviders
	authenticator authenticator.IAuthenticator
	photo         photoregistry.IRegistry
}

func NewProfile(
	profile user_provider.IProfileProvider,
	authenticator authenticator.IAuthenticator,
	photo photoregistry.IRegistry,
) user.IProfileService {
	return &profileService{repoProfileProviders{profile}, authenticator, photo}
}

func mapProfile(value *models.UserProfile) user.UserProfile {
	var user user.UserProfile

	gmap(value.Name, &user.Name)
	gmap(value.Surname, &user.Surname)
	gmap(value.Patronymic, &user.Patronymic)
	gmapf(value.BirthDate, &user.BirthDate, func(v *time.Time) date.Date {
		return date.New(*v)
	})
	gmap(value.PhotoId, &user.PhotoId)

	return user
}

func unmapProfile(value *user.UserProfile, profile *models.UserProfile) {
	gunmap(value.Name, &profile.Name)
	gunmap(value.Surname, &profile.Surname)
	gunmap(value.Patronymic, &profile.Patronymic)
	gunmapf(value.BirthDate, &profile.BirthDate, func(v *date.Date) time.Time {
		return v.Time
	})
	gunmap(value.PhotoId, &profile.PhotoId)
}

func (self *profileService) GetUserProfile(
	userId uuid.UUID,
) (user.UserProfile, error) {
	var user user.UserProfile
	repo := self.repos.profile.GetUserProfileRepository()
	muser, err := repo.GetByUserId(userId)

	if nil == err {
		user = mapProfile(&muser)
	} else if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return user, err
}

func (self *profileService) GetSelfUserProfile(
	token token.Token,
) (user.UserProfile, error) {
	var profile user.UserProfile
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		profile, err = self.GetUserProfile(user.Id)
	}

	return profile, err
}

func (self *profileService) UpdateSelfUserProfile(
	token token.Token,
	data user.UserProfile,
) error {
	var profile models.UserProfile
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err && nil != data.PhotoId {
		*data.PhotoId, err = self.photo.MoveFromTemp(*data.PhotoId)
	}

	if nil == err {
		repo := self.repos.profile.GetUserProfileRepository()
		profile, err = repo.GetByUserId(user.Id)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			unmapProfile(&data, &profile)
			profile.UserId = user.Id

			_, err = repo.Create(profile)
		} else if nil == err {
			unmapProfile(&data, &profile)
			err = repo.Update(profile)
		}

		if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return err
}

type repoFavoriteProviders struct {
	favorite user_provider.IFavoriteProvider
}

type favoriteService struct {
	repos         repoFavoriteProviders
	authenticator authenticator.IAuthenticator
}

func NewFavorite(
	favorite user_provider.IFavoriteProvider,
	authenticator authenticator.IAuthenticator,
) user.IFavoriteService {
	return &favoriteService{repoFavoriteProviders{favorite}, authenticator}
}

func mapFavorite(value *models.UserFavoritePickUpPoint) user.UserFavoritePickUpPoint {
	var favorite user.UserFavoritePickUpPoint

	gmap(value.PickUpPointId, &favorite.PickUpPointId)

	return favorite
}

func unmapFavorite(
	value *user.UserFavoritePickUpPoint,
	favorite *models.UserFavoritePickUpPoint,
) {
	gunmap(value.PickUpPointId, &favorite.PickUpPointId)
}

func (self *favoriteService) GetUserFavorite(
	userId uuid.UUID,
) (user.UserFavoritePickUpPoint, error) {
	var favorite user.UserFavoritePickUpPoint
	repo := self.repos.favorite.GetUserFavoriteRepository()
	mfavorite, err := repo.GetByUserId(userId)

	if nil == err {
		favorite = mapFavorite(&mfavorite)
	} else if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return favorite, err
}

func (self *favoriteService) GetSelfUserFavorite(
	token token.Token,
) (user.UserFavoritePickUpPoint, error) {
	var favorite user.UserFavoritePickUpPoint
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		favorite, err = self.GetUserFavorite(user.Id)
	}

	return favorite, err
}

func (self *favoriteService) UpdateSelfUserFavorite(
	token token.Token,
	data user.UserFavoritePickUpPoint,
) error {
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		var favorite models.UserFavoritePickUpPoint
		repo := self.repos.favorite.GetUserFavoriteRepository()
		favorite, err = repo.GetByUserId(user.Id)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			unmapFavorite(&data, &favorite)
			favorite.UserId = user.Id

			_, err = repo.Create(favorite)
		} else if nil == err {
			unmapFavorite(&data, &favorite)
			err = repo.Update(favorite)
		}

		if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return err
}

type repoRoleProviders struct {
	renter role.IRenterProvider
}

type roleService struct {
	repos         repoRoleProviders
	authenticator authenticator.IAuthenticator
	authorizer    authorizer.IAuthorizer
}

func NewRole(
	authenticator authenticator.IAuthenticator,
	authorizer authorizer.IAuthorizer,
	renter role.IRenterProvider,
) user.IRoleService {
	return &roleService{repoRoleProviders{renter}, authenticator, authorizer}
}

func (self *roleService) IsRenter(token token.Token) error {
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		_, err = self.authorizer.IsRenter(user.Id)
	}

	return err
}

func (self *roleService) RegisterAsRenter(token token.Token) error {
	usr, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		repo := self.repos.renter.GetRenterRepository()
		_, err = repo.Create(usr.Id)

		if cerr := (repo_errors.ErrorDuplicate{}); errors.As(err, &cerr) {
			err = user.AlreadyRenter(usr.Email)
		} else if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return err
}

func (self *roleService) IsAdministrator(token token.Token) error {
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		_, err = self.authorizer.IsAdministrator(user.Id)
	}

	return err
}

func mapSK(sk *models.Storekeeper) user.StoreKeeper {
	return user.StoreKeeper{
		PickUpPointId: sk.PickUpPointId,
	}
}

func (self *roleService) IsStoreKeeper(token token.Token) (user.StoreKeeper, error) {
	var sk models.Storekeeper
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		sk, err = self.authorizer.IsStorekeeper(user.Id)
	}

	return mapSK(&sk), err
}

