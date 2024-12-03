package psql

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"os"
	"rent_service/builders/misc/generator"
	"rent_service/builders/misc/uuidgen"
	"rent_service/builders/repository/factory/v1/psql"
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	"rent_service/internal/repository/implementation/sql/technical/implementations/simple"
	"rent_service/misc/nullable"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	TEST_HOST     string = "TEST_DB_HOST"
	TEST_PORT     string = "TEST_DB_PORT"
	TEST_DATABASE string = "TEST_DB_NAME"
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
		database: getOr(TEST_DATABASE, "rent_market"),
		user:     getOr(TEST_USER, "postgres"),
		password: getOr(TEST_PASSWORD, "postgres"),
	}
}

func PSQLRepositoryFactory() *psql.Builder {
	conf := parse()

	return psql.New().
		WithHost(conf.host).
		WithPort(conf.port).
		WithDatabase(conf.database).
		WithUser(conf.user).
		WithPassword(conf.password).
		WithHasher(stubHasher).
		WithSetter(simple.New("test"))
}

func stubHasher(user *models.User) string {
	return user.Email + user.Name + user.Password
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

func GeneratorStepNewList[T any](
	t generator.IAllureProvider,
	name string,
	base func(uint) (T, uuid.UUID),
	inserter func(*T),
) generator.IGenerator {
	out := make([]T, 0)
	return GeneratorStepList(t, name, &out, base, inserter)
}

func GeneratorStepList[T any](
	t generator.IAllureProvider,
	name string,
	reference *[]T,
	base func(uint) (T, uuid.UUID),
	inserter func(*T),
) generator.IGenerator {
	i := uint(0)
	spy := nullable.Some(generator.Spy{})

	return generator.NewAllureWrap(
		t, fmt.Sprintf("Create and insert %v", name),
		generator.NewFuncWrapped(
			generator.FuncListWrap(
				reference,
				func() (T, uuid.UUID) {
					j := i
					i++
					return base(j)
				},
				func(items []T) {
					nullable.IfSome(spy, func(spy *generator.Spy) {
						spy.SniffValue(name, items)
					})
					BulkInsert(inserter, items...)
				},
			),
		),
		spy,
	)
}

func GeneratorStepValue[T any](
	t generator.IAllureProvider,
	name string,
	reference *T,
	base func() (T, uuid.UUID),
	inserter func(*T),
) generator.IGenerator {
	id := nullable.None[uuid.UUID]()
	spy := nullable.Some(generator.Spy{})

	return generator.NewAllureWrap(
		t, fmt.Sprintf("Create and insert %v", name),
		generator.NewFunc(
			func() uuid.UUID {
				return nullable.GetOrInsertFunc(id, func() uuid.UUID {
					var id uuid.UUID
					*reference, id = base()
					return id
				})
			},
			func() {
				nullable.IfSome(spy, func(spy *generator.Spy) {
					spy.SniffValue(name, *reference)
				})
				inserter(reference)
			},
		),
		spy,
	)
}

func GeneratorStepNewValue[T any](
	t generator.IAllureProvider,
	name string,
	base func() (T, uuid.UUID),
	inserter func(*T),
) generator.IGenerator {
	return GeneratorStepValue(t, name, new(T), base, inserter)
}

func GeneratorStep(
	t generator.IAllureProvider,
	name string,
	gen func(spy *nullable.Nullable[generator.Spy]) generator.IGenerator,
) generator.IGenerator {
	spy := nullable.Some(generator.Spy{})

	return generator.NewAllureWrap(
		t, fmt.Sprintf("Create and insert %v", name),
		gen(spy), spy,
	)
}

func prepareInsert(table string, columnNames ...string) string {
	lim := len(columnNames)
	var columns = make([]string, lim+2)
	var idx = make([]string, lim+2)

	for i, v := range columnNames {
		columns[i] = v
		idx[i] = fmt.Sprintf("$%v", i+1)
	}

	columns[lim] = "modification_date"
	idx[lim] = "now()"
	columns[lim+1] = "modification_source"
	idx[lim+1] = "'test_preset'"

	return fmt.Sprintf(
		"insert into %v (%v) values (%v)",
		table, strings.Join(columns, ","), strings.Join(idx, ","),
	)
}

type Inserter struct {
	db *sqlx.DB
}

func NewInserter() *Inserter {
	config := parse()
	db, err := config.getConnection()

	if nil != err {
		panic(err)
	}

	i := &Inserter{db}
	i.initDB()

	return i
}

func (self *Inserter) initDB() {
	var init bool
	_, err := self.db.Exec("create table if not exists public.initter(initialised boolean)")

	if nil == err {
		row := self.db.QueryRow("select * from public.initter")
		err = row.Scan(&init)

		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
	}

	if nil == err && !init {
		self.ClearDB()
		_, err = self.db.Exec("insert into public.initter(initialised) values (true)")
	}

	if nil != err {
		panic(err)
	}
}

func (self *Inserter) Close() {
	self.db.Close()
}

var currencies = map[string]uuid.UUID{
	"rub":     uuidgen.Generate(),
	"usd":     uuidgen.Generate(),
	"unknown": uuidgen.Generate(),
}

var tables = []string{
	"payments.payments",

	"records.renters_instances",
	"records.users_rents",
	"records.pick_up_points_instances",

	"deliveries.deliveries",

	"provisions.requests_pay_plans",
	"provisions.requests",
	"provisions.revokes",

	"rents.requests",
	"rents.returns",

	"roles.renters",
	"roles.administrators",
	"roles.storekeepers",

	"instances.reviews",

	"payments.users_methods",
	"users.favorite_pick_up_points",
	"users.profiles",
	"users.users",

	"pick_up_points.photos",
	"pick_up_points.working_hours",
	"pick_up_points.pick_up_points",

	"instances.pay_plans",
	"instances.photos",
	"instances.instances",

	"products.characteristics",
	"products.photos",
	"products.products",

	"addresses.addresses",
	"categories.categories",
	"periods.periods",
	"payments.methods",
	"deliveries.companies",
	"photos.photos",
	"photos.temp",
	"currencies.currencies",
}

func (self *Inserter) ClearDB() {
	for _, table := range tables {
		clearTable(self.db, table)
	}

	insertCurrencies(self.db)
}

func clearTable(db *sqlx.DB, table string) {
	callWrap(db, fmt.Sprintf("delete from %v", table))
}

func callWrap(db *sqlx.DB, query string, args ...any) {
	_, err := db.Exec(query, args...)

	if nil != err {
		panic(err)
	}
}

func nullWrap[T any, N any](value *T, f func(*T) N) N {
	var empty N
	return nullable.GetOr(
		nullable.Map(
			nullable.FromPtr(value),
			f,
		),
		empty,
	)
}

func NullString(value *string) sql.NullString {
	return nullWrap(value, func(s *string) sql.NullString {
		return sql.NullString{
			String: *s,
			Valid:  true,
		}
	})
}

func NullTime(value *time.Time) sql.NullTime {
	return nullWrap(value, func(t *time.Time) sql.NullTime {
		return sql.NullTime{
			Time:  *t,
			Valid: true,
		}
	})
}

func NullUUID(value *uuid.UUID) uuid.NullUUID {
	return nullWrap(value, func(id *uuid.UUID) uuid.NullUUID {
		return uuid.NullUUID{
			UUID:  *id,
			Valid: true,
		}
	})
}

func BulkInsert[T any](inserter func(*T), values ...T) {
	for _, v := range values {
		inserter(&v)
	}
}

func GetCurrency(name string) uuid.UUID {
	if value, found := currencies[name]; found {
		return value
	} else {
		return currencies["unknown"]
	}
}

func GetAllCurrencies() []uuid.UUID {
	out := make([]uuid.UUID, len(currencies))
	i := 0

	for _, v := range currencies {
		out[i] = v
		i++
	}

	return out
}

var insertCurrencyQuery = prepareInsert("currencies.currencies",
	"id", "name")

func insertCurrencies(db *sqlx.DB) {
	for name, id := range currencies {
		callWrap(db, insertCurrencyQuery, id, name)
	}
}

var insertAddressQuery = prepareInsert("addresses.addresses",
	"id", "country", "city", "street", "house", "flat")

func (self *Inserter) InsertAddress(value *models.Address) {
	callWrap(self.db, insertAddressQuery,
		value.Id, value.Country, value.City, value.Street, value.House,
		NullString(value.Flat),
	)
}

var insertAdministratorQuery = prepareInsert("roles.administrators",
	"id", "user_id")

func (self *Inserter) InsertAdministrator(value *models.Administrator) {
	callWrap(self.db, insertAdministratorQuery, value.Id, value.UserId)
}

var insertCategoryQuery = prepareInsert("categories.categories",
	"id", "name", "parent_id")

func (self *Inserter) InsertCategory(value *models.Category) {
	callWrap(self.db, insertCategoryQuery, value.Id, value.Name,
		NullUUID(value.ParentId))
}

var insertDeliveryCompanyQuery = prepareInsert("deliveries.companies",
	"id", "name", "site", "phone_bumber", "description")

func (self *Inserter) InsertDeliveryCompany(value *models.DeliveryCompany) {
	callWrap(self.db, insertDeliveryCompanyQuery, value.Id, value.Name,
		value.Site, value.PhoneNumber, value.Description)
}

var insertInstanceQuery = prepareInsert("instances.instances",
	"id", "product_id", "name", "description", "condition")

func (self *Inserter) InsertInstance(value *models.Instance) {
	callWrap(self.db, insertInstanceQuery, value.Id, value.ProductId,
		value.Name, value.Description, value.Condition)
}

var insertInstancePayPlansQuery = prepareInsert("instances.pay_plans",
	"id", "instance_id", "period_id", "currency_id", "price")

func (self *Inserter) InsertInstancePayPlans(value *models.InstancePayPlans) {
	for _, v := range value.Map {
		callWrap(self.db, insertInstancePayPlansQuery, v.Id, value.InstanceId,
			v.PeriodId, GetCurrency(v.Price.Name), v.Price.Value)
	}
}

var insertPaymentQuery = prepareInsert("payments.payments",
	"id", "rent_id", "pay_method_id", "payment_id", "period_strat",
	"period_end", "currency_id", "value", "status", "create_date",
	"payment_date")

func (self *Inserter) InsertPayment(value *models.Payment) {
	callWrap(self.db, insertPaymentQuery, value.Id, value.RentId,
		NullUUID(value.PayMethodId), NullString(value.PaymentId),
		value.PeriodStart, value.PeriodEnd, GetCurrency(value.Value.Name),
		value.Value.Value, value.Status, value.CreateDate,
		NullTime(value.PaymentDate),
	)
}

var insertPayMethodQuery = prepareInsert("payments.methods",
	"id", "name", "description")

func (self *Inserter) InsertPayMethod(value *models.PayMethod) {
	callWrap(self.db, insertPayMethodQuery, value.Id, value.Name, value.Description)
}

var insertPeriodQuery = prepareInsert("periods.periods",
	"id", "name", "duration")

func (self *Inserter) InsertPeriod(value *models.Period) {
	callWrap(self.db, insertPeriodQuery, value.Id, value.Name, value.Duration)
}

var insertPhotoQuery = prepareInsert("photos.photos",
	"id", "placeholder", "description", "path", "mime", "date")

func (self *Inserter) InsertPhoto(value *models.Photo) {
	callWrap(self.db, insertPhotoQuery, value.Id, value.Placeholder,
		value.Description, value.Path, value.Mime, value.Date)
}

var insertTempPhotoQuery = prepareInsert("photos.temp",
	"id", "placeholder", "description", "path", "mime", "date")

func (self *Inserter) InsertTempPhoto(value *models.TempPhoto) {
	callWrap(self.db, insertTempPhotoQuery, value.Id, value.Placeholder,
		value.Description, NullString(value.Path), value.Mime, value.Create)
}

var insertPickUpPointQuery = prepareInsert("pick_up_points.pick_up_points",
	"id", "address_id", "capacity")

func (self *Inserter) InsertPickUpPoint(value *models.PickUpPoint) {
	self.InsertAddress(&value.Address)
	callWrap(self.db, insertPickUpPointQuery, value.Id, value.Address.Id,
		value.Capacity)
}

var insertPickUpPointWorkingHoursQuery = prepareInsert(
	"pick_up_points.working_hours",
	"id", "pick_up_point_id", "day", "start_time", "end_time",
)

type ctime struct {
	Time time.Duration
}

func (self ctime) Value() (driver.Value, error) {
	return time.Time{}.Add(self.Time).Format("15:04:05"), nil
}

func (self *Inserter) InsertPickUpPointWorkingHours(value *models.PickUpPointWorkingHours) {
	for _, v := range value.Map {
		callWrap(self.db, insertPickUpPointWorkingHoursQuery, v.Id,
			value.PickUpPointId, v.Day, ctime{v.Begin}, ctime{v.End})
	}
}

var insertProductQuery = prepareInsert("products.products",
	"id", "name", "category_id", "description")

func (self *Inserter) InsertProduct(value *models.Product) {
	callWrap(self.db, insertProductQuery, value.Id, value.Name,
		value.CategoryId, value.Description)
}

var insertProductCharacteristicsQuery = prepareInsert("products.characteristics",
	"id", "product_id", "name", "value")

func (self *Inserter) InsertProductCharacteristics(value *models.ProductCharacteristics) {
	for _, v := range value.Map {
		callWrap(self.db, insertProductCharacteristicsQuery, v.Id,
			value.ProductId, v.Name, v.Value)
	}
}

var insertRenterQuery = prepareInsert("roles.renters",
	"id", "user_id")

func (self *Inserter) InsertRenter(value *models.Renter) {
	callWrap(self.db, insertRenterQuery, value.Id, value.UserId)
}

var insertReviewQuery = prepareInsert("instances.reviews",
	"id", "instance_id", "user_id", "content", "rating", "date")

func (self *Inserter) InsertReview(value *models.Review) {
	callWrap(self.db, insertReviewQuery, value.Id, value.InstanceId,
		value.UserId, value.Content, value.Rating, value.Date)
}

var insertStorekeeperQuery = prepareInsert("roles.storekeepers",
	"id", "user_id", "pick_up_point_id")

func (self *Inserter) InsertStorekeeper(value *models.Storekeeper) {
	callWrap(self.db, insertStorekeeperQuery, value.Id, value.UserId, value.PickUpPointId)
}

var insertUserQuery = prepareInsert("users.users",
	"id", "token", "name", "email", "password")

func (self *Inserter) InsertUser(value *models.User) {
	_, _ = self.db.Exec(insertUserQuery, value.Id, value.Token, value.Name,
		value.Email, value.Password)
}

var insertUserProfileQuery = prepareInsert("users.profiles",
	"id", "user_id", "name", "surname", "patronymic", "birth_date", "photo_id")

func (self *Inserter) InsertUserProfile(value *models.UserProfile) {
	callWrap(self.db, insertUserProfileQuery, value.Id, value.UserId,
		NullString(value.Name), NullString(value.Surname),
		NullString(value.Patronymic), NullTime(value.BirthDate),
		NullUUID(value.PhotoId))
}

var insertUserFavoritePickUpPointQuery = prepareInsert(
	"users.favorite_pick_up_points",
	"id", "user_id", "pick_up_point_id",
)

func (self *Inserter) InsertUserFavoritePickUpPoint(value *models.UserFavoritePickUpPoint) {
	callWrap(self.db, insertUserFavoritePickUpPointQuery, value.Id,
		value.UserId, NullUUID(value.PickUpPointId))
}

var insertUserPayMethodsQuery = prepareInsert("payments.users_methods",
	"id", "pay_method_id", "payer_id", "user_id", "name", "priority")

func (self *Inserter) InsertUserPayMethods(value *models.UserPayMethods) {
	for _, v := range value.Map {
		callWrap(self.db, insertUserPayMethodsQuery, uuidgen.Generate(),
			v.MethodId, v.PayerId, value.UserId, v.Name, v.Priority)
	}
}

var insertProvisionQuery = prepareInsert("records.renters_instances",
	"id", "renter_id", "instance_id", "start_date", "end_date")

func (self *Inserter) InsertProvision(value *records.Provision) {
	callWrap(self.db, insertProvisionQuery, value.Id, value.RenterId,
		value.InstanceId, value.StartDate, NullTime(value.EndDate))
}

var insertRentQuery = prepareInsert("records.users_rents",
	"id", "user_id", "instance_id", "start_date", "end_date",
	"payment_period_id")

func (self *Inserter) InsertRent(value *records.Rent) {
	callWrap(self.db, insertRentQuery, value.Id, value.UserId,
		value.InstanceId, value.StartDate, NullTime(value.EndDate),
		value.PaymentPeriodId)
}

var insertStorageQuery = prepareInsert("records.pick_up_points_instances",
	"id", "pick_up_point_id", "instance_id", "in_date", "out_date")

func (self *Inserter) InsertStorage(value *records.Storage) {
	callWrap(self.db, insertStorageQuery, value.Id, value.PickUpPointId,
		value.InstanceId, value.InDate, NullTime(value.OutDate))
}

var insertDeliveryQuery = prepareInsert("deliveries.deliveries",
	"id", "company_id", "instance_id", "from_id", "to_id", "delivery_id",
	"scheduled_begin_date", "actual_begin_date", "scheduled_end_date",
	"actual_end_date", "verification_code", "create_date")

func (self *Inserter) InsertDelivery(value *requests.Delivery) {
	callWrap(self.db, insertDeliveryQuery, value.Id, value.CompanyId,
		value.InstanceId, value.FromId, value.ToId, value.DeliveryId,
		value.ScheduledBeginDate, NullTime(value.ActualBeginDate),
		value.ScheduledEndDate, NullTime(value.ActualEndDate),
		value.VerificationCode, value.CreateDate)
}

var insertProvideReuqestQuery = prepareInsert("provisions.requests",
	"id", "product_id", "renter_id", "pick_up_point_id", "name", "description",
	"condition", "verification_code", "create_date")
var insertProvideReuqestPPQuery = prepareInsert("provisions.requests_pay_plans",
	"id", "request_id", "period_id", "currency_id", "value")

func (self *Inserter) InsertProvideRequest(value *requests.Provide) {
	callWrap(self.db, insertProvideReuqestQuery, value.Id, value.ProductId,
		value.ProductId, value.PickUpPointId, value.Name, value.Description,
		value.Condition, value.VerificationCode, value.CreateDate)

	for _, v := range value.PayPlans {
		callWrap(self.db, insertProvideReuqestPPQuery, v.Id, value.Id,
			v.PeriodId, GetCurrency(v.Price.Name), v.Price.Value)
	}
}

var insertRentRequestQuery = prepareInsert("rents.requests",
	"id", "instance_id", "user_id", "pick_up_point_id", "payment_period_id",
	"verification_code", "create_date")

func (self *Inserter) InsertRentRequest(value *requests.Rent) {
	callWrap(self.db, insertRentRequestQuery, value.Id,
		value.InstanceId, value.UserId, value.PickUpPointId,
		value.PaymentPeriodId, value.VerificationCode, value.CreateDate)
}

var insertReturnRequestQuery = prepareInsert("rents.returns",
	"id", "instance_id", "user_id", "pick_up_point_id", "rent_end_date",
	"verification_code", "create_date")

func (self *Inserter) InsertReturnRequest(value *requests.Return) {
	callWrap(self.db, insertReturnRequestQuery, value.Id, value.InstanceId,
		value.UserId, value.PickUpPointId, value.RentEndDate,
		value.VerificationCode, value.CreateDate)
}

var insertRevokeRequestQuery = prepareInsert("provisions.revokes",
	"id", "instance_id", "renter_id", "pick_up_point_id", "verification_code",
	"create_date")

func (self *Inserter) InsertRevokeRequest(value *requests.Revoke) {
	callWrap(self.db, insertRevokeRequestQuery, value.Id, value.InstanceId,
		value.RenterId, value.PickUpPointId, value.VerificationCode,
		value.CreateDate)
}

type Photo struct {
	Id     *nullable.Nullable[uuid.UUID]
	Target uuid.UUID
	Photo  uuid.UUID
}

func NewPhoto(id *nullable.Nullable[uuid.UUID], target uuid.UUID, photo uuid.UUID) *Photo {
	return &Photo{id, target, photo}
}

func (self *Photo) getId() uuid.UUID {
	return nullable.GetOrFunc(self.Id, uuidgen.Generate)
}

var insertInstancePhotoQuery = prepareInsert("instances.photos",
	"id", "instance_id", "photo_id")

func (self *Inserter) InsertInstancePhoto(pair *Photo) {
	callWrap(self.db, insertInstancePhotoQuery, pair.getId(),
		pair.Target, pair.Photo)
}

var insertPickUpPointPhotoQuery = prepareInsert("pick_up_points.photos",
	"id", "pick_up_point_id", "photo_id")

func (self *Inserter) InsertPickUpPointPhoto(pair *Photo) {
	callWrap(self.db, insertPickUpPointPhotoQuery, pair.getId(),
		pair.Target, pair.Photo)
}

var insertProductPhotoQuery = prepareInsert("products.photos",
	"id", "product_id", "photo_id")

func (self *Inserter) InsertProductPhoto(pair *Photo) {
	callWrap(self.db, insertProductPhotoQuery, pair.getId(),
		pair.Target, pair.Photo)
}

