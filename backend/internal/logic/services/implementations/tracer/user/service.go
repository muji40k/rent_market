package user

import (
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/user"
	"rent_service/internal/logic/services/types/token"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/misc/contextholder"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type service struct {
	wrapped user.IService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped user.IService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) user.IService {
	return &service{wrapped, hl, tracer}
}

func (self *service) GetSelfUserInfo(token token.Token) (user.Info, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetSelfUserInfo",
		func(span trace.Span) (user.Info, error) {
			return self.wrapped.GetSelfUserInfo(token)
		},
		cmnerrors.Internal,
	)
}

func (self *service) UpdateSelfUserInfo(
	token token.Token,
	form user.UpdateForm,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.UpdateSelfUserInfo",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("UpdateForm", form),
			)
			return self.wrapped.UpdateSelfUserInfo(token, form)
		},
		cmnerrors.Internal,
	)
}

func (self *service) UpdateSelfUserPassword(
	token token.Token,
	old_password string,
	new_password string,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.UpdateSelfUserPassword",
		func(_ trace.Span) error {
			return self.wrapped.UpdateSelfUserPassword(
				token,
				old_password,
				new_password,
			)
		},
		cmnerrors.Internal,
	)
}

type profileService struct {
	wrapped user.IProfileService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewProfile(
	wrapped user.IProfileService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) user.IProfileService {
	return &profileService{wrapped, hl, tracer}
}

func (self *profileService) GetUserProfile(
	userId uuid.UUID,
) (user.UserProfile, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetUserProfile",
		func(span trace.Span) (user.UserProfile, error) {
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return self.wrapped.GetUserProfile(userId)
		},
		cmnerrors.Internal,
	)
}

func (self *profileService) GetSelfUserProfile(
	token token.Token,
) (user.UserProfile, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetSelfUserProfile",
		func(span trace.Span) (user.UserProfile, error) {
			return self.wrapped.GetSelfUserProfile(token)
		},
		cmnerrors.Internal,
	)
}

func (self *profileService) UpdateSelfUserProfile(
	token token.Token,
	data user.UserProfile,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.UpdateSelfUserProfile",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("Data", data),
			)
			return self.wrapped.UpdateSelfUserProfile(token, data)
		},
		cmnerrors.Internal,
	)
}

type favoriteService struct {
	wrapped user.IFavoriteService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewFavorite(
	wrapped user.IFavoriteService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) user.IFavoriteService {
	return &favoriteService{wrapped, hl, tracer}
}

func (self *favoriteService) GetUserFavorite(
	userId uuid.UUID,
) (user.UserFavoritePickUpPoint, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetUserFavorite",
		func(span trace.Span) (user.UserFavoritePickUpPoint, error) {
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return self.wrapped.GetUserFavorite(userId)
		},
		cmnerrors.Internal,
	)
}

func (self *favoriteService) GetSelfUserFavorite(
	token token.Token,
) (user.UserFavoritePickUpPoint, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetSelfUserFavorite",
		func(span trace.Span) (user.UserFavoritePickUpPoint, error) {
			return self.wrapped.GetSelfUserFavorite(token)
		},
		cmnerrors.Internal,
	)
}

func (self *favoriteService) UpdateSelfUserFavorite(
	token token.Token,
	data user.UserFavoritePickUpPoint,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.UpdateSelfUserFavorite",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("Data", data),
			)
			return self.wrapped.UpdateSelfUserFavorite(token, data)
		},
		cmnerrors.Internal,
	)
}

type roleService struct {
	wrapped user.IRoleService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewRole(
	wrapped user.IRoleService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) user.IRoleService {
	return &roleService{wrapped, hl, tracer}
}

func (self *roleService) IsRenter(token token.Token) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.IsRenter",
		func(span trace.Span) error {
			return self.wrapped.IsRenter(token)
		},
		cmnerrors.Internal,
	)
}

func (self *roleService) RegisterAsRenter(token token.Token) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.RegisterAsRenter",
		func(span trace.Span) error {
			return self.wrapped.RegisterAsRenter(token)
		},
		cmnerrors.Internal,
	)
}

func (self *roleService) IsAdministrator(token token.Token) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.IsAdministrator",
		func(span trace.Span) error {
			return self.wrapped.IsAdministrator(token)
		},
		cmnerrors.Internal,
	)
}

func (self *roleService) IsStoreKeeper(
	token token.Token,
) (user.StoreKeeper, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.IsStoreKeeper",
		func(span trace.Span) (user.StoreKeeper, error) {
			return self.wrapped.IsStoreKeeper(token)
		},
		cmnerrors.Internal,
	)
}

