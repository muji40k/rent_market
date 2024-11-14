package delivery

import (
	"database/sql"
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/requests"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/errors/cmnerrors"
	sqlCollection "rent_service/internal/repository/implementation/sql/collection"
	"rent_service/internal/repository/implementation/sql/exist"
	gen_uuid "rent_service/internal/repository/implementation/sql/generate/uuid"
	"rent_service/internal/repository/implementation/sql/repositories/instance"
	"rent_service/internal/repository/implementation/sql/repositories/pickuppoint"
	"rent_service/internal/repository/implementation/sql/technical"
	"rent_service/internal/repository/implementation/sql/utctime"
	"rent_service/internal/repository/interfaces/delivery"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Delivery struct {
	Id                 uuid.UUID                 `db:"id"`
	CompanyId          uuid.UUID                 `db:"company_id"`
	InstanceId         uuid.UUID                 `db:"instance_id"`
	FromId             uuid.UUID                 `db:"from_id"`
	ToId               uuid.UUID                 `db:"to_id"`
	DeliveryId         string                    `db:"delivery_id"`
	ScheduledBeginDate utctime.UTCTime           `db:"scheduled_begin_date"`
	ActualBeginDate    sql.Null[utctime.UTCTime] `db:"actual_begin_date"`
	ScheduledEndDate   utctime.UTCTime           `db:"scheduled_end_date"`
	ActualEndDate      sql.Null[utctime.UTCTime] `db:"actual_end_date"`
	VerificationCode   string                    `db:"verification_code"`
	CreateDate         utctime.UTCTime           `db:"create_date"`
	technical.Info
}

type DeliveryCompany struct {
	Id          uuid.UUID      `db:"id"`
	Name        string         `db:"name"`
	Site        sql.NullString `db:"site"`
	PhoneNumber sql.NullString `db:"phone_bumber"`
	Description sql.NullString `db:"description"`
	technical.Info
}

type repository struct {
	connection *sqlx.DB
	setter     technical.ISetter
}

func New(connection *sqlx.DB, setter technical.ISetter) delivery.IRepository {
	return &repository{connection, setter}
}

func mapf(value *Delivery) requests.Delivery {
	out := requests.Delivery{
		Id:                 value.Id,
		CompanyId:          value.CompanyId,
		InstanceId:         value.InstanceId,
		FromId:             value.FromId,
		ToId:               value.ToId,
		DeliveryId:         value.DeliveryId,
		ScheduledBeginDate: value.ScheduledBeginDate.Time,
		ScheduledEndDate:   value.ScheduledEndDate.Time,
		VerificationCode:   value.VerificationCode,
		CreateDate:         value.CreateDate.Time,
	}

	if value.ActualBeginDate.Valid {
		out.ActualBeginDate = new(time.Time)
		*out.ActualBeginDate = value.ActualBeginDate.V.Time
	}

	if value.ActualEndDate.Valid {
		out.ActualEndDate = new(time.Time)
		*out.ActualEndDate = value.ActualEndDate.V.Time
	}

	return out
}

func unmapf(value *requests.Delivery) Delivery {
	out := Delivery{
		Id:                 value.Id,
		CompanyId:          value.CompanyId,
		InstanceId:         value.InstanceId,
		FromId:             value.FromId,
		ToId:               value.ToId,
		DeliveryId:         value.DeliveryId,
		ScheduledBeginDate: utctime.FromTime(value.ScheduledBeginDate),
		ScheduledEndDate:   utctime.FromTime(value.ScheduledEndDate),
		VerificationCode:   value.VerificationCode,
		CreateDate:         utctime.FromTime(value.CreateDate),
	}

	if nil != value.ActualBeginDate {
		out.ActualBeginDate.Valid = true
		out.ActualBeginDate.V = utctime.FromTime(*value.ActualBeginDate)
	}

	if nil != value.ActualEndDate {
		out.ActualEndDate.Valid = true
		out.ActualEndDate.V = utctime.FromTime(*value.ActualEndDate)
	}

	return out
}

func commonChecks(db *sqlx.DB, delivery *requests.Delivery) error {
	var err error

	if nil == delivery.ActualBeginDate && nil != delivery.ActualEndDate {
		err = errors.New("Delivery broken state...")
	}

	if nil == err {
		err = CheckCompanyExistsById(db, delivery.CompanyId)
	}

	if nil == err {
		err = instance.CheckExistsById(db, delivery.InstanceId)
	}

	if nil == err {
		err = pickuppoint.CheckExistsById(db, delivery.FromId)
	}

	if nil == err {
		err = pickuppoint.CheckExistsById(db, delivery.ToId)
	}

	return err
}

const insert_query string = `
    insert into deliveries.deliveries (
        id, company_id, instance_id, from_id, to_id, delivery_id,
        scheduled_begin_date, actual_begin_date, scheduled_end_date,
        actual_end_date, verification_code, create_date, modification_date,
        modification_source
    ) values (
        :id, :company_id, :instance_id, :from_id, :to_id, :delivery_id,
        :scheduled_begin_date, :actual_begin_date, :scheduled_end_date,
        :actual_end_date, :verification_code, :create_date, :modification_date,
        :modification_source
    )
`

func (self *repository) Create(
	delivery requests.Delivery,
) (requests.Delivery, error) {
	err := commonChecks(self.connection, &delivery)

	if nil == err && nil == delivery.ActualEndDate {
		err = CheckExistsActiveByInstanceId(
			self.connection,
			delivery.InstanceId,
		)

		if nil == err {
			err = cmnerrors.Duplicate("delivery_active_instance_id")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	if nil == err {
		delivery.Id, err = gen_uuid.GenerateAvailable(
			self.connection,
			CheckExistsById,
		)
	}

	if nil == err {
		mapped := unmapf(&delivery)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(insert_query, mapped)
	}

	return delivery, err
}

const update_query string = `
    update deliveries.deliveries
    set company_id=:company_id, instance_id=:instance_id, from_id=:from_id,
        to_id=:to_id, delivery_id=:delivery_id,
        scheduled_begin_date=:scheduled_begin_date,
        actual_begin_date=:actual_begin_date,
        scheduled_end_date=:scheduled_end_date,
        actual_end_date=:actual_end_date, verification_code=:verification_code,
        create_date=:create_date, modification_date=:modification_date,
        modification_source=:modification_source
    where id=:id;
`

func (self *repository) Update(delivery requests.Delivery) error {
	err := CheckExistsById(self.connection, delivery.Id)

	if nil == err {
		err = commonChecks(self.connection, &delivery)
	}

	if nil == err {
		mapped := unmapf(&delivery)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(update_query, mapped)
	}

	return err
}

const get_by_id_query string = `
    select * from deliveries.deliveries where id = $1
`

func (self *repository) GetById(
	deliveryId uuid.UUID,
) (requests.Delivery, error) {
	var delivery Delivery
	err := CheckExistsById(self.connection, deliveryId)

	if nil == err {
		err = self.connection.Get(&delivery, get_by_id_query, deliveryId)
	}

	return mapf(&delivery), err
}

const get_active_by_pick_up_point_id_query string = `
    select *
    from deliveries.deliveries
    where (from_id = $1 or to_id = $1) and actual_end_date is null
    offset $2
`

func (self *repository) GetActiveByPickUpPointId(
	pickUpPointId uuid.UUID,
) (collection.Collection[requests.Delivery], error) {
	if err := pickuppoint.CheckExistsById(self.connection, pickUpPointId); nil != err {
		return nil, err
	}

	return collection.MapCollection(
		mapf,
		sqlCollection.New[Delivery](func(offset uint) (*sqlx.Rows, error) {
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
    from deliveries.deliveries
    where instance_id = $1 and actual_end_date is null
`

func (self *repository) GetActiveByInstanceId(
	instanceId uuid.UUID,
) (requests.Delivery, error) {
	var delivery Delivery
	err := instance.CheckExistsById(self.connection, instanceId)

	if nil == err {
		err = CheckExistsActiveByInstanceId(self.connection, instanceId)
	}

	if nil == err {
		err = self.connection.Get(
			&delivery,
			get_active_by_instance_id_query,
			instanceId,
		)
	}

	return mapf(&delivery), err
}

var count_by_id_query string = exist.GenericCounter("deliveries.deliveries")

func CheckExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("delivery_id", db, count_by_id_query, id)
}

var count_active_by_instance_id_query string = `
    select count(*)
    from deliveries.deliveries
    where instance_id = $1 and actual_end_date is null
`

func CheckExistsActiveByInstanceId(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check(
		"delivery_instance_id",
		db,
		count_active_by_instance_id_query,
		id,
	)
}

type companyRepository struct {
	connection *sqlx.DB
}

func NewCompany(connection *sqlx.DB) delivery.ICompanyRepository {
	return &companyRepository{connection}
}

func mapCompany(value *DeliveryCompany) models.DeliveryCompany {
	out := models.DeliveryCompany{
		Id:   value.Id,
		Name: value.Name,
	}

	if value.Site.Valid {
		out.Site = value.Site.String
	}

	if value.PhoneNumber.Valid {
		out.PhoneNumber = value.PhoneNumber.String
	}

	if value.Description.Valid {
		out.Description = value.Description.String
	}

	return out
}

const get_company_by_id_query string = `
    select * from deliveries.companies where id = $1
`

func (self *companyRepository) GetById(
	companyId uuid.UUID,
) (models.DeliveryCompany, error) {
	var out DeliveryCompany
	err := CheckExistsById(self.connection, companyId)

	if nil == err {
		err = self.connection.Get(&out, get_company_by_id_query, companyId)
	}

	return mapCompany(&out), err
}

const get_company_all_query string = `
    select * from deliveries.companies offset $1
`

func (self *companyRepository) GetAll() (collection.Collection[models.DeliveryCompany], error) {
	return collection.MapCollection(
		mapCompany,
		sqlCollection.New[DeliveryCompany](func(offset uint) (*sqlx.Rows, error) {
			return self.connection.Queryx(get_company_all_query, offset)
		}),
	), nil
}

var count_company_by_id_query string = exist.GenericCounter("deliveries.companies")

func CheckCompanyExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check(
		"delivery_company_id",
		db,
		count_company_by_id_query,
		id,
	)
}

