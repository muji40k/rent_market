package psql

import (
	"errors"
	"fmt"
	"rent_service/internal/factory/repositories/v1/psql"
	"rent_service/internal/repository/context/v1"
	"rent_service/internal/repository/implementation/sql/repositories/user"
	"rent_service/internal/repository/implementation/sql/technical"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Builder struct {
	connection struct {
		user     string
		password string
		host     string
		port     string
		database string
	}
	setter technical.ISetter
	hasher user.Hasher
}

func New() *Builder {
	return &Builder{}
}

func (self *Builder) WithUser(user string) *Builder {
	self.connection.user = user
	return self
}

func (self *Builder) WithPassword(password string) *Builder {
	self.connection.password = password
	return self
}

func (self *Builder) WithHost(host string) *Builder {
	self.connection.host = host
	return self
}

func (self *Builder) WithPort(port string) *Builder {
	self.connection.port = port
	return self
}

func (self *Builder) WithDatabase(database string) *Builder {
	self.connection.database = database
	return self
}

func (self *Builder) WithSetter(setter technical.ISetter) *Builder {
	self.setter = setter
	return self
}

func (self *Builder) WithHasher(hasher user.Hasher) *Builder {
	self.hasher = hasher
	return self
}

func (self *Builder) getConnString() string {
	return fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v",
		self.connection.user,
		self.connection.password,
		self.connection.host,
		self.connection.port,
		self.connection.database,
	)
}

func (self *Builder) Build() (v1.Factories, error) {
	var factory *psql.Factory
	db, err := sqlx.Connect("pgx", self.getConnString())

	if nil != err {
		err = fmt.Errorf("PSQLFactoryBuilder: %w", err)
	}

	if nil == err && nil == self.setter {
		err = errors.New("PSQLFactoryBuilder: Setter not set")
	}

	if nil == err && nil == self.hasher {
		err = errors.New("PSQLFactoryBuilder: Hasher not set")
	}

	if nil == err {
		factory = psql.New(db, self.setter, self.hasher)
	}

	return factory.ToFactories(), err
}

