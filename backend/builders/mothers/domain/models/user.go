package models

import (
	"fmt"
	modelsb "rent_service/builders/domain/models"
	"rent_service/builders/misc/dategen"
	"rent_service/builders/misc/uuidgen"
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	"rent_service/misc/nullable"

	"github.com/google/uuid"
)

func UserRandomId() *modelsb.UserBuilder {
	return modelsb.NewUser().
		WithId(uuidgen.Generate())
}

func UserExample(name string, email string, password string) *modelsb.UserBuilder {
	return UserRandomId().
		WithName(name).
		WithEmail(email).
		WithPassword(password).
		WithToken(name + email + password)
}

func getName(base string, prefix *nullable.Nullable[string]) string {
	return fmt.Sprintf("%v%v", nullable.GetOr(
		nullable.Map(prefix, func(prefix *string) string {
			return fmt.Sprintf("%v_%v", base, *prefix)
		}),
		base,
	), uuidgen.Generate())
}

func packUser(name string) *modelsb.UserBuilder {
	return UserExample(
		name,
		fmt.Sprintf("%v@mail.ru", name),
		fmt.Sprintf("%v_psswd", name),
	)
}

func UserDefault(prefix *nullable.Nullable[string]) *modelsb.UserBuilder {
	return packUser(getName("user", prefix))
}

func UserAdministrator(prefix *nullable.Nullable[string]) *modelsb.UserBuilder {
	return packUser(getName("admin", prefix))
}

func UserRenter(prefix *nullable.Nullable[string]) *modelsb.UserBuilder {
	return packUser(getName("renter", prefix))
}

func UserStorekeeper(prefix *nullable.Nullable[string]) *modelsb.UserBuilder {
	return packUser(getName("storekeeper", prefix))
}

func UserProfileRandomId() *modelsb.UserProfileBuilder {
	return modelsb.NewUserProfile().
		WithId(uuidgen.Generate())
}

func UserProfileEmpty(userId uuid.UUID) *modelsb.UserProfileBuilder {
	return UserProfileRandomId().
		WithUserId(userId)
}

var getBirthDate = dategen.CreateGetter(
	dategen.NewDate(1960, 1, 1),
	dategen.NewDate(2006, 1, 1),
)

func UserProfileFilledNoPhoto(userId uuid.UUID, prefix string) *modelsb.UserProfileBuilder {
	return UserProfileEmpty(userId).
		WithName(nullable.Some("Name" + prefix)).
		WithSurname(nullable.Some("Surname" + prefix)).
		WithPatronymic(nullable.Some("Patronymic" + prefix)).
		WithBirthDate(nullable.Some(getBirthDate()))
}

func UserFavoritePickUpPointRandomId() *modelsb.UserFavoritePickUpPointBuilder {
	return modelsb.NewUserFavoritePickUpPoint().
		WithId(uuidgen.Generate())
}

func UserFavoritePickUpPointEmpty(userId uuid.UUID) *modelsb.UserFavoritePickUpPointBuilder {
	return UserFavoritePickUpPointRandomId().
		WithUserId(userId)
}

func UserFavoritePickUpPointFilled(userId uuid.UUID, pickUpPointId uuid.UUID) *modelsb.UserFavoritePickUpPointBuilder {
	return UserFavoritePickUpPointEmpty(userId).
		WithPickUpPointId(nullable.Some(pickUpPointId))
}

func UserPayMethodExample(prefix string, methodId uuid.UUID, id string) *modelsb.UserPayMethodBuilder {
	return modelsb.NewUserPayMethod().
		WithMethodId(methodId).
		WithName("Example " + prefix).
		WithPayerId(prefix + ":" + id)
}

func UserPayMethodCollect(builders ...*modelsb.UserPayMethodBuilder) []models.UserPayMethod {
	return collection.Collect(
		collection.MapIterator(
			func(pair *collection.Pair[uint, *modelsb.UserPayMethodBuilder]) models.UserPayMethod {
				return pair.B.WithPriority(pair.A).Build()
			},
			collection.EnumerateIterator(
				collection.SliceIterator(builders),
			),
		),
	)
}

func UserPayMethods(userId uuid.UUID, methods ...models.UserPayMethod) *modelsb.UserPayMethodsBuilder {
	return modelsb.NewUserPayMethods().
		WithUserId(userId).
		WithMethods(methods...)
}

