package provision

import (
	"database/sql"
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/errors/cmnerrors"
	sqlCollection "rent_service/internal/repository/implementation/sql/collection"
	"rent_service/internal/repository/implementation/sql/exist"
	gen_uuid "rent_service/internal/repository/implementation/sql/generate/uuid"
	"rent_service/internal/repository/implementation/sql/repositories/currency"
	"rent_service/internal/repository/implementation/sql/repositories/instance"
	"rent_service/internal/repository/implementation/sql/repositories/period"
	"rent_service/internal/repository/implementation/sql/repositories/pickuppoint"
	"rent_service/internal/repository/implementation/sql/repositories/product"
	"rent_service/internal/repository/implementation/sql/repositories/role"
	"rent_service/internal/repository/implementation/sql/repositories/user"
	"rent_service/internal/repository/implementation/sql/technical"
	"rent_service/internal/repository/interfaces/provision"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Provision struct {
	Id         uuid.UUID    `db:"id"`
	RenterId   uuid.UUID    `db:"renter_id"`
	InstanceId uuid.UUID    `db:"instance_id"`
	StartDate  time.Time    `db:"start_date"`
	EndDate    sql.NullTime `db:"end_date"`
	technical.Info
}

type Request struct {
	Id               uuid.UUID `db:"id"`
	ProductId        uuid.UUID `db:"product_id"`
	RenterId         uuid.UUID `db:"renter_id"`
	PickUpPointId    uuid.UUID `db:"pick_up_point_id"`
	Name             string    `db:"name"`
	Description      string    `db:"description"`
	Condition        string    `db:"condition"`
	VerificationCode string    `db:"verification_code"`
	CreateDate       time.Time `db:"create_date"`
	technical.Info
}

type RequestPayPlans struct {
	Id         uuid.UUID `db:"id"`
	RequestId  uuid.UUID `db:"request_id"`
	PeriodId   uuid.UUID `db:"period_id"`
	CurrencyId uuid.UUID `db:"currency_id"`
	Value      float64   `db:"value"`
	technical.Info
}

type Revoke struct {
	Id               uuid.UUID `db:"id"`
	InstanceId       uuid.UUID `db:"instance_id"`
	RenterId         uuid.UUID `db:"renter_id"`
	PickUpPointId    uuid.UUID `db:"pick_up_point_id"`
	VerificationCode string    `db:"verification_code"`
	CreateDate       time.Time `db:"create_date"`
	technical.Info
}

type repository struct {
	connection *sqlx.DB
	setter     technical.ISetter
}

func New(connection *sqlx.DB, setter technical.ISetter) provision.IRepository {
	return &repository{connection, setter}
}

func mapf(value *Provision) records.Provision {
	out := records.Provision{
		Id:         value.Id,
		RenterId:   value.RenterId,
		InstanceId: value.InstanceId,
		StartDate:  value.StartDate,
	}

	if value.EndDate.Valid {
		out.EndDate = new(time.Time)
		*out.EndDate = value.EndDate.Time
	}

	return out
}

func unmapf(value *records.Provision) Provision {
	out := Provision{
		Id:         value.Id,
		RenterId:   value.RenterId,
		InstanceId: value.InstanceId,
		StartDate:  value.StartDate,
	}

	if nil != value.EndDate {
		out.EndDate.Valid = true
		out.EndDate.Time = *value.EndDate
	}

	return out
}

const insert_query string = `
    insert into records.renters_instances (
        id, renter_id, instance_id, start_date, end_date, modification_date,
        modification_source
    ) values (
        :id, :renter_id, :instance_id, :start_date, :end_date,
        :modification_date, :modification_source
    )
`

func (self *repository) Create(
	provision records.Provision,
) (records.Provision, error) {
	err := role.CheckRenterExistsById(self.connection, provision.RenterId)

	if nil == err {
		err = instance.CheckExistsById(self.connection, provision.InstanceId)
	}

	if nil == err && nil == provision.EndDate {
		err = CheckActiveExistsByInstanceId(self.connection, provision.InstanceId)

		if nil == err {
			err = cmnerrors.Duplicate("active_provision")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	if nil == err {
		provision.Id, err = gen_uuid.GenerateAvailable(
			self.connection,
			CheckExistsById,
		)
	}

	if nil == err {
		mapped := unmapf(&provision)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(insert_query, mapped)
	}

	return provision, err
}

const update_query string = `
    update records.renters_instances
    set renter_id=:renter_id, instance_id=:instance_id, start_date=:start_date,
        end_date=:end_date, modification_date=:modification_date,
        modification_source=:modification_source
    where id=:id
`

func (self *repository) Update(provision records.Provision) error {
	err := CheckExistsById(self.connection, provision.Id)

	if nil == err {
		err = role.CheckRenterExistsById(self.connection, provision.RenterId)
	}

	if nil == err {
		err = instance.CheckExistsById(self.connection, provision.InstanceId)
	}

	if nil == err && nil == provision.EndDate {
		err = CheckActiveExistsByInstanceId(self.connection, provision.InstanceId)

		if nil == err {
			err = cmnerrors.Duplicate("active_provision")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	if nil == err {
		mapped := unmapf(&provision)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(update_query, mapped)
	}

	return err
}

const get_by_id_query string = `
    select * from records.renters_instances where id = $1
`

func (self *repository) GetById(
	provisionId uuid.UUID,
) (records.Provision, error) {
	var out Provision
	err := CheckExistsById(self.connection, provisionId)

	if nil == err {
		err = self.connection.Get(&out, get_by_id_query, provisionId)
	}

	return mapf(&out), err
}

const get_active_by_instance_id_query string = `
    select *
    form records.renters_instances
    where end_date is null and instance_id = $1
`

func (self *repository) GetActiveByInstanceId(
	instanceId uuid.UUID,
) (records.Provision, error) {
	var out Provision
	err := instance.CheckExistsById(self.connection, instanceId)

	if nil == err {
		err = CheckActiveExistsByInstanceId(self.connection, instanceId)
	}

	if nil == err {
		err = self.connection.Get(&out, get_active_by_instance_id_query, instanceId)
	}

	return mapf(&out), err
}

const get_active_by_renters_user_id_query string = `
    select *
    from records.renters_instances
    where end_date is null
          and renter_id = (select id from roles.renters where user_id = $1)
    offset $2
`

func (self *repository) GetActiveByRenterUserId(
	userId uuid.UUID,
) (collection.Collection[records.Provision], error) {
	err := user.CheckExistsById(self.connection, userId)

	if nil == err {
		err = role.CheckRenterExistsByUserId(self.connection, userId)
	}

	if nil != err {
		return nil, err
	}

	return collection.MapCollection(
		mapf,
		sqlCollection.New[Provision](func(offset uint) (*sqlx.Rows, error) {
			return self.connection.Queryx(
				get_active_by_renters_user_id_query,
				userId,
				offset,
			)
		}),
	), nil
}

var count_by_id_query string = exist.GenericCounter("records.renters_instances")

func CheckExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("provision_id", db, count_by_id_query, id)
}

const count_active_by_instance_id_query string = `
    select count(*)
    from records.renters_instances
    where end_date is null and instance_id = $1
`

func CheckActiveExistsByInstanceId(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check(
		"provision_active_instance_id",
		db,
		count_active_by_instance_id_query,
		id,
	)
}

type requestRepository struct {
	connection *sqlx.DB
	setter     technical.ISetter
	currency   *currency.Repository
}

func NewRequest(
	connection *sqlx.DB,
	setter technical.ISetter,
	currency *currency.Repository,
) provision.IRequestRepository {
	return &requestRepository{connection, setter, currency}
}

func (self *requestRepository) mapRequest(
	value *Request,
	payPlans []RequestPayPlans,
) requests.Provide {
	out := requests.Provide{
		Id:               value.Id,
		ProductId:        value.ProductId,
		RenterId:         value.RenterId,
		PickUpPointId:    value.PickUpPointId,
		PayPlans:         make(map[uuid.UUID]models.PayPlan),
		Name:             value.Name,
		Description:      value.Description,
		Condition:        value.Condition,
		VerificationCode: value.VerificationCode,
		CreateDate:       value.CreateDate,
	}

	for _, plan := range payPlans {
		price, _ := self.currency.GetById(plan.CurrencyId)
		price.Value = plan.Value

		out.PayPlans[plan.PeriodId] = models.PayPlan{
			Id:       plan.Id,
			PeriodId: plan.PeriodId,
			Price:    price,
		}
	}

	return out
}

func (self *requestRepository) unmapRequest(
	value *requests.Provide,
) (Request, []RequestPayPlans, error) {
	request := Request{
		Id:               value.Id,
		ProductId:        value.ProductId,
		RenterId:         value.RenterId,
		PickUpPointId:    value.PickUpPointId,
		Name:             value.Name,
		Description:      value.Description,
		Condition:        value.Condition,
		VerificationCode: value.VerificationCode,
		CreateDate:       value.CreateDate,
	}

	payPlans := make([]RequestPayPlans, 0, len(value.PayPlans))

	for _, plan := range value.PayPlans {
		cid, err := self.currency.GetId(plan.Price.Name)

		if nil != err {
			return Request{}, nil, err
		}

		payPlans = append(payPlans, RequestPayPlans{
			Id:         plan.Id,
			RequestId:  value.Id,
			PeriodId:   plan.PeriodId,
			CurrencyId: cid,
			Value:      plan.Price.Value,
		})
	}

	return request, payPlans, nil
}

const insert_request_query string = `
    insert into provisions.requests (
        id, product_id, renter_id, pick_up_point_id, "name", description,
        "condition", verification_code, create_date, modification_date,
        modification_source
    ) values (
        :id, :product_id, :renter_id, :pick_up_point_id, :name, :description,
        :condition, :verification_code, :create_date, :modification_date,
        :modification_source
    )
`

const insert_request_pay_plan_query string = `
    insert into provisions.requests_pay_plans (
        id, request_id, period_id, currency_id, value, modification_date,
        modification_source
    ) values (
        :id, :request_id, :period_id, :currency_id, :value, :modification_date,
        :modification_source
    )
`

func (self *requestRepository) Create(
	request requests.Provide,
) (requests.Provide, error) {
	var mr Request
	var mrpp []RequestPayPlans
	err := product.CheckExistsById(self.connection, request.ProductId)

	if nil == err {
		err = role.CheckRenterExistsById(self.connection, request.RenterId)
	}

	if nil == err {
		err = pickuppoint.CheckExistsById(self.connection, request.PickUpPointId)
	}

	if nil == err {
		for k, _ := range request.PayPlans {
			if nil == err {
				err = period.CheckExistsById(self.connection, k)
			}
		}
	}

	if nil == err {
		request.Id, err = gen_uuid.GenerateAvailable(
			self.connection,
			CheckRequestExistsById,
		)
	}

	if nil == err {
		for k, v := range request.PayPlans {
			if nil == err {
				v.Id, err = gen_uuid.GenerateAvailable(
					self.connection,
					CheckRequestPayPlanExistsById,
				)
			}

			if nil == err {
				request.PayPlans[k] = v
			}
		}
	}

	if nil == err {
		mr, mrpp, err = self.unmapRequest(&request)
	}

	if nil == err {
		self.setter.Update(&mr.Info)
		_, err = self.connection.NamedExec(insert_request_query, mr)
	}

	for i := 0; nil == err && len(mrpp) > i; i++ {
		self.setter.Update(&mrpp[i].Info)
		_, err = self.connection.NamedExec(
			insert_request_pay_plan_query,
			mrpp[i],
		)
	}

	return request, err
}

const get_request_by_id_query string = `
    select * from provisions.requests where id = $1
`

const get_request_pay_plan_by_request_id_query string = `
    select * from provisions.requests_pay_plans where request_id = $1
`

func (self *requestRepository) GetById(
	requestId uuid.UUID,
) (requests.Provide, error) {
	var out Request
	var outpp []RequestPayPlans
	err := CheckRequestExistsById(self.connection, requestId)

	if nil == err {
		err = self.connection.Get(&out, get_request_by_id_query, requestId)
	}

	if nil == err {
		err = self.connection.Select(
			&outpp,
			get_request_pay_plan_by_request_id_query,
			requestId,
		)
	}

	return self.mapRequest(&out, outpp), err
}

const get_request_by_renters_user_id_query string = `
    select *
    from provisions.requests
    where renter_id = (select id from roles.renters where user_id = $1)
    offset $2
`

func (self *requestRepository) GetByUserId(
	userId uuid.UUID,
) (collection.Collection[requests.Provide], error) {
	err := user.CheckExistsById(self.connection, userId)

	if nil == err {
		err = role.CheckRenterExistsByUserId(self.connection, userId)
	}

	if nil != err {
		return nil, err
	}

	return self.getRequest(func(offset uint) (*sqlx.Rows, error) {
		return self.connection.Queryx(
			get_request_by_renters_user_id_query,
			userId,
			offset,
		)
	}), nil
}

func (self *requestRepository) GetByInstanceId(
	instanceId uuid.UUID,
) (requests.Provide, error) {
	return requests.Provide{}, errors.New(
		"Этого здесь быть не должно... но останется, так как какая разница?",
	)
}

const get_request_by_pick_up_point_id_query string = `
    select * from provisions.requests where pick_up_point_id = $1 offset $2
`

func (self *requestRepository) GetByPickUpPointId(
	pickUpPointId uuid.UUID,
) (collection.Collection[requests.Provide], error) {
	if err := pickuppoint.CheckExistsById(self.connection, pickUpPointId); nil != err {
		return nil, err
	}

	return self.getRequest(func(offset uint) (*sqlx.Rows, error) {
		return self.connection.Queryx(
			get_request_by_pick_up_point_id_query,
			pickUpPointId,
			offset,
		)
	}), nil
}

const delete_request_by_id_query string = `
    delete from provisions.requests where id = $1
`

const delete_request_pay_plan_by_request_id_query string = `
    delete from provisions.requests_pay_plans where request_id = $1
`

func (self *requestRepository) Remove(requestId uuid.UUID) error {
	err := CheckRequestExistsById(self.connection, requestId)

	if nil == err {
		_, err1 := self.connection.Exec(delete_request_by_id_query, requestId)
		_, err2 := self.connection.Exec(
			delete_request_pay_plan_by_request_id_query,
			requestId,
		)

		if nil != err1 {
			err = err1
		} else if nil != err2 {
			err = err2
		}
	}

	return err
}

var count_request_by_id_query string = exist.GenericCounter("provisions.requests")

func CheckRequestExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check(
		"provision_request_id",
		db,
		count_request_by_id_query,
		id,
	)
}

var count_request_pay_plan_by_id_query string = exist.GenericCounter("provisions.requests_pay_plans")

func CheckRequestPayPlanExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check(
		"provision_request_pay_plan_id",
		db,
		count_request_pay_plan_by_id_query,
		id,
	)
}

func (self *requestRepository) getRequest(
	query sqlCollection.Query,
) collection.Collection[requests.Provide] {
	return collection.MapCollection(
		func(pair *collection.Pair[*Request, []RequestPayPlans]) requests.Provide {
			return self.mapRequest(pair.A, pair.B)
		},
		collection.MapCollection(
			func(rq *Request) collection.Pair[*Request, []RequestPayPlans] {
				var plans []RequestPayPlans
				err := self.connection.Select(
					&plans,
					get_request_pay_plan_by_request_id_query,
					rq.Id,
				)

				if nil != err {
					plans = nil
				}

				return collection.Pair[*Request, []RequestPayPlans]{
					A: rq,
					B: plans,
				}
			},
			sqlCollection.New[Request](query),
		),
	)
}

type revokeRepository struct {
	connection *sqlx.DB
	setter     technical.ISetter
}

func NewRevoke(
	connection *sqlx.DB,
	setter technical.ISetter,
) provision.IRevokeRepository {
	return &revokeRepository{connection, setter}
}

func mapRevoke(value *Revoke) requests.Revoke {
	return requests.Revoke{
		Id:               value.Id,
		InstanceId:       value.InstanceId,
		RenterId:         value.RenterId,
		PickUpPointId:    value.PickUpPointId,
		VerificationCode: value.VerificationCode,
		CreateDate:       value.CreateDate,
	}
}

func unmapRevoke(value *requests.Revoke) Revoke {
	return Revoke{
		Id:               value.Id,
		InstanceId:       value.InstanceId,
		RenterId:         value.RenterId,
		PickUpPointId:    value.PickUpPointId,
		VerificationCode: value.VerificationCode,
		CreateDate:       value.CreateDate,
	}
}

const insert_revoke_query string = `
    insert into provisions.revokes (
        id, instance_id, renter_id, pick_up_point_id, verification_code,
        create_date, modification_date, modification_source
    ) values (
        :id, :instance_id, :renter_id, :pick_up_point_id, :verification_code,
        :create_date, :modification_date, :modification_source
    )
`

func (self *revokeRepository) Create(
	request requests.Revoke,
) (requests.Revoke, error) {
	err := instance.CheckExistsById(self.connection, request.InstanceId)

	if nil == err {
		err = role.CheckRenterExistsById(self.connection, request.RenterId)
	}

	if nil == err {
		err = pickuppoint.CheckExistsById(self.connection, request.PickUpPointId)
	}

	if nil == err {
		err = CheckRevokeExistsByInstanceId(self.connection, request.InstanceId)

		if nil == err {
			err = cmnerrors.Duplicate("active_revoke")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	if nil == err {
		request.Id, err = gen_uuid.GenerateAvailable(
			self.connection,
			CheckRevokeExistsById,
		)
	}

	if nil == err {
		mapped := unmapRevoke(&request)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(insert_revoke_query, mapped)
	}

	return request, err
}

const get_revoke_by_id_query string = `
    select * from provisions.revokes where id = $1
`

func (self *revokeRepository) GetById(
	requestId uuid.UUID,
) (requests.Revoke, error) {
	var out Revoke
	err := CheckRevokeExistsById(self.connection, requestId)

	if nil == err {
		err = self.connection.Get(&out, get_revoke_by_id_query, requestId)
	}

	return mapRevoke(&out), err
}

const get_revoke_by_renters_user_id_query string = `
    select *
    from provisions.revokes
    where renter_id = (select id from roles.renters where user_id = $1)
    offset $2
`

func (self *revokeRepository) GetByUserId(
	userId uuid.UUID,
) (collection.Collection[requests.Revoke], error) {
	err := user.CheckExistsById(self.connection, userId)

	if nil == err {
		err = role.CheckRenterExistsByUserId(self.connection, userId)
	}

	if nil != err {
		return nil, err
	}

	return collection.MapCollection(
		mapRevoke,
		sqlCollection.New[Revoke](func(offset uint) (*sqlx.Rows, error) {
			return self.connection.Queryx(
				get_revoke_by_renters_user_id_query,
				userId,
				offset,
			)
		}),
	), nil
}

const get_revoke_by_instance_id_query string = `
    select * from provisions.revokes where instance_id = $1
`

func (self *revokeRepository) GetByInstanceId(
	instanceId uuid.UUID,
) (requests.Revoke, error) {
	var out Revoke
	err := instance.CheckExistsById(self.connection, instanceId)

	if nil == err {
		err = CheckRevokeExistsByInstanceId(self.connection, instanceId)
	}

	if nil == err {
		err = self.connection.Get(
			&out,
			get_revoke_by_instance_id_query,
			instanceId,
		)
	}

	return mapRevoke(&out), err
}

const get_revoke_by_pick_up_point_id_query string = `
    select * from provisions.revokes where pick_up_point_id = $1 offset $2
`

func (self *revokeRepository) GetByPickUpPointId(
	pickUpPointId uuid.UUID,
) (collection.Collection[requests.Revoke], error) {
	if err := pickuppoint.CheckExistsById(self.connection, pickUpPointId); nil != err {
		return nil, err
	}

	return collection.MapCollection(
		mapRevoke,
		sqlCollection.New[Revoke](func(offset uint) (*sqlx.Rows, error) {
			return self.connection.Queryx(
				get_revoke_by_pick_up_point_id_query,
				pickUpPointId,
				offset,
			)
		}),
	), nil
}

const delete_revoke_by_id_query string = `
    delete from provisions.revokes where id = $1
`

func (self *revokeRepository) Remove(requestId uuid.UUID) error {
	err := CheckRevokeExistsById(self.connection, requestId)

	if nil == err {
		_, err = self.connection.Exec(delete_revoke_by_id_query, requestId)
	}

	return err
}

var count_revoke_by_id_query string = exist.GenericCounter("provisions.revokes")

func CheckRevokeExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("provision", db, count_revoke_by_id_query, id)
}

const count_revoke_by_instance_id string = `
    select count(*) from provisions.revokes where instance_id = $1
`

func CheckRevokeExistsByInstanceId(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check(
		"provision_revoke_instance_id",
		db,
		count_revoke_by_instance_id,
		id,
	)
}

