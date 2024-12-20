package role

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/repository/interfaces/role"
	"rent_service/misc/contextholder"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type administratorRepository struct {
	wrapped role.IAdministratorRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewAdministrator(
	wrapped role.IAdministratorRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) role.IAdministratorRepository {
	return &administratorRepository{wrapped, hl, tracer}
}

func (self *administratorRepository) GetByUserId(
	userId uuid.UUID,
) (models.Administrator, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.RoleAdministrator.GetByUserId",
		func(span trace.Span) (models.Administrator, error) {
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return self.wrapped.GetByUserId(userId)
		},
		func(err error) error { return err },
	)
}

type renterRepository struct {
	wrapped role.IRenterRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewRenter(
	wrapped role.IRenterRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) role.IRenterRepository {
	return &renterRepository{wrapped, hl, tracer}
}

func (self *renterRepository) Create(userId uuid.UUID) (models.Renter, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.RoleRenter.Create",
		func(span trace.Span) (models.Renter, error) {
			out, err := self.wrapped.Create(userId)
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return out, err
		},
		func(err error) error { return err },
	)
}

func (self *renterRepository) GetById(
	renterId uuid.UUID,
) (models.Renter, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.RoleRenter.GetById",
		func(span trace.Span) (models.Renter, error) {
			span.SetAttributes(
				attribute.Stringer("RenterId", renterId),
			)
			return self.wrapped.GetById(renterId)
		},
		func(err error) error { return err },
	)
}

func (self *renterRepository) GetByUserId(
	userId uuid.UUID,
) (models.Renter, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.RoleRenter.GetByUserId",
		func(span trace.Span) (models.Renter, error) {
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return self.wrapped.GetByUserId(userId)
		},
		func(err error) error { return err },
	)
}

type storekeeperRepository struct {
	wrapped role.IStorekeeperRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewStorekeeper(
	wrapped role.IStorekeeperRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) role.IStorekeeperRepository {
	return &storekeeperRepository{wrapped, hl, tracer}
}

func (self *storekeeperRepository) GetByUserId(
	userId uuid.UUID,
) (models.Storekeeper, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.RoleStorekeeper.GetByUserId",
		func(span trace.Span) (models.Storekeeper, error) {
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return self.wrapped.GetByUserId(userId)
		},
		func(err error) error { return err },
	)
}

