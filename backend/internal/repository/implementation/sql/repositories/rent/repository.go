package rent

import (
	"database/sql"
	"errors"
	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/errors/cmnerrors"
	sqlCollection "rent_service/internal/repository/implementation/sql/collection"
	"rent_service/internal/repository/implementation/sql/exist"
	gen_uuid "rent_service/internal/repository/implementation/sql/generate/uuid"
	"rent_service/internal/repository/implementation/sql/repositories/instance"
	"rent_service/internal/repository/implementation/sql/repositories/period"
	"rent_service/internal/repository/implementation/sql/repositories/pickuppoint"
	"rent_service/internal/repository/implementation/sql/repositories/user"
	"rent_service/internal/repository/implementation/sql/technical"
	"rent_service/internal/repository/implementation/sql/utctime"
	"rent_service/internal/repository/interfaces/rent"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Rent struct {
	Id              uuid.UUID                 `db:"id"`
	UserId          uuid.UUID                 `db:"user_id"`
	InstanceId      uuid.UUID                 `db:"instance_id"`
	StartDate       utctime.UTCTime           `db:"start_date"`
	EndDate         sql.Null[utctime.UTCTime] `db:"end_date"`
	PaymentPeriodId uuid.UUID                 `db:"payment_period_id"`
	technical.Info
}

type Request struct {
	Id               uuid.UUID       `db:"id"`
	InstanceId       uuid.UUID       `db:"instance_id"`
	UserId           uuid.UUID       `db:"user_id"`
	PickUpPointId    uuid.UUID       `db:"pick_up_point_id"`
	PaymentPeriodId  uuid.UUID       `db:"payment_period_id"`
	VerificationCode string          `db:"verification_code"`
	CreateDate       utctime.UTCTime `db:"create_date"`
	technical.Info
}

type Return struct {
	Id               uuid.UUID       `db:"id"`
	InstanceId       uuid.UUID       `db:"instance_id"`
	UserId           uuid.UUID       `db:"user_id"`
	PickUpPointId    uuid.UUID       `db:"pick_up_point_id"`
	RentEndDate      utctime.UTCTime `db:"rent_end_date"`
	VerificationCode string          `db:"verification_code"`
	CreateDate       utctime.UTCTime `db:"create_date"`
	technical.Info
}

type repository struct {
	connection *sqlx.DB
	setter     technical.ISetter
}

func New(connection *sqlx.DB, setter technical.ISetter) rent.IRepository {
	return &repository{connection, setter}
}

func mapf(value *Rent) records.Rent {
	out := records.Rent{
		Id:              value.Id,
		UserId:          value.UserId,
		InstanceId:      value.InstanceId,
		StartDate:       value.StartDate.Time,
		PaymentPeriodId: value.PaymentPeriodId,
	}

	if value.EndDate.Valid {
		out.EndDate = new(time.Time)
		*out.EndDate = value.EndDate.V.Time
	}

	return out
}

func unmapf(value *records.Rent) Rent {
	out := Rent{
		Id:              value.Id,
		UserId:          value.UserId,
		InstanceId:      value.InstanceId,
		StartDate:       utctime.FromTime(value.StartDate),
		PaymentPeriodId: value.PaymentPeriodId,
	}

	if nil != value.EndDate {
		out.EndDate.Valid = true
		out.EndDate.V = utctime.FromTime(*value.EndDate)
	}

	return out
}

const insert_query string = `
    insert into records.users_rents (
        id, user_id, instance_id, start_date, end_date, payment_period_id,
        modification_date, modification_source
    ) values (
        :id, :user_id, :instance_id, :start_date, :end_date, :payment_period_id,
        :modification_date, :modification_source
    )
`

func (self *repository) Create(rent records.Rent) (records.Rent, error) {
	err := user.CheckExistsById(self.connection, rent.UserId)

	if nil == err {
		err = instance.CheckExistsById(self.connection, rent.InstanceId)
	}

	if nil == err {
		err = period.CheckExistsById(self.connection, rent.PaymentPeriodId)
	}

	if nil == err && nil == rent.EndDate {
		err = CheckActiveExistsByInstanceId(self.connection, rent.InstanceId)

		if nil == err {
			err = cmnerrors.Duplicate("active_rent")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	if nil == err {
		rent.Id, err = gen_uuid.GenerateAvailable(
			self.connection,
			CheckExistsById,
		)
	}

	if nil == err {
		mapped := unmapf(&rent)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(insert_query, mapped)
	}

	return rent, err
}

const update_query string = `
    update records.users_rents
    set user_id=:user_id, instance_id=:instance_id, start_date=:start_date,
        end_date=:end_date, payment_period_id=:payment_period_id,
        modification_date=:modification_date,
        modification_source=:modification_source
    where id=:id
`

func (self *repository) Update(rent records.Rent) error {
	err := CheckExistsById(self.connection, rent.Id)

	if nil == err {
		err = user.CheckExistsById(self.connection, rent.UserId)
	}

	if nil == err {
		err = instance.CheckExistsById(self.connection, rent.InstanceId)
	}

	if nil == err {
		err = period.CheckExistsById(self.connection, rent.PaymentPeriodId)
	}

	if nil == err && nil == rent.EndDate {
		err = CheckActiveExistsByInstanceId(self.connection, rent.InstanceId)

		if nil == err {
			err = cmnerrors.Duplicate("active_rent")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	if nil == err {
		mapped := unmapf(&rent)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(update_query, mapped)
	}

	return err
}

const get_by_id_query string = `
    select * from records.users_rents where id = $1
`

func (self *repository) GetById(rentId uuid.UUID) (records.Rent, error) {
	var out Rent
	err := CheckExistsById(self.connection, rentId)

	if nil == err {
		err = self.connection.Get(&out, get_by_id_query, rentId)
	}

	return mapf(&out), err
}

const get_by_user_id_query string = `
    select *
    from records.users_rents
    where user_id = $1
    offset $2
`

func (self *repository) GetByUserId(
	userId uuid.UUID,
) (collection.Collection[records.Rent], error) {
	if err := user.CheckExistsById(self.connection, userId); nil != err {
		return nil, err
	}

	return collection.MapCollection(
		mapf,
		sqlCollection.New[Rent](func(offset uint) (*sqlx.Rows, error) {
			return self.connection.Queryx(get_by_user_id_query, userId, offset)
		}),
	), nil
}

const get_active_by_instance_id_query string = `
    select *
    from records.users_rents
    where end_date is null and instance_id = $1
`

func (self *repository) GetActiveByInstanceId(
	instanceId uuid.UUID,
) (records.Rent, error) {
	var out Rent
	err := instance.CheckExistsById(self.connection, instanceId)

	if nil == err {
		err = CheckActiveExistsByInstanceId(self.connection, instanceId)
	}

	if nil == err {
		err = self.connection.Get(&out, get_active_by_instance_id_query, instanceId)
	}

	return mapf(&out), err
}

const get_past_by_user_id_query string = `
    select *
    from records.users_rents
    where end_date is not null and user_id = $1
    offset $2
`

func (self *repository) GetPastByUserId(
	userId uuid.UUID,
) (collection.Collection[records.Rent], error) {
	if err := user.CheckExistsById(self.connection, userId); nil != err {
		return nil, err
	}

	return collection.MapCollection(
		mapf,
		sqlCollection.New[Rent](func(offset uint) (*sqlx.Rows, error) {
			return self.connection.Queryx(get_past_by_user_id_query, userId, offset)
		}),
	), nil
}

var count_by_id_query string = exist.GenericCounter("records.users_rents")

func CheckExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("rent_id", db, count_by_id_query, id)
}

const count_active_by_instance_id_query string = `
    select count(*)
    from records.users_rents
    where end_date is null and instance_id = $1
`

func CheckActiveExistsByInstanceId(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check(
		"rent_active_instance_id",
		db,
		count_active_by_instance_id_query,
		id,
	)
}

type requestRepository struct {
	connection *sqlx.DB
	setter     technical.ISetter
}

func NewRequest(
	connection *sqlx.DB,
	setter technical.ISetter,
) rent.IRequestRepository {
	return &requestRepository{connection, setter}
}

func mapRequest(value *Request) requests.Rent {
	return requests.Rent{
		Id:               value.Id,
		InstanceId:       value.InstanceId,
		UserId:           value.UserId,
		PickUpPointId:    value.PickUpPointId,
		PaymentPeriodId:  value.PaymentPeriodId,
		VerificationCode: value.VerificationCode,
		CreateDate:       value.CreateDate.Time,
	}
}

func unmapRequest(value *requests.Rent) Request {
	return Request{
		Id:               value.Id,
		InstanceId:       value.InstanceId,
		UserId:           value.UserId,
		PickUpPointId:    value.PickUpPointId,
		PaymentPeriodId:  value.PaymentPeriodId,
		VerificationCode: value.VerificationCode,
		CreateDate:       utctime.FromTime(value.CreateDate),
	}
}

const insert_request_query string = `
    insert into rents.requests (
        id, instance_id, user_id, pick_up_point_id, payment_period_id,
        verification_code, create_date, modification_date, modification_source
    ) values (
        :id, :instance_id, :user_id, :pick_up_point_id, :payment_period_id,
        :verification_code, :create_date, :modification_date,
        :modification_source
    )
`

func (self *requestRepository) Create(
	request requests.Rent,
) (requests.Rent, error) {
	err := instance.CheckExistsById(self.connection, request.InstanceId)

	if nil == err {
		err = user.CheckExistsById(self.connection, request.UserId)
	}

	if nil == err {
		err = pickuppoint.CheckExistsById(self.connection, request.PickUpPointId)
	}

	if nil == err {
		err = period.CheckExistsById(self.connection, request.PaymentPeriodId)
	}

	if nil == err {
		err = CheckRequestExistsByInstanceId(self.connection, request.InstanceId)

		if nil == err {
			err = cmnerrors.Duplicate("active_request")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	if nil == err {
		request.Id, err = gen_uuid.GenerateAvailable(
			self.connection,
			CheckRequestExistsById,
		)
	}

	if nil == err {
		mapped := unmapRequest(&request)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(insert_request_query, mapped)
	}

	return request, err
}

const get_request_by_id_query string = `
    select * from rents.requests where id = $1
`

func (self *requestRepository) GetById(
	requestId uuid.UUID,
) (requests.Rent, error) {
	var out Request
	err := CheckRequestExistsById(self.connection, requestId)

	if nil == err {
		err = self.connection.Get(&out, get_request_by_id_query, requestId)
	}

	return mapRequest(&out), err
}

const get_request_by_user_id_query string = `
    select * from rents.requests where user_id = $1 offset $2
`

func (self *requestRepository) GetByUserId(
	userId uuid.UUID,
) (collection.Collection[requests.Rent], error) {
	if err := user.CheckExistsById(self.connection, userId); nil != err {
		return nil, err
	}

	return collection.MapCollection(
		mapRequest,
		sqlCollection.New[Request](func(offset uint) (*sqlx.Rows, error) {
			return self.connection.Queryx(
				get_request_by_user_id_query,
				userId,
				offset,
			)
		}),
	), nil
}

const get_request_by_instance_id_query string = `
    select * from rents.requests where instance_id = $1
`

func (self *requestRepository) GetByInstanceId(
	instanceId uuid.UUID,
) (requests.Rent, error) {
	var out Request
	err := instance.CheckExistsById(self.connection, instanceId)

	if nil == err {
		err = CheckRequestExistsByInstanceId(self.connection, instanceId)
	}

	if nil == err {
		err = self.connection.Get(
			&out,
			get_request_by_instance_id_query,
			instanceId,
		)
	}

	return mapRequest(&out), err
}

const get_request_by_pick_up_point_id_query string = `
    select * from rents.requests where pick_up_point_id = $1 offset $2
`

func (self *requestRepository) GetByPickUpPointId(
	pickUpPointId uuid.UUID,
) (collection.Collection[requests.Rent], error) {
	if err := pickuppoint.CheckExistsById(self.connection, pickUpPointId); nil != err {
		return nil, err
	}

	return collection.MapCollection(
		mapRequest,
		sqlCollection.New[Request](func(offset uint) (*sqlx.Rows, error) {
			return self.connection.Queryx(
				get_request_by_pick_up_point_id_query,
				pickUpPointId,
				offset,
			)
		}),
	), nil
}

const delete_request_by_id_query string = `
    delete from rents.requests where id = $1
`

func (self *requestRepository) Remove(requestId uuid.UUID) error {
	err := CheckRequestExistsById(self.connection, requestId)

	if nil == err {
		_, err = self.connection.Exec(delete_request_by_id_query, requestId)
	}

	return err
}

var count_request_by_id_query string = exist.GenericCounter("rents.requests")

func CheckRequestExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("rent_request_id", db, count_request_by_id_query, id)
}

const count_request_by_instance_id string = `
    select count(*) from rents.requests where instance_id = $1
`

func CheckRequestExistsByInstanceId(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check(
		"rent_request_instance_id",
		db,
		count_request_by_instance_id,
		id,
	)
}

type returnRepository struct {
	connection *sqlx.DB
	setter     technical.ISetter
}

func NewReturn(
	connection *sqlx.DB,
	setter technical.ISetter,
) rent.IReturnRepository {
	return &returnRepository{connection, setter}
}

func mapReturn(value *Return) requests.Return {
	return requests.Return{
		Id:               value.Id,
		InstanceId:       value.InstanceId,
		UserId:           value.UserId,
		PickUpPointId:    value.PickUpPointId,
		RentEndDate:      value.RentEndDate.Time,
		VerificationCode: value.VerificationCode,
		CreateDate:       value.CreateDate.Time,
	}
}

func unmapReturn(value *requests.Return) Return {
	return Return{
		Id:               value.Id,
		InstanceId:       value.InstanceId,
		UserId:           value.UserId,
		PickUpPointId:    value.PickUpPointId,
		RentEndDate:      utctime.FromTime(value.RentEndDate),
		VerificationCode: value.VerificationCode,
		CreateDate:       utctime.FromTime(value.CreateDate),
	}
}

const insert_return_query string = `
    insert into rents."returns" (
        id, instance_id, user_id, pick_up_point_id, rent_end_date,
        verification_code, create_date, modification_date, modification_source
    ) values (
        :id, :instance_id, :user_id, :pick_up_point_id, :rent_end_date,
        :verification_code, :create_date, :modification_date,
        :modification_source
    )
`

func (self *returnRepository) Create(
	request requests.Return,
) (requests.Return, error) {
	err := instance.CheckExistsById(self.connection, request.InstanceId)

	if nil == err {
		err = user.CheckExistsById(self.connection, request.UserId)
	}

	if nil == err {
		err = pickuppoint.CheckExistsById(self.connection, request.PickUpPointId)
	}

	if nil == err {
		err = CheckReturnExistsByInstanceId(self.connection, request.InstanceId)

		if nil == err {
			err = cmnerrors.Duplicate("active_return")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	if nil == err {
		request.Id, err = gen_uuid.GenerateAvailable(
			self.connection,
			CheckReturnExistsById,
		)
	}

	if nil == err {
		mapped := unmapReturn(&request)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(insert_return_query, mapped)
	}

	return request, err
}

const get_return_by_id_query string = `
    select * from rents."returns" where id = $1
`

func (self *returnRepository) GetById(
	requestId uuid.UUID,
) (requests.Return, error) {
	var out Return
	err := CheckReturnExistsById(self.connection, requestId)

	if nil == err {
		err = self.connection.Get(&out, get_return_by_id_query, requestId)
	}

	return mapReturn(&out), err
}

const get_return_by_user_id_query string = `
    select * from rents."returns" where user_id = $1 offset $2
`

func (self *returnRepository) GetByUserId(
	userId uuid.UUID,
) (collection.Collection[requests.Return], error) {
	if err := user.CheckExistsById(self.connection, userId); nil != err {
		return nil, err
	}

	return collection.MapCollection(
		mapReturn,
		sqlCollection.New[Return](func(offset uint) (*sqlx.Rows, error) {
			return self.connection.Queryx(
				get_return_by_user_id_query,
				userId,
				offset,
			)
		}),
	), nil
}

const get_return_by_instance_id_query string = `
    select * from rents."returns" where instance_id = $1
`

func (self *returnRepository) GetByInstanceId(
	instanceId uuid.UUID,
) (requests.Return, error) {
	var out Return
	err := instance.CheckExistsById(self.connection, instanceId)

	if nil == err {
		err = CheckReturnExistsByInstanceId(self.connection, instanceId)
	}

	if nil == err {
		err = self.connection.Get(
			&out,
			get_return_by_instance_id_query,
			instanceId,
		)
	}

	return mapReturn(&out), err
}

const get_return_by_pick_up_point_id_query string = `
    select * from rents."returns" where pick_up_point_id = $1 offset $2
`

func (self *returnRepository) GetByPickUpPointId(
	pickUpPointId uuid.UUID,
) (collection.Collection[requests.Return], error) {
	if err := pickuppoint.CheckExistsById(self.connection, pickUpPointId); nil != err {
		return nil, err
	}

	return collection.MapCollection(
		mapReturn,
		sqlCollection.New[Return](func(offset uint) (*sqlx.Rows, error) {
			return self.connection.Queryx(
				get_return_by_pick_up_point_id_query,
				pickUpPointId,
				offset,
			)
		}),
	), nil
}

const delete_return_by_id_query string = `
    delete from rents."returns" where id=$1
`

func (self *returnRepository) Remove(
	requestId uuid.UUID,
) error {
	err := CheckReturnExistsById(self.connection, requestId)

	if nil == err {
		_, err = self.connection.Exec(delete_return_by_id_query, requestId)
	}

	return err
}

var count_return_by_id_query string = exist.GenericCounter("rents.\"returns\"")

func CheckReturnExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("rent_return_id", db, count_return_by_id_query, id)
}

const count_return_by_instance_id string = `
    select count(*) from rents."returns" where instance_id = $1
`

func CheckReturnExistsByInstanceId(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check(
		"rent_return_instance_id",
		db,
		count_return_by_instance_id,
		id,
	)
}

