package psql

import (
	"fmt"
	"rent_service/server/authenticator/implementations/apikey"
	"rent_service/server/authenticator/implementations/apikey/repositories/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Builder struct {
	user     string
	password string
	host     string
	port     string
	database string
}

func (self *Builder) WithUser(user string) *Builder {
	self.user = user
	return self
}

func (self *Builder) WithPassword(password string) *Builder {
	self.password = password
	return self
}

func (self *Builder) WithHost(host string) *Builder {
	self.host = host
	return self
}

func (self *Builder) WithPort(port string) *Builder {
	self.port = port
	return self
}

func (self *Builder) WithDatabase(database string) *Builder {
	self.database = database
	return self
}

func (self *Builder) getConnString() string {
	return fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v",
		self.user,
		self.password,
		self.host,
		self.port,
		self.database,
	)
}

func (self *Builder) Build() (apikey.ITokenRepository, error) {
	var repo apikey.ITokenRepository
	db, err := sqlx.Connect("pgx", self.getConnString())

	if nil == err {
		repo = sql.New(db)
	}

	return repo, err
}

