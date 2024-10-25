package storage

import (
	"database/sql"
	"errors"
	"rent_service/internal/domain/records"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/errors/cmnerrors"
	sqlCollection "rent_service/internal/repository/implementation/sql/collection"
	"rent_service/internal/repository/implementation/sql/exist"
	gen_uuid "rent_service/internal/repository/implementation/sql/generate/uuid"
	"rent_service/internal/repository/implementation/sql/repositories/instance"
	"rent_service/internal/repository/implementation/sql/repositories/pickuppoint"
	"rent_service/internal/repository/implementation/sql/technical"
	"rent_service/internal/repository/interfaces/storage"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	Id            uuid.UUID    `db:"id"`
	PickUpPointId uuid.UUID    `db:"pick_up_point_id"`
	InstanceId    uuid.UUID    `db:"instance_id"`
	InDate        time.Time    `db:"in_date"`
	OutDate       sql.NullTime `db:"out_date"`
	technical.Info
}

type repository struct {
	connection *sqlx.DB
	setter     technical.ISetter
}

func New(connection *sqlx.DB, setter technical.ISetter) storage.IRepository {
	return &repository{connection, setter}
}

func mapf(value *Storage) records.Storage {
	out := records.Storage{
		Id:            value.Id,
		PickUpPointId: value.PickUpPointId,
		InstanceId:    value.InstanceId,
		InDate:        value.InDate,
	}

	if value.OutDate.Valid {
		out.OutDate = new(time.Time)
		*out.OutDate = value.OutDate.Time
	}

	return out
}

func unmapf(value *records.Storage) Storage {
	out := Storage{
		Id:            value.Id,
		PickUpPointId: value.PickUpPointId,
		InstanceId:    value.InstanceId,
		InDate:        value.InDate,
	}

	if nil != value.OutDate {
		out.OutDate.Valid = true
		out.OutDate.Time = *value.OutDate
	}

	return out
}

const insert_query string = `
    insert into records.pick_up_points_instances (
        id, pick_up_point_id, instance_id, in_date, out_date,
        modification_date, modification_source
    ) values (
        :id, :pick_up_point_id, :instance_id, :in_date, :out_date,
        :modification_date, :modification_source
    )
`

func (self *repository) Create(
	storage records.Storage,
) (records.Storage, error) {
	err := pickuppoint.CheckExistsById(self.connection, storage.PickUpPointId)

	if nil == err {
		err = instance.CheckExistsById(self.connection, storage.InstanceId)
	}

	if nil == err && nil == storage.OutDate {
		err = CheckExistsByActiveInstance(self.connection, storage.InstanceId)

		if nil == err {
			err = cmnerrors.Duplicate("storage_active_instance_id")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	if nil == err {
		storage.Id, err = gen_uuid.GenerateAvailable(
			self.connection,
			CheckExistsById,
		)
	}

	if nil == err {
		mapped := unmapf(&storage)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(insert_query, mapped)
	}

	return storage, err
}

const update_query string = `
    update records.pick_up_points_instances
    set pick_up_point_id=:pick_up_point_id, instance_id=:instance_id,
        in_date=:in_date, out_date=:out_date,
        modification_date=:modification_date,
        modification_source=:modification_source
    where id=:id;
`

func (self *repository) Update(storage records.Storage) error {
	err := CheckExistsById(self.connection, storage.Id)

	if nil == err {
		err = pickuppoint.CheckExistsById(
			self.connection,
			storage.PickUpPointId,
		)
	}

	if nil == err {
		err = instance.CheckExistsById(self.connection, storage.InstanceId)
	}

	if nil == err && nil == storage.OutDate {
		err = CheckExistsByActiveInstance(self.connection, storage.InstanceId)

		if nil == err {
			err = cmnerrors.Duplicate("storage_active_instance_id")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	if nil == err {
		mapped := unmapf(&storage)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(update_query, mapped)
	}

	return err
}

const get_active_by_pick_up_point_id_query string = `
    select *
    from records.pick_up_points_instances
    where pick_up_point_id = $1 and out_date is null
    offset $2
`

func (self *repository) GetActiveByPickUpPointId(
	pickUpPointId uuid.UUID,
) (collection.Collection[records.Storage], error) {
	if err := pickuppoint.CheckExistsById(self.connection, pickUpPointId); nil != err {
		return nil, err
	}

	return collection.MapCollection(
		mapf,
		sqlCollection.New[Storage](func(offset uint) (*sqlx.Rows, error) {
			return self.connection.Queryx(
				get_active_by_pick_up_point_id_query,
				pickUpPointId,
				offset,
			)
		}),
	), nil
}

const get_active_by_instance_id_query string = `
    select *
    from records.pick_up_points_instances
    where instance_id = $1 and out_date is null
`

func (self *repository) GetActiveByInstanceId(instanceId uuid.UUID) (records.Storage, error) {
	var out Storage
	err := instance.CheckExistsById(self.connection, instanceId)

	if nil == err {
		err = CheckExistsByActiveInstance(self.connection, instanceId)
	}

	if nil == err {
		err = self.connection.Get(
			&out,
			get_active_by_instance_id_query,
			instanceId,
		)
	}

	return mapf(&out), err
}

const count_by_id_query string = `
    select count(*) from records.pick_up_points_instances where id = $1
`

func CheckExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("storage_id", db, count_by_id_query, id)
}

const count_by_active_instance_query string = `
    select count(*)
    from records.pick_up_points_instances
    where instance_id = $1 and out_date is null
`

func CheckExistsByActiveInstance(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check(
		"storage_active_instance_id",
		db,
		count_by_active_instance_query,
		id,
	)
}

