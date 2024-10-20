package payment

import (
	"database/sql"
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	sqlCollection "rent_service/internal/repository/implementation/sql/collection"
	"rent_service/internal/repository/implementation/sql/repositories/currency"
	"rent_service/internal/repository/implementation/sql/repositories/instance"
	"rent_service/internal/repository/implementation/sql/repositories/rent"
	"rent_service/internal/repository/implementation/sql/technical"
	"rent_service/internal/repository/interfaces/payment"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Payment struct {
	Id          uuid.UUID      `db:"id"`
	RentId      uuid.UUID      `db:"rent_id"`
	PayMethodId uuid.NullUUID  `db:"pay_method_id"`
	PaymentId   sql.NullString `db:"payment_id"`
	PeriodStart time.Time      `db:"period_strat"`
	PeriodEnd   time.Time      `db:"period_end"`
	CurrencyId  uuid.UUID      `db:"currncy_id"`
	Value       float64        `db:"value"`
	Status      sql.NullString `db:"status"`
	CreateDate  time.Time      `db:"create_date"`
	PaymentDate sql.NullTime   `db:"payment_date"`
	technical.Info
}

type repository struct {
	connection *sqlx.DB
	currency   *currency.Repository
}

func New(
	connection *sqlx.DB,
	currency *currency.Repository,
) payment.IRepository {
	return &repository{connection, currency}
}

func (self *repository) mapPayment(value *Payment) models.Payment {
	out := models.Payment{
		Id:          value.Id,
		RentId:      value.RentId,
		PeriodStart: value.PeriodStart,
		PeriodEnd:   value.PeriodEnd,
		CreateDate:  value.CreateDate,
	}

	if value.Status.Valid {
		out.Status = value.Status.String
	}

	if value.PayMethodId.Valid {
		out.PayMethodId = new(uuid.UUID)
		*out.PayMethodId = value.PayMethodId.UUID
	}

	if value.PaymentId.Valid {
		out.PaymentId = new(string)
		*out.PaymentId = value.PaymentId.String
	}

	if value.PaymentDate.Valid {
		out.PaymentDate = new(time.Time)
		*out.PaymentDate = value.PaymentDate.Time
	}

	out.Value, _ = self.currency.GetById(value.CurrencyId)
	out.Value.Value = value.Value

	return out
}

const get_by_instance_id_query string = `
    select payments.id, payments.rent_id, payments.pay_method_id,
           payments.payment_id, payments.period_strat, payments.period_end,
           payments.currency_id, payments.value, payments.status,
           payments.create_date, payments.payment_date,
           payments.modification_date, payments.modification_source
    from payments.payments
    join (select * from records.users_rents where instance_id = $1) as instance_rents
        on payments.payments.rent_id = instance_rents.id
    offset $2
`

func (self *repository) GetByInstanceId(
	instanceId uuid.UUID,
) (collection.Collection[models.Payment], error) {
	if err := instance.CheckExistsById(self.connection, instanceId); nil != err {
		return nil, err
	}

	return collection.MapCollection(
		self.mapPayment,
		sqlCollection.New[Payment](func(offset uint) (*sqlx.Rows, error) {
			return self.connection.Queryx(
				get_by_instance_id_query,
				instanceId,
				offset,
			)
		}),
	), nil
}

const get_by_rent_id string = `
    select * from payments.payments where ren_id = $1 offset $2
`

func (self *repository) GetByRentId(
	rentId uuid.UUID,
) (collection.Collection[models.Payment], error) {
	if err := rent.CheckExistsById(self.connection, rentId); nil != err {
		return nil, err
	}

	return collection.MapCollection(
		self.mapPayment,
		sqlCollection.New[Payment](func(offset uint) (*sqlx.Rows, error) {
			return self.connection.Queryx(
				get_by_rent_id,
				rentId,
				offset,
			)
		}),
	), nil
}

