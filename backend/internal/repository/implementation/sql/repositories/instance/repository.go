package instance

import (
	"errors"
	"fmt"
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/errors/cmnerrors"
	sqlCollection "rent_service/internal/repository/implementation/sql/collection"
	"rent_service/internal/repository/implementation/sql/exist"
	gen_uuid "rent_service/internal/repository/implementation/sql/generate/uuid"
	"rent_service/internal/repository/implementation/sql/repositories/currency"
	"rent_service/internal/repository/implementation/sql/repositories/period"
	"rent_service/internal/repository/implementation/sql/repositories/photo"
	"rent_service/internal/repository/implementation/sql/repositories/product"
	"rent_service/internal/repository/implementation/sql/technical"
	"rent_service/internal/repository/interfaces/instance"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Instance struct {
	Id          uuid.UUID `db:"id"`
	ProductId   uuid.UUID `db:"product_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Condition   string    `db:"condition"`
	technical.Info
}

type InstancePayPlan struct {
	Id         uuid.UUID `db:"id"`
	InstanceId uuid.UUID `db:"instance_id"`
	PeriodId   uuid.UUID `db:"period_id"`
	CurrencyId uuid.UUID `db:"currency_id"`
	Price      float64   `db:"price"`
	technical.Info
}

type repository struct {
	connection *sqlx.DB
	setter     technical.ISetter
}

func New(connection *sqlx.DB, setter technical.ISetter) instance.IRepository {
	return &repository{connection, setter}
}

func mapf(value *Instance) models.Instance {
	return models.Instance{
		Id:          value.Id,
		ProductId:   value.ProductId,
		Name:        value.Name,
		Description: value.Description,
		Condition:   value.Condition,
	}
}

func unmapf(value *models.Instance) Instance {
	return Instance{
		Id:          value.Id,
		ProductId:   value.ProductId,
		Name:        value.Name,
		Description: value.Description,
		Condition:   value.Condition,
	}
}

const insert_query string = `
    insert into instances.instances (
        id, product_id, "name", description, "condition", modification_date,
        modification_source
    ) values (
        :id, :product_id, :name, :description, :condition, :modification_date,
        :modification_source
    )
`

func (self *repository) Create(
	instance models.Instance,
) (models.Instance, error) {
	err := product.CheckExistsById(self.connection, instance.ProductId)

	if nil == err {
		instance.Id, err = gen_uuid.GenerateAvailable(
			self.connection,
			CheckExistsById,
		)
	}

	if nil == err {
		mapped := unmapf(&instance)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(insert_query, mapped)
	}

	return instance, err
}

const update_query string = `
    update instances.instances
    set product_id=:product_id, "name"=:name, description=:description,
        "condition"=:condition, modification_date=:modification_date,
        modification_source=:modification_source
    where id=:id
`

func (self *repository) Update(instance models.Instance) error {
	err := CheckExistsById(self.connection, instance.Id)

	if nil == err {
		err = product.CheckExistsById(self.connection, instance.ProductId)
	}

	if nil == err {
		mapped := unmapf(&instance)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(update_query, mapped)
	}

	return err
}

const get_by_id_query string = `
    select * from instances.instances where id = $1
`

func (self *repository) GetById(instanceId uuid.UUID) (models.Instance, error) {
	var out Instance
	err := CheckExistsById(self.connection, instanceId)

	if nil == err {
		err = self.connection.Get(&out, get_by_id_query, instanceId)
	}

	return mapf(&out), err
}

type query func(*sqlx.DB, uuid.UUID, uint) (*sqlx.Rows, error)

const count_active_query string = `
    select count(*)
    from (
        select *
        from records.pick_up_points_instances
        where out_date is null and instance_id = instances.id
    ) as storages
    join (
        select *
        from records.renters_instances
        where end_date is null and instance_id = instances.id
    ) as provisions
    on storages.instance_id = provisions.instance_id
    where not exists(
        select *
        from deliveries.deliveries
        where actual_end_date is null and instance_id = instances.id
    ) and not exists (
        select * from rents.requests where instance_id = instances.id
    ) and not exists (
        select * from provisions.revokes where instance_id = instances.id
    ) and not exists (
        select *
        from records.users_rents
        where end_date is null and instance_id = instances.Id
    )
`

const base_get_by_category_id_oreder_query string = `
    select id, product_id, "name", description, "condition", modification_date,
           modification_source
    from (
        select id, product_id, "name",
               description, "condition",
               modification_date, modification_source,
               (` + count_active_query + `) as active,
               (%v) as ordering
        from instances.instances
        where product_id = $1
    ) where active = 1
    order by ordering %v
    offset $2
`

var get_by_category_id_order_by_rating_query string = fmt.Sprintf(
	base_get_by_category_id_oreder_query,
	`select avg(rating) from instances.reviews where instance_id = instances.id`,
	"%v",
)

var get_by_category_id_order_by_date_query string = fmt.Sprintf(
	base_get_by_category_id_oreder_query,
	`select start_date
     from records.renters_instances
     where end_date is null and instance_id = instances.id`,
	"%v",
)

var get_by_category_id_order_by_price_query string = fmt.Sprintf(
	base_get_by_category_id_oreder_query,
	`select min(price) from instances.pay_plans where instance_id = instances.id`,
	"%v",
)

var get_by_category_id_order_by_usage_query string = fmt.Sprintf(
	base_get_by_category_id_oreder_query,
	`select count(*)
     from records.users_rents
     where end_date is null and instance_id = instances.id`,
	"%v",
)

func getQueryBySort(sort instance.Sort) query {
	var bquery, direction string
	switch sort {
	case instance.SORT_NONE:
		bquery, direction = get_by_category_id_order_by_date_query, "desc"
	case instance.SORT_RATING_ASC:
		bquery, direction = get_by_category_id_order_by_rating_query, "asc"
	case instance.SORT_RATING_DSC:
		bquery, direction = get_by_category_id_order_by_rating_query, "desc"
	case instance.SORT_DATE_ASC:
		bquery, direction = get_by_category_id_order_by_date_query, "asc"
	case instance.SORT_DATE_DSC:
		bquery, direction = get_by_category_id_order_by_date_query, "desc"
	case instance.SORT_PRICE_ASC:
		bquery, direction = get_by_category_id_order_by_price_query, "asc"
	case instance.SORT_PRICE_DSC:
		bquery, direction = get_by_category_id_order_by_price_query, "desc"
	case instance.SORT_USAGE_ASC:
		bquery, direction = get_by_category_id_order_by_usage_query, "asc"
	case instance.SORT_USAGE_DSC:
		bquery, direction = get_by_category_id_order_by_usage_query, "desc"
	default:
		bquery, direction = get_by_category_id_order_by_date_query, "desc"
	}

	query := fmt.Sprintf(bquery, direction)

	return func(db *sqlx.DB, productId uuid.UUID, offset uint) (*sqlx.Rows, error) {
		return db.Queryx(query, productId, offset)
	}
}

func (self *repository) GetWithFilter(
	filter instance.Filter,
	sort instance.Sort,
) (collection.Collection[models.Instance], error) {
	if err := product.CheckExistsById(self.connection, filter.ProductId); nil != err {
		return nil, err
	}

	query := getQueryBySort(sort)

	return collection.MapCollection(
		mapf,
		sqlCollection.New[Instance](func(offset uint) (*sqlx.Rows, error) {
			return query(self.connection, filter.ProductId, offset)
		}),
	), nil
}

var count_by_id_query string = exist.GenericCounter("instances.instances")

func CheckExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("instance_id", db, count_by_id_query, id)
}

type payPlansRepository struct {
	connection *sqlx.DB
	setter     technical.ISetter
	currency   *currency.Repository
}

func NewPayPlans(
	connection *sqlx.DB,
	setter technical.ISetter,
	currency *currency.Repository,
) instance.IPayPlansRepository {
	return &payPlansRepository{connection, setter, currency}
}

func (self *payPlansRepository) mapPayPlan(
	plan InstancePayPlan,
) models.PayPlan {
	price, _ := self.currency.GetById(plan.CurrencyId)
	price.Value = plan.Price

	return models.PayPlan{
		Id:       plan.Id,
		PeriodId: plan.PeriodId,
		Price:    price,
	}
}

func (self *payPlansRepository) unmapPayPlans(
	value *models.InstancePayPlans,
) ([]InstancePayPlan, error) {
	payPlans := make([]InstancePayPlan, 0, len(value.Map))

	for _, plan := range value.Map {
		mapped, err := self.unmapPayPlan(value.InstanceId, &plan)

		if nil != err {
			return nil, err
		}

		payPlans = append(payPlans, mapped)
	}

	return payPlans, nil
}

func (self *payPlansRepository) unmapPayPlan(
	instanceId uuid.UUID,
	value *models.PayPlan,
) (InstancePayPlan, error) {
	cid, err := self.currency.GetId(value.Price.Name)

	if nil != err {
		return InstancePayPlan{}, err
	} else {
		return InstancePayPlan{
			Id:         value.Id,
			InstanceId: instanceId,
			PeriodId:   value.PeriodId,
			CurrencyId: cid,
			Price:      value.Price.Value,
		}, nil
	}
}

const insert_pay_plan_query string = `
    insert into instances.pay_plans (
        id, instance_id, period_id, currency_id, price, modification_date,
        modification_source
    ) values (
        :id, :instance_id, :period_id, :currency_id, :price, :modification_date,
        :modification_source
    )
`

func (self *payPlansRepository) Create(
	payPlans models.InstancePayPlans,
) (models.InstancePayPlans, error) {
	err := CheckExistsById(self.connection, payPlans.InstanceId)

	if nil == err {
		err = CheckPayPlanExistsByInstnaceId(self.connection, payPlans.InstanceId)

		if nil == err {
			err = cmnerrors.Duplicate("instance_pay_plans")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	if nil == err {
		for k := range payPlans.Map {
			if nil == err {
				err = period.CheckExistsById(self.connection, k)
			}
		}
	}

	if nil == err {
		for k, v := range payPlans.Map {
			if nil == err {
				v.Id, err = gen_uuid.GenerateAvailable(
					self.connection,
					CheckPayPlanExistsById,
				)
			}

			if nil == err {
				payPlans.Map[k] = v
			}
		}
	}

	var mpp []InstancePayPlan
	if nil == err {
		mpp, err = self.unmapPayPlans(&payPlans)
	}

	for i := 0; nil == err && len(mpp) > i; i++ {
		self.setter.Update(&mpp[i].Info)
		_, err = self.connection.NamedExec(insert_pay_plan_query, mpp[i])
	}

	return payPlans, err
}

func (self *payPlansRepository) AddPayPlan(
	instanceId uuid.UUID,
	plan models.PayPlan,
) (models.InstancePayPlans, error) {
	err := CheckExistsById(self.connection, instanceId)

	if nil == err {
		err = period.CheckExistsById(self.connection, plan.PeriodId)
	}

	if nil == err {
		err = CheckPayPlanExistsByInstanceIdAndPeriodId(
			self.connection,
			instanceId, plan.PeriodId,
		)

		if nil == err {
			err = cmnerrors.Duplicate("instance_pay_plan_instance")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	var rows *sqlx.Rows
	var out models.InstancePayPlans
	var mapped InstancePayPlan

	if nil == err {
		mapped, err = self.unmapPayPlan(instanceId, &plan)
	}

	if nil == err {
		mapped.Id, err = gen_uuid.GenerateAvailable(
			self.connection,
			CheckPayPlanExistsById,
		)
	}

	if nil == err {
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(insert_pay_plan_query, mapped)
	}

	if nil == err {
		rows, err = self.getByInstanceId(instanceId)
	}

	if nil == err {
		out, err = self.fold(instanceId, rows)
	}

	return out, err
}

const update_pay_plan_query string = `
    update instances.pay_plans
    set instance_id=:instance_id, period_id=:period_id,
        currency_id=:currency_id, price=:price,
        modification_date=:modification_date,
        modification_source=:modification_source
    where id=:id
`

const delete_pay_plan_query string = `
    delete from instances.pay_plans where id = $1
`

func (self *payPlansRepository) Update(payPlans models.InstancePayPlans) error {
	var rows *sqlx.Rows
	var plan, mapped InstancePayPlan
	err := CheckExistsById(self.connection, payPlans.InstanceId)

	if nil == err {
		rows, err = self.getByInstanceId(payPlans.InstanceId)
	}

	for nil == err && rows.Next() {
		if err = rows.StructScan(&plan); nil == err {
			if uplan, found := payPlans.Map[plan.PeriodId]; !found {
				_, err = self.connection.Exec(
					delete_pay_plan_query, plan.Id,
				)
			} else {
				delete(payPlans.Map, plan.PeriodId)
				mapped, err = self.unmapPayPlan(payPlans.InstanceId, &uplan)

				if nil == err {
					self.setter.Update(&mapped.Info)
					_, err = self.connection.NamedExec(
						update_pay_plan_query, mapped,
					)
				}
			}
		}
	}

	if nil == err && 0 != len(payPlans.Map) {
		unknown := make([]uuid.UUID, 0, len(payPlans.Map))

		for _, plan := range payPlans.Map {
			unknown = append(unknown, plan.PeriodId)
		}

		err = fmt.Errorf("Found unknown pay periods: %v", unknown)
	}

	return err
}

const get_pay_plan_by_instance_id_query string = `
    select * from instances.pay_plans where instance_id = $1
`

func (self *payPlansRepository) GetByInstanceId(
	instanceId uuid.UUID,
) (models.InstancePayPlans, error) {
	var out models.InstancePayPlans
	rows, err := self.getByInstanceId(instanceId)

	if nil == err {
		out, err = self.fold(instanceId, rows)
	}

	return out, err
}

func (self *payPlansRepository) fold(instanceId uuid.UUID, rows *sqlx.Rows) (models.InstancePayPlans, error) {
	var plan InstancePayPlan
	var err error
	var out = models.InstancePayPlans{
		InstanceId: instanceId,
		Map:        make(map[uuid.UUID]models.PayPlan),
	}

	for nil == err && rows.Next() {
		if err = rows.StructScan(&plan); nil == err {
			mapped := self.mapPayPlan(plan)
			out.Map[mapped.PeriodId] = mapped
		}
	}

	if nil != err {
		out = models.InstancePayPlans{}
	}

	return out, err
}

func (self *payPlansRepository) getByInstanceId(
	instanceId uuid.UUID,
) (*sqlx.Rows, error) {
	var rows *sqlx.Rows
	err := CheckExistsById(self.connection, instanceId)

	if nil == err {
		err = CheckPayPlanExistsByInstnaceId(self.connection, instanceId)
	}

	if nil == err {
		rows, err = self.connection.Queryx(
			get_pay_plan_by_instance_id_query,
			instanceId,
		)
	}

	return rows, err
}

var count_pay_plan_by_id_query string = exist.GenericCounter("instances.pay_plans")

func CheckPayPlanExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check(
		"instance_pay_plan_id",
		db,
		count_pay_plan_by_id_query,
		id,
	)
}

var count_pay_plan_by_instance_id_query string = `
    select count(*) from instances.pay_plans where instance_id = $1
`

func CheckPayPlanExistsByInstnaceId(db *sqlx.DB, id uuid.UUID) error {
	return exist.CheckMultiple(
		"instance_pay_plan_instance_id",
		db,
		count_pay_plan_by_instance_id_query,
		id,
	)
}

var count_pay_plan_by_period_id_query string = `
    select count(*) from instances.pay_plans where period_id = $1
`

func CheckPayPlanExistsByPeriodId(db *sqlx.DB, id uuid.UUID) error {
	return exist.CheckMultiple(
		"instance_pay_plan_period_id",
		db,
		count_pay_plan_by_period_id_query,
		id,
	)
}

const count_pay_plan_by_instance_id_and_period_id_query string = `
    select count(*)
    from instances.pay_plans
    where instance_id = $1 and period_id = $2
`

func CheckPayPlanExistsByInstanceIdAndPeriodId(
	db *sqlx.DB,
	instanceId uuid.UUID,
	periodId uuid.UUID,
) error {
	return exist.Check(
		"instance_pay_plan_instance",
		db,
		count_pay_plan_by_instance_id_and_period_id_query,
		instanceId, periodId,
	)
}

type photoRepository struct {
	connection *sqlx.DB
	setter     technical.ISetter
}

func NewPhoto(
	connection *sqlx.DB,
	setter technical.ISetter,
) instance.IPhotoRepository {
	return &photoRepository{connection, setter}
}

const insert_photo_query string = `
    insert into instances.photos (
        id, instance_id, photo_id, modification_date, modification_source
    ) values (
        $1, $2, $3, $4, $5
    )
`

func (self *photoRepository) Create(
	instanceId uuid.UUID,
	photoId uuid.UUID,
) error {
	var id uuid.UUID
	err := CheckExistsById(self.connection, instanceId)

	if nil == err {
		err = photo.CheckExistsById(self.connection, photoId)
	}

	if nil == err {
		id, err = gen_uuid.GenerateAvailable(
			self.connection,
			CheckPhotoExistsById,
		)
	}

	if nil == err {
		info := technical.Info{}
		self.setter.Update(&info)
		_, err = self.connection.Exec(
			insert_photo_query,
			id, instanceId, photoId, info.MDate, info.MSource,
		)
	}

	return err
}

const get_photo_by_id_query string = `
    select photo_id
    from instances.photos
    where instance_id = $1
    offset $2
`

func (self *photoRepository) GetByInstanceId(
	instanceId uuid.UUID,
) (collection.Collection[uuid.UUID], error) {
	if err := CheckExistsById(self.connection, instanceId); nil != err {
		return nil, err
	}

	return sqlCollection.NewDirect[uuid.UUID](func(offset uint) (*sqlx.Rows, error) {
		return self.connection.Queryx(get_photo_by_id_query, instanceId, offset)
	}), nil
}

var count_photo_by_id_query string = exist.GenericCounter("instances.photos")

func CheckPhotoExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("instance_photo_id", db, count_photo_by_id_query, id)
}

