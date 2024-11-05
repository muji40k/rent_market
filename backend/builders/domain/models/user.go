package models

import (
	"rent_service/builders/misc/nullcommon"
	"rent_service/internal/domain/models"
	"rent_service/misc/nullable"
	"time"

	"github.com/google/uuid"
)

type UserBuilder struct {
	id       uuid.UUID
	name     string
	email    string
	password string
	token    string
}

func NewUser() *UserBuilder {
	return &UserBuilder{}
}

func (self *UserBuilder) WithId(id uuid.UUID) *UserBuilder {
	self.id = id
	return self
}

func (self *UserBuilder) WithName(name string) *UserBuilder {
	self.name = name
	return self
}

func (self *UserBuilder) WithEmail(email string) *UserBuilder {
	self.email = email
	return self
}

func (self *UserBuilder) WithPassword(password string) *UserBuilder {
	self.password = password
	return self
}

func (self *UserBuilder) WithToken(token string) *UserBuilder {
	self.token = token
	return self
}

func (self *UserBuilder) Build() models.User {
	return models.User{
		Id:       self.id,
		Name:     self.name,
		Email:    self.email,
		Password: self.password,
		Token:    models.Token(self.token),
	}
}

type UserProfileBuilder struct {
	id         uuid.UUID
	userId     uuid.UUID
	name       *string
	surname    *string
	patronymic *string
	birthDate  *time.Time
	photoId    *uuid.UUID
}

func NewUserProfile() *UserProfileBuilder {
	return &UserProfileBuilder{}
}

func (self *UserProfileBuilder) WithId(id uuid.UUID) *UserProfileBuilder {
	self.id = id
	return self
}

func (self *UserProfileBuilder) WithUserId(userId uuid.UUID) *UserProfileBuilder {
	self.userId = userId
	return self
}

func (self *UserProfileBuilder) WithName(name *nullable.Nullable[string]) *UserProfileBuilder {
	self.name = nullcommon.CopyPtrIfSome(name)
	return self
}

func (self *UserProfileBuilder) WithSurname(surname *nullable.Nullable[string]) *UserProfileBuilder {
	self.surname = nullcommon.CopyPtrIfSome(surname)
	return self
}

func (self *UserProfileBuilder) WithPatronymic(patronymic *nullable.Nullable[string]) *UserProfileBuilder {
	self.patronymic = nullcommon.CopyPtrIfSome(patronymic)
	return self
}

func (self *UserProfileBuilder) WithBirthDate(birthDate *nullable.Nullable[time.Time]) *UserProfileBuilder {
	self.birthDate = nullcommon.CopyPtrIfSome(birthDate)
	return self
}

func (self *UserProfileBuilder) WithPhotoId(photoId *nullable.Nullable[uuid.UUID]) *UserProfileBuilder {
	self.photoId = nullcommon.CopyPtrIfSome(photoId)
	return self
}

func (self *UserProfileBuilder) Build() models.UserProfile {
	return models.UserProfile{
		Id:         self.id,
		UserId:     self.userId,
		Name:       self.name,
		Surname:    self.surname,
		Patronymic: self.patronymic,
		BirthDate:  self.birthDate,
		PhotoId:    self.photoId,
	}
}

type UserFavoritePickUpPointBuilder struct {
	id            uuid.UUID
	userId        uuid.UUID
	pickUpPointId *uuid.UUID
}

func NewUserFavoritePickUpPoint() *UserFavoritePickUpPointBuilder {
	return &UserFavoritePickUpPointBuilder{}
}

func (self *UserFavoritePickUpPointBuilder) WithId(id uuid.UUID) *UserFavoritePickUpPointBuilder {
	self.id = id
	return self
}

func (self *UserFavoritePickUpPointBuilder) WithUserId(userId uuid.UUID) *UserFavoritePickUpPointBuilder {
	self.userId = userId
	return self
}

func (self *UserFavoritePickUpPointBuilder) WithPickUpPointId(pickUpPointId *nullable.Nullable[uuid.UUID]) *UserFavoritePickUpPointBuilder {
	self.pickUpPointId = nullcommon.CopyPtrIfSome(pickUpPointId)
	return self
}

func (self *UserFavoritePickUpPointBuilder) Build() models.UserFavoritePickUpPoint {
	return models.UserFavoritePickUpPoint{
		Id:            self.id,
		UserId:        self.userId,
		PickUpPointId: self.pickUpPointId,
	}
}

type UserPayMethodBuilder struct {
	name     string
	methodId uuid.UUID
	payerId  string
	priority uint
}

func NewUserPayMethod() *UserPayMethodBuilder {
	return &UserPayMethodBuilder{}
}

func (self *UserPayMethodBuilder) WithName(name string) *UserPayMethodBuilder {
	self.name = name
	return self
}

func (self *UserPayMethodBuilder) WithMethodId(methodId uuid.UUID) *UserPayMethodBuilder {
	self.methodId = methodId
	return self
}

func (self *UserPayMethodBuilder) WithPayerId(payerId string) *UserPayMethodBuilder {
	self.payerId = payerId
	return self
}

func (self *UserPayMethodBuilder) WithPriority(priority uint) *UserPayMethodBuilder {
	self.priority = priority
	return self
}

func (self *UserPayMethodBuilder) Build() models.UserPayMethod {
	return models.UserPayMethod{
		Name:     self.name,
		MethodId: self.methodId,
		PayerId:  self.payerId,
		Priority: self.priority,
	}
}

type UserPayMethodsBuilder struct {
	userId uuid.UUID
	mmap   map[uuid.UUID]models.UserPayMethod
}

func NewUserPayMethods() *UserPayMethodsBuilder {
	return &UserPayMethodsBuilder{}
}

func (self *UserPayMethodsBuilder) WithUserId(userId uuid.UUID) *UserPayMethodsBuilder {
	self.userId = userId
	return self
}

func (self *UserPayMethodsBuilder) WithMethods(methods ...models.UserPayMethod) *UserPayMethodsBuilder {
	self.mmap = make(map[uuid.UUID]models.UserPayMethod, len(methods))

	for _, m := range methods {
		self.mmap[m.MethodId] = m
	}

	return self
}

func (self *UserPayMethodsBuilder) Build() models.UserPayMethods {
	return models.UserPayMethods{
		UserId: self.userId,
		Map:    self.mmap,
	}
}

