package server

import (
	"fmt"
	"os"
	"rent_service/builders/application/web/authenticator/apikey"
	"rent_service/builders/application/web/authenticator/apikey/repository/psql"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/context/v1"
	"rent_service/misc/contextholder"
	"rent_service/server"
	"rent_service/server/authenticator"
	"rent_service/server/controllers/categories"
	"rent_service/server/controllers/deliveries"
	"rent_service/server/controllers/deliverycompanies"
	"rent_service/server/controllers/instances"
	"rent_service/server/controllers/payments"
	"rent_service/server/controllers/paymethods"
	"rent_service/server/controllers/periods"
	"rent_service/server/controllers/photos"
	"rent_service/server/controllers/pickuppoints"
	"rent_service/server/controllers/products"
	"rent_service/server/controllers/profiles"
	prequests "rent_service/server/controllers/provision/requests"
	preturns "rent_service/server/controllers/provision/returns"
	"rent_service/server/controllers/provisions"
	rrequests "rent_service/server/controllers/rent/requests"
	rreturns "rent_service/server/controllers/rent/returns"
	"rent_service/server/controllers/rents"
	"rent_service/server/controllers/roles"
	"rent_service/server/controllers/sessions"
	"rent_service/server/controllers/storedinstances"
	"rent_service/server/controllers/users"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

const (
	ACCESS_TIME time.Duration = 24 * time.Hour
	RENEW_TIME  time.Duration = 7 * 24 * time.Hour
)

const (
	TEST_HOST     string = "TEST_DB_HOST"
	TEST_PORT     string = "TEST_DB_PORT"
	TEST_DATABASE string = "TEST_DB_AUTH_NAME"
	TEST_USER     string = "TEST_DB_USERNAME"
	TEST_PASSWORD string = "TEST_DB_PASSWORD"
)

func getOr(variable string, def string) string {
	if value := os.Getenv(variable); "" != value {
		return value
	} else {
		return def
	}
}

type dbConfig struct {
	host     string
	port     string
	database string
	user     string
	password string
}

func parse() dbConfig {
	return dbConfig{
		host:     getOr(TEST_HOST, "localhost"),
		port:     getOr(TEST_PORT, "5432"),
		database: getOr(TEST_DATABASE, "authentication"),
		user:     getOr(TEST_USER, "postgres"),
		password: getOr(TEST_PASSWORD, "postgres"),
	}
}

type ServerExtender interface{}

type ControllerCreator func(authenticator.IAuthenticator, *v1.Context) server.IController
type EngindeExtender func(*gin.Engine)

func TestServer(
	factories v1.Factories,
	extenders ...ServerExtender,
) *gin.Engine {
	ctx := v1.New(factories)
	config := parse()
	repo, err := psql.New().
		WithHost(config.host).
		WithPort(config.port).
		WithDatabase(config.database).
		WithUser(config.user).
		WithPassword(config.password).
		Build()

	if nil != err {
		panic(err)
	}

	a, err := apikey.New().
		WithAccesTime(ACCESS_TIME).
		WithRenewTime(RENEW_TIME).
		WithLogin(ctx).
		WithTokenRepository(repo).
		Build()

	if nil != err {
		panic(err)
	}

	engine := gin.New()

	for _, ext := range extenders {
		switch ext := ext.(type) {
		case ControllerCreator:
			ext(a, ctx).Register(engine)
		case EngindeExtender:
			ext(engine)
		default:
			panic("Unknown extender")
		}
	}

	return engine
}

type Inserter struct {
	db *sqlx.DB
}

func (self *dbConfig) getConnection() (*sqlx.DB, error) {
	return sqlx.Connect("pgx",
		fmt.Sprintf(
			"postgres://%v:%v@%v:%v/%v",
			self.user,
			self.password,
			self.host,
			self.port,
			self.database,
		),
	)
}

func NewInserter() *Inserter {
	config := parse()
	db, err := config.getConnection()

	if nil != err {
		panic(err)
	}

	return &Inserter{db}
}

func (self *Inserter) Close() {
	self.db.Close()
}

func (self *Inserter) ClearDB() {
	_, err := self.db.Exec("delete from public.sessions")

	if nil != err {
		panic(err)
	}
}

func (self *Inserter) InsertSession(user models.User, access string, renew string) {
	now := time.Now()
	_, err := self.db.Exec(`
        INSERT INTO public.sessions (
            access_token, access_valid_to, renew_token, renew_valid_to, token
        ) VALUES (
            $1, $2, $3, $4, $5
        )`,
		access, now.Add(ACCESS_TIME), renew, now.Add(RENEW_TIME),
		string(user.Token),
	)

	if nil != err {
		panic(err)
	}
}

func TracerExtender(tr trace.TracerProvider, hl *contextholder.Holder) ServerExtender {
	return EngindeExtender(server.TracerExtender(tr, hl))
}

func DefaultControllers() []ServerExtender {
	return []ServerExtender{
		ControllerCreator(Categories), ControllerCreator(Deliveries),
		ControllerCreator(DeliveryCompanies), ControllerCreator(Instances),
		ControllerCreator(Payments), ControllerCreator(PayMethods),
		ControllerCreator(Periods), ControllerCreator(Photos),
		ControllerCreator(PickUpPoints), ControllerCreator(Products),
		ControllerCreator(Profiles), ControllerCreator(ProvisionRequests),
		ControllerCreator(ProvisionReturns), ControllerCreator(Provisions),
		ControllerCreator(RentRequests), ControllerCreator(RentReturns),
		ControllerCreator(Rents), ControllerCreator(Roles),
		ControllerCreator(Sessions), ControllerCreator(StoredInstances),
		ControllerCreator(Users),
	}
}

func Categories(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return categories.New(c)
}

func Deliveries(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return deliveries.New(a, c)
}
func DeliveryCompanies(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return deliverycompanies.New(c, a)
}

func Instances(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return instances.New(a, c, c, c, c)
}

func Payments(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return payments.New(a, c)
}

func PayMethods(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return paymethods.New(c)
}

func Periods(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return periods.New(c)
}

func Photos(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return photos.New(c, a)
}

func PickUpPoints(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return pickuppoints.New(c, c, c)
}

func Products(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return products.New(c, c, c)
}

func Profiles(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return profiles.New(a, c, c, c)
}

func ProvisionRequests(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return prequests.New(a, c, c)
}

func ProvisionReturns(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return preturns.New(a, c, c)
}

func Provisions(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return provisions.New(a, c)
}

func RentRequests(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return rrequests.New(a, c, c)
}

func RentReturns(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return rreturns.New(a, c, c)
}

func Rents(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return rents.New(a, c)
}

func Roles(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return roles.New(a, c)
}

func Sessions(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return sessions.New(a)
}

func StoredInstances(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return storedinstances.New(a, c)
}

func Users(a authenticator.IAuthenticator, c *v1.Context) server.IController {
	return users.New(a, c, c)
}

