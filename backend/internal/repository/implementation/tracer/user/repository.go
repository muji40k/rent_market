package user

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/repository/interfaces/user"
	"rent_service/misc/contextholder"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type repository struct {
	wrapped user.IRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped user.IRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) user.IRepository {
	return &repository{wrapped, hl, tracer}
}

type maskedUser struct {
	Id    uuid.UUID
	Name  string
	Email string
}

func maskUser(user models.User) maskedUser {
	return maskedUser{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}
}

func (self *repository) Create(user models.User) (models.User, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.User.Create",
		func(span trace.Span) (models.User, error) {
			out, err := self.wrapped.Create(user)
			span.SetAttributes(
				wrap.AttributeJSON("User", maskUser(out)),
			)
			return out, err
		},
		func(err error) error { return err },
	)
}

func (self *repository) Update(user models.User) error {
	return wrap.SpanCall(self.hl, self.tracer, "Repository.User.Update",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("User", maskUser(user)),
			)
			return self.wrapped.Update(user)
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetById(userId uuid.UUID) (models.User, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.User.GetById",
		func(span trace.Span) (models.User, error) {
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return self.wrapped.GetById(userId)
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetByEmail(email string) (models.User, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.User.GetByEmail",
		func(span trace.Span) (models.User, error) {
			span.SetAttributes(
				attribute.String("Email", email),
			)
			return self.wrapped.GetByEmail(email)
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetByToken(token models.Token) (models.User, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.User.GetByToken",
		func(span trace.Span) (models.User, error) {
			return self.wrapped.GetByToken(token)
		},
		func(err error) error { return err },
	)
}

type profileRepository struct {
	wrapped user.IProfileRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewProfile(
	wrapped user.IProfileRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) user.IProfileRepository {
	return &profileRepository{wrapped, hl, tracer}
}

func (self *profileRepository) Create(
	profile models.UserProfile,
) (models.UserProfile, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.UserProfile.Create",
		func(span trace.Span) (models.UserProfile, error) {
			span.SetAttributes(
				attribute.Stringer("UserId", profile.UserId),
			)
			return self.wrapped.Create(profile)
		},
		func(err error) error { return err },
	)
}

func (self *profileRepository) Update(
	profile models.UserProfile,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Repository.UserProfile.Update",
		func(span trace.Span) error {
			span.SetAttributes(
				attribute.Stringer("UserId", profile.UserId),
			)
			return self.wrapped.Update(profile)
		},
		func(err error) error { return err },
	)
}

func (self *profileRepository) GetByUserId(
	userId uuid.UUID,
) (models.UserProfile, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.UserProfile.GetByUserId",
		func(span trace.Span) (models.UserProfile, error) {
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return self.wrapped.GetByUserId(userId)
		},
		func(err error) error { return err },
	)
}

type favoriteRepository struct {
	wrapped user.IFavoriteRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewFavorite(
	wrapped user.IFavoriteRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) user.IFavoriteRepository {
	return &favoriteRepository{wrapped, hl, tracer}
}

func (self *favoriteRepository) Create(
	profile models.UserFavoritePickUpPoint,
) (models.UserFavoritePickUpPoint, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.UserFavorite.Create",
		func(span trace.Span) (models.UserFavoritePickUpPoint, error) {
			out, err := self.wrapped.Create(profile)
			span.SetAttributes(
				attribute.Stringer("UserId", out.UserId),
			)
			return out, err
		},
		func(err error) error { return err },
	)
}

func (self *favoriteRepository) Update(
	profile models.UserFavoritePickUpPoint,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Repository.UserFavorite.Update",
		func(span trace.Span) error {
			span.SetAttributes(
				attribute.Stringer("UserId", profile.UserId),
			)
			return self.wrapped.Update(profile)
		},
		func(err error) error { return err },
	)
}

func (self *favoriteRepository) GetByUserId(
	userId uuid.UUID,
) (models.UserFavoritePickUpPoint, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.UserFavorite.GetByUserId",
		func(span trace.Span) (models.UserFavoritePickUpPoint, error) {
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return self.wrapped.GetByUserId(userId)
		},
		func(err error) error { return err },
	)
}

type payMethodsRepository struct {
	wrapped user.IPayMethodsRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewPayMethods(
	wrapped user.IPayMethodsRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) user.IPayMethodsRepository {
	return &payMethodsRepository{wrapped, hl, tracer}
}

func (self *payMethodsRepository) CreatePayMethod(
	userId uuid.UUID,
	payMethod models.UserPayMethod,
) (models.UserPayMethods, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.UserPayMethods.CreatePayMethod",
		func(span trace.Span) (models.UserPayMethods, error) {
			out, err := self.wrapped.CreatePayMethod(userId, payMethod)
			span.SetAttributes(
				wrap.AttributeJSON("UserId", userId),
			)
			return out, err
		},
		func(err error) error { return err },
	)
}

func (self *payMethodsRepository) Update(
	payMethods models.UserPayMethods,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Repository.UserPayMethods.Update",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("UserId", payMethods.UserId),
			)
			return self.wrapped.Update(payMethods)
		},
		func(err error) error { return err },
	)
}

func (self *payMethodsRepository) GetByUserId(
	userId uuid.UUID,
) (models.UserPayMethods, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.UserPayMethods.GetByUserId",
		func(span trace.Span) (models.UserPayMethods, error) {
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return self.wrapped.GetByUserId(userId)
		},
		func(err error) error { return err },
	)
}

