package address

import (
	"database/sql"
	"rent_service/internal/domain/models"
	"rent_service/internal/repository/implementation/sql/exist"
	"rent_service/internal/repository/implementation/sql/technical"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Address struct {
	Id      uuid.UUID      `db:"id"`
	Country string         `db:"country"`
	City    string         `db:"city"`
	Street  string         `db:"street"`
	House   string         `db:"house"`
	Flat    sql.NullString `db:"flat"`
	technical.Info
}

type Repository struct {
	connection *sqlx.DB
}

func New(connection *sqlx.DB) *Repository {
	return &Repository{connection}
}

func mapf(value *Address) models.Address {
	out := models.Address{
		Id:      value.Id,
		Country: value.Country,
		City:    value.City,
		Street:  value.Street,
		House:   value.House,
	}

	if value.Flat.Valid {
		out.Flat = new(string)
		*out.Flat = value.Flat.String
	}

	return out
}

var get_by_id_query string = `
    select name from addresses.addresses where id = $1
`

func (self *Repository) GetById(
	addressId uuid.UUID,
) (models.Address, error) {
	var out Address
	err := CheckExistsById(self.connection, addressId)

	if nil == err {
		err = self.connection.Get(&out, get_by_id_query, addressId)
	}

	return mapf(&out), err
}

var count_by_id_query string = exist.GenericCounter("addresses.addresses")

func CheckExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("address_id", db, count_by_id_query, id)
}

