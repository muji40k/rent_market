package period

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	sqlCollection "rent_service/internal/repository/implementation/sql/collection"
	"rent_service/internal/repository/implementation/sql/exist"
	"rent_service/internal/repository/implementation/sql/technical"
	"rent_service/internal/repository/interfaces/period"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Period struct {
	Id       uuid.UUID     `db:"id"`
	Name     string        `db:"name"`
	Duration time.Duration `db:"duration"`
	technical.Info
}

type repository struct {
	connection *sqlx.DB
}

func New(connection *sqlx.DB) period.IRepository {
	return &repository{connection}
}

func mapPeriod(value *Period) models.Period {
	return models.Period{
		Id:       value.Id,
		Name:     value.Name,
		Duration: value.Duration,
	}
}

const get_by_id_query string = `
    select * from periods.periods where id = $1
`

func (self *repository) GetById(periodId uuid.UUID) (models.Period, error) {
	var out Period
	err := CheckExistsById(self.connection, periodId)

	if nil == err {
		err = self.connection.Get(&out, get_by_id_query, periodId)
	}

	return mapPeriod(&out), err
}

const get_all_query string = "select * from periods.periods offset $1"

func (self *repository) GetAll() (collection.Collection[models.Period], error) {
	return collection.MapCollection(
		mapPeriod,
		sqlCollection.New[Period](func(offset uint) (*sqlx.Rows, error) {
			return self.connection.Queryx(get_all_query, offset)
		}),
	), nil
}

var count_by_id_query string = exist.GenericCounter("periods.periods")

func CheckExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("period_id", db, count_by_id_query, id)
}

