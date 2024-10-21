package paymethod

import (
	"database/sql"
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	sqlCollection "rent_service/internal/repository/implementation/sql/collection"
	"rent_service/internal/repository/implementation/sql/exist"
	"rent_service/internal/repository/implementation/sql/technical"
	"rent_service/internal/repository/interfaces/paymethod"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PayMethod struct {
	Id          uuid.UUID      `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	technical.Info
}

type repository struct {
	connection *sqlx.DB
}

func New(connection *sqlx.DB) paymethod.IRepository {
	return &repository{connection}
}

func mapPayMethod(value *PayMethod) models.PayMethod {
	out := models.PayMethod{
		Id:   value.Id,
		Name: value.Name,
	}

	if value.Description.Valid {
		out.Description = value.Description.String
	}

	return out
}

const get_all_query string = "select * from payments.methods offset $1"

func (self *repository) GetAll() (collection.Collection[models.PayMethod], error) {
	return collection.MapCollection(
		mapPayMethod,
		sqlCollection.New[PayMethod](func(offset uint) (*sqlx.Rows, error) {
			return self.connection.Queryx(get_all_query, offset)
		}),
	), nil
}

var count_by_id_query = exist.GenericCounter("payments.methods")

func CheckExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("pay_method_id", db, count_by_id_query, id)
}

