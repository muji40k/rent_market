package currency

import (
	"rent_service/internal/misc/types/currency"
	"rent_service/internal/repository/implementation/sql/exist"
	"rent_service/misc/mapfuncs"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const VALID_PERIOD time.Duration = 7 * 24 * time.Hour

type cell struct {
	currency currency.Currency
	validTo  time.Time
}

type Repository struct {
	connection *sqlx.DB
	cache      map[uuid.UUID]cell
}

func New(connection *sqlx.DB) *Repository {
	return &Repository{connection, make(map[uuid.UUID]cell)}
}

var get_by_id_query string = `
    select name from currencies.currencies where id = $1
`

func (self *Repository) GetById(
	currencyId uuid.UUID,
) (currency.Currency, error) {
	now := time.Now()

	if c, found := self.cache[currencyId]; found && now.After(c.validTo) {
		return c.currency, nil
	}

	var out = currency.Currency{Name: "unknown", Value: 0}
	err := CheckExistsById(self.connection, currencyId)

	if nil == err {
		err = self.connection.Get(&out.Name, get_by_id_query, currencyId)
	}

	if nil == err {
		self.cache[currencyId] = cell{out, now.Add(VALID_PERIOD)}
	} else {
		delete(self.cache, currencyId)
	}

	return out, err
}

const get_by_name_query string = `
    select id from currencies.currencies where name = $1
`

func (self *Repository) GetId(currencyName string) (uuid.UUID, error) {
	now := time.Now()

	id, found := mapfuncs.FindByValueF(self.cache, func(value *cell) bool {
		return value.currency.Name == currencyName
	})

	if found && now.After(self.cache[id].validTo) {
		return id, nil
	}

	err := CheckExistsByName(self.connection, currencyName)

	if nil == err {
		err = self.connection.Get(&id, get_by_name_query, currencyName)
	}

	return id, err
}

func (self *Repository) CleanCache() {
	self.cache = make(map[uuid.UUID]cell)
}

var count_by_id_query string = exist.GenericCounter("currencies.currencies")

func CheckExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("currency_id", db, count_by_id_query, id)
}

const count_by_name_query string = `
    select count(*) from currencies.currencies where name = $1
`

func CheckExistsByName(db *sqlx.DB, name string) error {
	return exist.Check("currency_name", db, count_by_name_query, name)
}

