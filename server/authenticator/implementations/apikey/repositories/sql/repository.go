package sql

import (
	"database/sql"
	"errors"
	"rent_service/internal/logic/services/types/token"
	"rent_service/server/authenticator"
	"rent_service/server/authenticator/implementations/apikey"
	"time"

	"github.com/jmoiron/sqlx"
)

type session struct {
	AccessToken   string    `db:"access_token"`
	AccessValidTo time.Time `db:"access_valid_to"`
	RenewToken    string    `db:"renew_token"`
	RenewValidTo  time.Time `db:"renew_valid_to"`
	Token         string    `db:"token"`
}

type repository struct {
	connection *sqlx.DB
}

func New(connection *sqlx.DB) apikey.ITokenRepository {
	return &repository{connection}
}

const delete_by_access_token_query string = `
    delete from public.sessions where access_token = $1
`

func (self *repository) DeleteToken(access string) error {
	_, err := self.connection.Exec(delete_by_access_token_query, access)

	if nil != err {
		err = apikey.ErrorDataAccess
	}

	return err
}

const get_by_access_token_query string = `
    select * from public.sessions where access_token = $1
`

func (self *repository) GetToken(access string) (token.Token, error) {
	var s session
	err := self.checkExistsByAccess(access)

	if nil == err {
		err = self.connection.Get(&s, get_by_access_token_query, access)

		if nil != err {
			err = apikey.ErrorDataAccess
		}
	}

	if nil == err && time.Now().After(s.AccessValidTo) {
		err = apikey.ErrorNotFound
	}

	if nil == err {
		return token.Token(s.Token), nil
	} else {
		return token.Token(""), err
	}
}

const get_by_access_and_renew_token_query string = `
    select * from public.sessions where access_token = $1 and renew_token = $2
`

func (self *repository) RenewToken(apiToken authenticator.ApiToken) (token.Token, error) {
	var s session
	err := self.checkExistsByAccess(apiToken.Access)

	if nil == err {
		err = self.connection.Get(
			&s,
			get_by_access_and_renew_token_query,
			apiToken.Access,
			apiToken.Renew,
		)

		self.connection.Exec(delete_by_access_token_query, apiToken.Access)

		if errors.Is(err, sql.ErrNoRows) {
			err = apikey.ErrorNotFound
		} else if nil != err {
			err = apikey.ErrorDataAccess
		}
	}

	if nil == err && time.Now().After(s.RenewValidTo) {
		err = apikey.ErrorNotFound
	}

	if nil == err {
		return token.Token(s.Token), nil
	} else {
		return token.Token(""), err
	}
}

const insert_query string = `
    insert into public.sessions (
        access_token, access_valid_to, renew_token, renew_valid_to, token
    ) values (
        :access_token, :access_valid_to, :renew_token, :renew_valid_to, :token
    )
`

func (self *repository) WriteToken(
	token token.Token,
	access apikey.TokenHandle,
	renew apikey.TokenHandle,
) error {
	if nil == self.checkExistsByAccess(access.Value) {
		return apikey.ErrorDataAccess
	}

	var s = session{
		AccessToken:   access.Value,
		AccessValidTo: access.ValidTo,
		RenewToken:    renew.Value,
		RenewValidTo:  renew.ValidTo,
		Token:         string(token),
	}

	_, err := self.connection.NamedExec(insert_query, s)

	if nil != err {
		err = apikey.ErrorDataAccess
	}

	return err
}

const count_by_access_token_query string = `
    select count(*) from public.sessions where access_token = $1
`

func (self *repository) checkExistsByAccess(access string) error {
	var count uint

	err := self.connection.Get(&count, count_by_access_token_query, access)

	if nil == err && 0 == count {
		err = apikey.ErrorNotFound
	} else if (nil == err && 1 != count) || nil != err {
		err = apikey.ErrorDataAccess
	}

	return err
}

func (self *repository) Clear() {
	if nil != self.connection {
		self.connection.Close()
	}
}

