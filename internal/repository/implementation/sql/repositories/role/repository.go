package role

import (
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/repository/errors/cmnerrors"
	"rent_service/internal/repository/implementation/sql/exist"
	gen_uuid "rent_service/internal/repository/implementation/sql/generate/uuid"
	"rent_service/internal/repository/implementation/sql/repositories/user"
	"rent_service/internal/repository/implementation/sql/technical"
	"rent_service/internal/repository/interfaces/role"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Administrator struct {
	Id     uuid.UUID `db:"id"`
	UserId uuid.UUID `db:"user_id"`
	technical.Info
}

type Renter struct {
	Id     uuid.UUID `db:"id"`
	UserId uuid.UUID `db:"user_id"`
	technical.Info
}

type Storekeeper struct {
	Id            uuid.UUID `db:"id"`
	UserId        uuid.UUID `db:"user_id"`
	PickUpPointId uuid.UUID `db:"pick_up_point_id"`
	technical.Info
}

type administratorRepository struct {
	connection *sqlx.DB
}

func NewAdministrator(connection *sqlx.DB) role.IAdministratorRepository {
	return &administratorRepository{connection}
}

func mapAdministrator(value *Administrator) models.Administrator {
	return models.Administrator{
		Id:     value.Id,
		UserId: value.UserId,
	}
}

const get_administrator_by_user_id_query string = `
    select * from roles.administrators where user_id = $1
`

func (self *administratorRepository) GetByUserId(
	userId uuid.UUID,
) (models.Administrator, error) {
	var administrator Administrator
	err := user.CheckExistsById(self.connection, userId)

	if nil == err {
		err = CheckAdministratorExistsByUserId(self.connection, userId)
	}

	if nil == err {
		err = self.connection.Get(
			&administrator,
			get_administrator_by_user_id_query,
			userId,
		)
	}

	return mapAdministrator(&administrator), err
}

const count_administrators_by_user_id_query string = `
    select count(*) from roles.administrators where user_id = $1
`

func CheckAdministratorExistsByUserId(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check(
		"administrator_user_id",
		db,
		count_administrators_by_user_id_query,
		id,
	)
}

type renterRepository struct {
	connection *sqlx.DB
	setter     technical.ISetter
}

func NewRenter(
	connection *sqlx.DB,
	setter technical.ISetter,
) role.IRenterRepository {
	return &renterRepository{connection, setter}
}

func mapRenter(value *Renter) models.Renter {
	return models.Renter{
		Id:     value.Id,
		UserId: value.UserId,
	}
}

const insert_renter_query string = `
    insert into roles.renters (id, user_id, modification_date, modification_source)
    values (:id, :user_id, :modification_date, :modification_source);
`

func (self *renterRepository) Create(userId uuid.UUID) (models.Renter, error) {
	var renter Renter
	err := user.CheckExistsById(self.connection, userId)

	if nil == err {
		err = CheckRenterExistsByUserId(self.connection, userId)

		if nil == err {
			err = cmnerrors.Duplicate("renter_user_id")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	if nil == err {
		renter.Id, err = gen_uuid.GenerateAvailable(
			self.connection,
			CheckRenterExistsById,
		)
	}

	if nil == err {
		renter.UserId = userId
		self.setter.Update(&renter.Info)
		_, err = self.connection.NamedExec(insert_renter_query, renter)
	}

	return mapRenter(&renter), err
}

const get_renter_by_id_query string = `
    select * from roles.renters where id = $1
`

func (self *renterRepository) GetById(
	renterId uuid.UUID,
) (models.Renter, error) {
	var renter Renter
	err := CheckRenterExistsById(self.connection, renterId)

	if nil == err {
		err = self.connection.Get(
			&renter,
			get_renter_by_id_query,
			renterId,
		)
	}

	return mapRenter(&renter), err
}

const get_renter_by_user_id_query string = `
    select * from roles.renters where user_id = $1
`

func (self *renterRepository) GetByUserId(
	userId uuid.UUID,
) (models.Renter, error) {
	var renter Renter
	err := user.CheckExistsById(self.connection, userId)

	if nil == err {
		err = CheckRenterExistsByUserId(self.connection, userId)
	}

	if nil == err {
		err = self.connection.Get(
			&renter,
			get_renter_by_user_id_query,
			userId,
		)
	}

	return mapRenter(&renter), err
}

var count_renters_by_id_query = exist.GenericCounter("roles.renters")

func CheckRenterExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("renter_id", db, count_renters_by_id_query, id)
}

const count_renters_by_user_id_query string = `
    select count(*) from roles.renters where user_id = $1
`

func CheckRenterExistsByUserId(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check(
		"renter_user_id",
		db,
		count_renters_by_user_id_query,
		id,
	)
}

type storekeeperRepository struct {
	connection *sqlx.DB
}

func NewStorekeeper(connection *sqlx.DB) role.IStorekeeperRepository {
	return &storekeeperRepository{connection}
}

func mapStorekeeper(value *Storekeeper) models.Storekeeper {
	return models.Storekeeper{
		Id:            value.Id,
		UserId:        value.UserId,
		PickUpPointId: value.PickUpPointId,
	}
}

const get_storekeeper_by_user_id_query string = `
    select * from roles.storekeepers where user_id = $1
`

func (self *storekeeperRepository) GetByUserId(
	userId uuid.UUID,
) (models.Storekeeper, error) {
	var storekeeper Storekeeper
	err := user.CheckExistsById(self.connection, userId)

	if nil == err {
		err = CheckStorekeeperExistsByUserId(self.connection, userId)
	}

	if nil == err {
		err = self.connection.Get(
			&storekeeper,
			get_storekeeper_by_user_id_query,
			userId,
		)
	}

	return mapStorekeeper(&storekeeper), err
}

const count_storekeepers_by_user_id_query string = `
    select count(*) from roles.storekeepers where user_id = $1
`

func CheckStorekeeperExistsByUserId(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check(
		"storekeeper_user_id",
		db,
		count_storekeepers_by_user_id_query,
		id,
	)
}

