package user

import (
	"database/sql"
	"errors"
	"fmt"
	"rent_service/internal/domain/models"
	"rent_service/internal/repository/errors/cmnerrors"
	"rent_service/internal/repository/implementation/sql/exist"
	gen_uuid "rent_service/internal/repository/implementation/sql/generate/uuid"
	"rent_service/internal/repository/implementation/sql/repositories/paymethod"
	"rent_service/internal/repository/implementation/sql/repositories/photo"
	"rent_service/internal/repository/implementation/sql/repositories/pickuppoint"
	"rent_service/internal/repository/implementation/sql/technical"
	"rent_service/internal/repository/implementation/sql/utctime"
	"rent_service/internal/repository/interfaces/user"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Hasher func(*models.User) string

type User struct {
	Id       uuid.UUID `db:"id"`
	Token    string    `db:"token"`
	Name     string    `db:"name"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
	technical.Info
}

type Profile struct {
	Id         uuid.UUID                 `db:"id"`
	UserId     uuid.UUID                 `db:"user_id"`
	Name       sql.NullString            `db:"name"`
	Surname    sql.NullString            `db:"surname"`
	Patronymic sql.NullString            `db:"patronymic"`
	BirthDate  sql.Null[utctime.UTCTime] `db:"birth_date"`
	PhotoId    uuid.NullUUID             `db:"photo_id"`
	technical.Info
}

type FavoritePickUpPoint struct {
	Id            uuid.UUID     `db:"id"`
	UserId        uuid.UUID     `db:"user_id"`
	PickUpPointId uuid.NullUUID `db:"pick_up_point_id"`
	technical.Info
}

type PayMethod struct {
	Id          uuid.UUID `db:"id"`
	PayMethodId uuid.UUID `db:"pay_method_id"`
	PayerId     string    `db:"payer_id"`
	UserId      uuid.UUID `db:"user_id"`
	Name        string    `db:"name"`
	Priority    uint      `db:"priority"`
	technical.Info
}

type repository struct {
	connection *sqlx.DB
	setter     technical.ISetter
	hasher     Hasher
}

func New(
	connection *sqlx.DB,
	setter technical.ISetter,
	hasher Hasher,
) user.IRepository {
	return &repository{connection, setter, hasher}
}

func mapf(value *User) models.User {
	return models.User{
		Id:       value.Id,
		Name:     value.Name,
		Email:    value.Email,
		Password: value.Password,
		Token:    models.Token(value.Token),
	}
}

func unmapf(value *models.User) User {
	return User{
		Id:       value.Id,
		Name:     value.Name,
		Email:    value.Email,
		Password: value.Password,
		Token:    string(value.Token),
	}
}

const insert_query string = `
    insert into users.users (
        id, "token", "name", email, "password",
        modification_date, modification_source
    ) values (
        :id, :token, :name, :email, :password,
        :modification_date, :modification_source
    )
`

func (self *repository) checkCreate(user models.User) error {
	err := CheckExistsByEmail(self.connection, user.Email)

	if nil == err {
		err = cmnerrors.Duplicate("user_email")
	} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = nil
	}

	if nil == err {
		err = CheckExistsByName(self.connection, user.Name)

		if nil == err {
			err = cmnerrors.Duplicate("user_name")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	return err
}

func (self *repository) Create(user models.User) (models.User, error) {
	err := self.checkCreate(user)

	if nil == err {
		user.Id, err = gen_uuid.GenerateAvailable(
			self.connection,
			CheckExistsById,
		)
	}

	if nil == err {
		token := self.hasher(&user)
		err = CheckExistsByToken(self.connection, token)

		if nil == err {
			err = cmnerrors.Duplicate("user_token")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
			user.Token = models.Token(token)
		}
	}

	if nil == err {
		mapped := unmapf(&user)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(insert_query, mapped)
	}

	return user, err
}

const update_query string = `
    update users.users
    set "token"=:token, "name"=:name, email=:email, "password"=:password,
        modification_date=:modification_date,
        modification_source=:modification_source
    where id=:id
`

func (self *repository) checkUserExists(user models.User) error {
	err := CheckExistsById(self.connection, user.Id)

	if nil == err {
		err = CheckExistsByIdAndEmail(self.connection, user.Id, user.Email)

		if nil == err {
			err = cmnerrors.Duplicate("user_email")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	if nil == err {
		err = CheckExistsByIdAndName(self.connection, user.Id, user.Name)

		if nil == err {
			err = cmnerrors.Duplicate("user_name")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	if nil == err {
		token := self.hasher(&user)
		err = CheckExistsByToken(self.connection, token)

		if nil == err {
			err = cmnerrors.Duplicate("user_token")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
			user.Token = models.Token(token)
		}
	}

	return err
}

func (self *repository) Update(user models.User) error {
	err := self.checkUserExists(user)

	if nil == err {
		mapped := unmapf(&user)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(update_query, mapped)
	}

	return err
}

const get_by_id_query string = `
    select * from users.users where id = $1
`

func (self *repository) GetById(userId uuid.UUID) (models.User, error) {
	var user User
	err := CheckExistsById(self.connection, userId)

	if nil == err {
		err = self.connection.Get(&user, get_by_id_query, userId)
	}

	return mapf(&user), err
}

const get_by_email_query string = `
    select * from users.users where email = $1
`

func (self *repository) GetByEmail(email string) (models.User, error) {
	var user User
	err := CheckExistsByEmail(self.connection, email)

	if nil == err {
		err = self.connection.Get(&user, get_by_email_query, email)
	}

	return mapf(&user), err
}

const get_by_token_query string = `
    select * from users.users where token = $1
`

func (self *repository) GetByToken(token models.Token) (models.User, error) {
	var user User
	err := CheckExistsByToken(self.connection, string(token))

	if nil == err {
		err = self.connection.Get(&user, get_by_token_query, token)
	}

	return mapf(&user), err
}

var count_by_id_query string = exist.GenericCounter("users.users")

func CheckExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("user_id", db, count_by_id_query, id)
}

var count_by_token_query string = `
    select count(*) from users.users where token = $1
`

func CheckExistsByToken(db *sqlx.DB, token string) error {
	return exist.Check("user_token", db, count_by_token_query, token)
}

var count_by_email_query string = `
    select count(*) from users.users where email = $1
`

func CheckExistsByEmail(db *sqlx.DB, email string) error {
	return exist.Check("user_email", db, count_by_email_query, email)
}

var count_by_id_and_email_query string = `
    select count(*) from users.users where id != $1 and email = $2
`

func CheckExistsByIdAndEmail(db *sqlx.DB, id uuid.UUID, email string) error {
	return exist.Check("user_email", db, count_by_id_and_email_query, id, email)
}

var count_by_name_query string = `
    select count(*) from users.users where name = $1
`

func CheckExistsByName(db *sqlx.DB, name string) error {
	return exist.Check("user_name", db, count_by_name_query, name)
}

var count_by_id_and_name_query string = `
    select count(*) from users.users where id != $1 and name = $2
`

func CheckExistsByIdAndName(db *sqlx.DB, id uuid.UUID, name string) error {
	return exist.Check("user_name", db, count_by_id_and_name_query, id, name)
}

type profileRepository struct {
	connection *sqlx.DB
	setter     technical.ISetter
}

func NewProfile(
	connection *sqlx.DB,
	setter technical.ISetter,
) user.IProfileRepository {
	return &profileRepository{connection, setter}
}

func setter[T any](dest **T, valid bool, value T) {
	if valid {
		*dest = new(T)
		**dest = value
	}
}

func mapProfile(value *Profile) models.UserProfile {
	out := models.UserProfile{
		Id:     value.Id,
		UserId: value.UserId,
	}

	setter(&out.Name, value.Name.Valid, value.Name.String)
	setter(&out.Surname, value.Surname.Valid, value.Surname.String)
	setter(&out.Patronymic, value.Patronymic.Valid, value.Patronymic.String)
	setter(&out.BirthDate, value.BirthDate.Valid, value.BirthDate.V.Time)
	setter(&out.PhotoId, value.PhotoId.Valid, value.PhotoId.UUID)

	return out
}

func unmapProfile(value *models.UserProfile) Profile {
	out := Profile{
		Id:     value.Id,
		UserId: value.UserId,
	}

	if nil != value.Name {
		out.Name.Valid = true
		out.Name.String = *value.Name
	}

	if nil != value.Surname {
		out.Surname.Valid = true
		out.Surname.String = *value.Surname
	}

	if nil != value.Patronymic {
		out.Patronymic.Valid = true
		out.Patronymic.String = *value.Patronymic
	}

	if nil != value.BirthDate {
		out.BirthDate.Valid = true
		out.BirthDate.V = utctime.FromTime(*value.BirthDate)
	}

	if nil != value.PhotoId {
		out.PhotoId.Valid = true
		out.PhotoId.UUID = *value.PhotoId
	}

	return out
}

const insert_profile_query string = `
    insert into users.profiles (
        id, user_id, "name", surname, patronymic, birth_date, photo_id,
        modification_date, modification_source
    ) values (
        :id, :user_id, :name, :surname, :patronymic, :birth_date, :photo_id,
        :modification_date, :modification_source
    )
`

func (self *profileRepository) Create(
	profile models.UserProfile,
) (models.UserProfile, error) {
	err := CheckExistsById(self.connection, profile.UserId)

	if nil == err {
		err = CheckProfileExistsByUserId(self.connection, profile.UserId)

		if nil == err {
			err = cmnerrors.Duplicate("user_profile_user_id")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	if nil == err && nil != profile.PhotoId {
		err = photo.CheckExistsById(self.connection, *profile.PhotoId)
	}

	if nil == err {
		profile.Id, err = gen_uuid.GenerateAvailable(
			self.connection,
			CheckProfileExistsById,
		)
	}

	if nil == err {
		mapped := unmapProfile(&profile)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(insert_profile_query, mapped)
	}

	return profile, err
}

const update_profile_query string = `
    update users.profiles
    set user_id=:user_id, "name"=:name, surname=:surname,
        patronymic=:patronymic, birth_date=:birth_date, photo_id=:photo_id,
        modification_date=:modification_date,
        modification_source=:modification_source
    where id=:id;
`

func (self *profileRepository) Update(profile models.UserProfile) error {
	err := CheckProfileExistsByIdAndUserId(
		self.connection,
		profile.Id,
		profile.UserId,
	)

	if nil == err && nil != profile.PhotoId {
		err = photo.CheckExistsById(self.connection, *profile.PhotoId)
	}

	if nil == err {
		mapped := unmapProfile(&profile)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(update_profile_query, mapped)
	}

	return err
}

const get_profile_by_user_id = `
    select * from users.profiles where user_id = $1
`

func (self *profileRepository) GetByUserId(
	userId uuid.UUID,
) (models.UserProfile, error) {
	var profile Profile
	err := CheckExistsById(self.connection, userId)

	if nil == err {
		err = CheckProfileExistsByUserId(self.connection, userId)
	}

	if nil == err {
		err = self.connection.Get(&profile, get_profile_by_user_id, userId)
	}

	return mapProfile(&profile), err
}

var count_profile_by_id_query string = exist.GenericCounter("users.profiles")

func CheckProfileExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("user_profile_id", db, count_profile_by_id_query, id)
}

var count_profile_by_user_id_query string = `
    select count(*) from users.profiles where user_id = $1
`

func CheckProfileExistsByUserId(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check(
		"user_profile_user_id",
		db,
		count_profile_by_user_id_query,
		id,
	)
}

var count_profile_by_id_and_user_id_query string = `
    select count(*) from users.profiles where id = $1 and user_id = $2
`

func CheckProfileExistsByIdAndUserId(
	db *sqlx.DB,
	id uuid.UUID,
	userId uuid.UUID,
) error {
	return exist.Check(
		"user_profile_instance",
		db,
		count_profile_by_id_and_user_id_query,
		id, userId,
	)
}

type favoriteRepository struct {
	connection *sqlx.DB
	setter     technical.ISetter
}

func NewFavorite(
	connection *sqlx.DB,
	setter technical.ISetter,
) user.IFavoriteRepository {
	return &favoriteRepository{connection, setter}
}

func mapFavorite(value *FavoritePickUpPoint) models.UserFavoritePickUpPoint {
	out := models.UserFavoritePickUpPoint{
		Id:     value.Id,
		UserId: value.UserId,
	}

	setter(&out.PickUpPointId, value.PickUpPointId.Valid, value.PickUpPointId.UUID)

	return out
}

func unmapFavorite(value *models.UserFavoritePickUpPoint) FavoritePickUpPoint {
	out := FavoritePickUpPoint{
		Id:     value.Id,
		UserId: value.UserId,
	}

	if nil != value.PickUpPointId {
		out.PickUpPointId.Valid = true
		out.PickUpPointId.UUID = *value.PickUpPointId
	}

	return out
}

const insert_favorite_query string = `
    insert into users.favorite_pick_up_points (
        id, user_id, pick_up_point_id, modification_date,
        modification_source
    ) values (
        :id, :user_id, :pick_up_point_id, :modification_date,
        :modification_source
    )
`

func (self *favoriteRepository) Create(
	profile models.UserFavoritePickUpPoint,
) (models.UserFavoritePickUpPoint, error) {
	err := CheckExistsById(self.connection, profile.UserId)

	if nil == err && nil != profile.PickUpPointId {
		err = pickuppoint.CheckExistsById(
			self.connection,
			*profile.PickUpPointId,
		)
	}

	if nil == err {
		err = CheckFavoriteExistsByUserId(self.connection, profile.UserId)

		if nil == err {
			err = cmnerrors.Duplicate("user_favorite_user_id")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	if nil == err {
		profile.Id, err = gen_uuid.GenerateAvailable(
			self.connection,
			CheckFavoriteExistsById,
		)
	}

	if nil == err {
		mapped := unmapFavorite(&profile)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(insert_favorite_query, mapped)
	}

	return profile, err
}

const update_favorite_query string = `
    update users.favorite_pick_up_points
    set user_id=:user_id, pick_up_point_id=:pick_up_point_id,
        modification_date=:modification_date,
        modification_source=:modification_source
    where id=:id
`

func (self *favoriteRepository) Update(
	profile models.UserFavoritePickUpPoint,
) error {
	err := CheckFavoriteExistsByIdAndUserId(
		self.connection,
		profile.Id,
		profile.UserId,
	)

	if nil == err && nil != profile.PickUpPointId {
		err = pickuppoint.CheckExistsById(
			self.connection,
			*profile.PickUpPointId,
		)
	}

	if nil == err {
		mapped := unmapFavorite(&profile)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(update_favorite_query, mapped)
	}

	return err
}

const get_favorite_by_user_id_query string = `
    select * from users.favorite_pick_up_points where user_id = $1
`

func (self *favoriteRepository) GetByUserId(
	userId uuid.UUID,
) (models.UserFavoritePickUpPoint, error) {
	var favorite FavoritePickUpPoint
	err := CheckExistsById(self.connection, userId)

	if nil == err {
		err = CheckFavoriteExistsByUserId(self.connection, userId)
	}

	if nil == err {
		err = self.connection.Get(
			&favorite,
			get_favorite_by_user_id_query,
			userId,
		)
	}

	return mapFavorite(&favorite), err
}

var count_favorite_by_id_query string = exist.GenericCounter("users.favorite_pick_up_points")

func CheckFavoriteExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("user_favorite_id", db, count_favorite_by_id_query, id)
}

const count_favorite_by_user_id_query string = `
    select count(*) from users.favorite_pick_up_points where user_id = $1
`

func CheckFavoriteExistsByUserId(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check(
		"user_favorite_user_id",
		db,
		count_favorite_by_user_id_query,
		id,
	)
}

var count_favorite_by_id_and_user_id_query string = `
    select count(*)
    from users.favorite_pick_up_points
    where id = $1 && user_id = $2
`

func CheckFavoriteExistsByIdAndUserId(
	db *sqlx.DB,
	id uuid.UUID,
	userId uuid.UUID,
) error {
	return exist.Check(
		"user_favorite_instance",
		db,
		count_favorite_by_id_and_user_id_query,
		id, userId,
	)
}

type payMethodsRepository struct {
	connection *sqlx.DB
	setter     technical.ISetter
}

func NewPayMethod(
	connection *sqlx.DB,
	setter technical.ISetter,
) user.IPayMethodsRepository {
	return &payMethodsRepository{connection, setter}
}

func mapPayMethod(value *PayMethod) models.UserPayMethod {
	return models.UserPayMethod{
		Name:     value.Name,
		MethodId: value.PayMethodId,
		PayerId:  value.PayerId,
		Priority: value.Priority,
	}
}

func unmapPayMethod(value *models.UserPayMethod) PayMethod {
	return PayMethod{
		PayMethodId: value.MethodId,
		PayerId:     value.PayerId,
		Name:        value.Name,
		Priority:    value.Priority,
	}
}

const insert_pay_method_query string = `
    insert into payments.users_methods (
        id, pay_method_id, payer_id, user_id, "name", priority,
        modification_date, modification_source
    ) values (
        :id, :pay_method_id, :payer_id, :user_id, :name, :priority,
        :modification_date, :modification_source
    )
`

func (self *payMethodsRepository) CreatePayMethod(
	userId uuid.UUID,
	payMethod models.UserPayMethod,
) (models.UserPayMethods, error) {
	err := CheckExistsById(self.connection, userId)

	if nil == err {
		err = paymethod.CheckExistsById(self.connection, payMethod.MethodId)
	}

	if nil == err {
		err = CheckPayMethodExistsByUserIdAndMethodId(
			self.connection,
			userId, payMethod.MethodId,
		)

		if nil == err {
			err = cmnerrors.Duplicate("user_pay_method_instance")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	var rows *sqlx.Rows
	var out models.UserPayMethods
	var mapped = unmapPayMethod(&payMethod)

	if nil == err {
		mapped.Id, err = gen_uuid.GenerateAvailable(
			self.connection,
			CheckPayMethodExistsById,
		)
	}

	if nil == err {
		mapped.UserId = userId
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(insert_pay_method_query, mapped)
	}

	if nil == err {
		rows, err = self.getByUserId(userId)
	}

	if nil == err {
		out, err = fold(userId, rows)
	}

	return out, err
}

const update_pay_method_query string = `
    update payments.users_methods
    set pay_method_id=:pay_method_id, payer_id=:payer_id, user_id=:user_id,
        name=:name, priority=:priority,
        modification_date=:modification_date,
        modification_source=:modification_source
    where id=:id;
`

const delete_pay_method_query string = `
    delete from payments.users_methods where id=$1
`

func (self *payMethodsRepository) Update(
	payMethods models.UserPayMethods, // <- Update to Match this
) error {
	var rows *sqlx.Rows
	var method PayMethod
	err := CheckExistsById(self.connection, payMethods.UserId)

	if nil == err {
		rows, err = self.getByUserId(payMethods.UserId)
	}

	for nil == err && rows.Next() {
		if err = rows.StructScan(&method); nil == err {
			if umethod, found := payMethods.Map[method.PayMethodId]; !found {
				_, err = self.connection.Exec(
					delete_pay_method_query, method.Id,
				)
			} else {
				delete(payMethods.Map, method.PayMethodId)
				mapped := unmapPayMethod(&umethod)
				mapped.Id = method.Id
				mapped.UserId = method.UserId
				self.setter.Update(&mapped.Info)
				_, err = self.connection.NamedExec(
					update_pay_method_query, mapped,
				)
			}
		}
	}

	if nil == err && 0 != len(payMethods.Map) {
		unknown := make([]uuid.UUID, 0, len(payMethods.Map))

		for _, method := range payMethods.Map {
			unknown = append(unknown, method.MethodId)
		}

		err = fmt.Errorf("Found unknown pay methods: %v", unknown)
	}

	return err
}

const get_pay_method_by_user_id_query string = `
    select * from payments.users_methods where user_id = $1
`

func (self *payMethodsRepository) GetByUserId(
	userId uuid.UUID,
) (models.UserPayMethods, error) {
	var out models.UserPayMethods
	rows, err := self.getByUserId(userId)

	if nil == err {
		out, err = fold(userId, rows)
	}

	return out, err
}

func fold(userId uuid.UUID, rows *sqlx.Rows) (models.UserPayMethods, error) {
	var method PayMethod
	var err error
	var out = models.UserPayMethods{
		UserId: userId,
		Map:    make(map[uuid.UUID]models.UserPayMethod),
	}

	for nil == err && rows.Next() {
		if err = rows.StructScan(&method); nil == err {
			mapped := mapPayMethod(&method)
			out.Map[mapped.MethodId] = mapped
		}
	}

	if nil != err {
		out = models.UserPayMethods{}
	}

	return out, err
}

func (self *payMethodsRepository) getByUserId(
	userId uuid.UUID,
) (*sqlx.Rows, error) {
	var rows *sqlx.Rows
	err := CheckExistsById(self.connection, userId)

	if nil == err {
		err = CheckPayMethodExistsByUserId(self.connection, userId)
	}

	if nil == err {
		rows, err = self.connection.Queryx(
			get_pay_method_by_user_id_query,
			userId,
		)
	}

	return rows, err
}

var count_pay_method_by_id_query string = exist.GenericCounter("payments.users_methods")

func CheckPayMethodExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check(
		"user_pay_method_id",
		db,
		count_pay_method_by_id_query,
		id,
	)
}

const count_pay_method_by_user_id_query string = `
    select count(*) from payments.users_methods where user_id = $1
`

func CheckPayMethodExistsByUserId(
	db *sqlx.DB,
	userId uuid.UUID,
) error {
	return exist.CheckMultiple(
		"user_pay_method_user_id",
		db,
		count_pay_method_by_user_id_query,
		userId,
	)
}

const count_pay_method_by_user_id_and_method_id_query string = `
    select count(*)
    from payments.users_methods
    where user_id = $1 and pay_method_id = $2
`

func CheckPayMethodExistsByUserIdAndMethodId(
	db *sqlx.DB,
	userId uuid.UUID,
	methodId uuid.UUID,
) error {
	return exist.Check(
		"user_pay_method_instance",
		db,
		count_pay_method_by_user_id_and_method_id_query,
		userId, methodId,
	)
}

