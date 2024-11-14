package photo

import (
	"database/sql"
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/repository/errors/cmnerrors"
	"rent_service/internal/repository/implementation/sql/exist"
	gen_uuid "rent_service/internal/repository/implementation/sql/generate/uuid"
	"rent_service/internal/repository/implementation/sql/technical"
	"rent_service/internal/repository/implementation/sql/utctime"
	"rent_service/internal/repository/interfaces/photo"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Photo struct {
	Id          uuid.UUID       `db:"id"`
	Placeholder string          `db:"placeholder"`
	Description string          `db:"description"`
	Path        string          `db:"path"`
	Mime        string          `db:"mime"`
	Date        utctime.UTCTime `db:"date"`
	technical.Info
}

type TempPhoto struct {
	Id          uuid.UUID       `db:"id"`
	Placeholder string          `db:"placeholder"`
	Description string          `db:"description"`
	Path        sql.NullString  `db:"path"`
	Mime        string          `db:"mime"`
	Create      utctime.UTCTime `db:"date"`
	technical.Info
}

type repository struct {
	connection *sqlx.DB
	setter     technical.ISetter
}

func New(connection *sqlx.DB, setter technical.ISetter) photo.IRepository {
	return &repository{connection, setter}
}

func mapf(value *Photo) models.Photo {
	return models.Photo{
		Id:          value.Id,
		Path:        value.Path,
		Mime:        value.Mime,
		Placeholder: value.Placeholder,
		Description: value.Description,
		Date:        value.Date.Time,
	}
}

func unmapf(value *models.Photo) Photo {
	return Photo{
		Id:          value.Id,
		Path:        value.Path,
		Mime:        value.Mime,
		Placeholder: value.Placeholder,
		Description: value.Description,
		Date:        utctime.FromTime(value.Date),
	}
}

const insert_query string = `
    insert into photos.photos (id, placeholder, description, "path", mime, "date", modification_date, modification_source)
    values (:id, :placeholder, :description, :path, :mime, :date, :modification_date, :modification_source)
`

func (self *repository) Create(photo models.Photo) (models.Photo, error) {
	err := CheckExistsByPath(self.connection, photo.Path)

	if nil == err {
		err = cmnerrors.Duplicate("photo_path")
	} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = nil
	}

	if nil == err {
		photo.Id, err = gen_uuid.GenerateAvailable(
			self.connection,
			CheckExistsById,
		)
	}

	if nil == err {
		photo.Date = time.Now()
		mapped := unmapf(&photo)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(insert_query, mapped)
	}

	return photo, err
}

const get_by_id_query string = `
    select * from photos.photos where id = $1
`

func (self *repository) GetById(photoId uuid.UUID) (models.Photo, error) {
	var out Photo
	err := CheckExistsById(self.connection, photoId)

	if nil == err {
		err = self.connection.Get(&out, get_by_id_query, photoId)
	}

	return mapf(&out), err
}

var count_by_id_query string = exist.GenericCounter("photos.photos")

func CheckExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("photo_id", db, count_by_id_query, id)
}

const count_by_path_query string = `
    select count(*) from photos.photos where path = $1
`

func CheckExistsByPath(db *sqlx.DB, path string) error {
	return exist.Check("photo_path", db, count_by_path_query, path)
}

type tempRepository struct {
	connection *sqlx.DB
	setter     technical.ISetter
}

func NewTemp(
	connection *sqlx.DB,
	setter technical.ISetter,
) photo.ITempRepository {
	return &tempRepository{connection, setter}
}

func mapTemp(value *TempPhoto) models.TempPhoto {
	out := models.TempPhoto{
		Id:          value.Id,
		Mime:        value.Mime,
		Placeholder: value.Placeholder,
		Description: value.Description,
		Create:      value.Create.Time,
	}

	if value.Path.Valid {
		out.Path = new(string)
		*out.Path = value.Path.String
	}

	return out
}

func unmapTemp(value *models.TempPhoto) TempPhoto {
	out := TempPhoto{
		Id:          value.Id,
		Mime:        value.Mime,
		Placeholder: value.Placeholder,
		Description: value.Description,
		Create:      utctime.FromTime(value.Create),
	}

	if nil != value.Path {
		out.Path.Valid = true
		out.Path.String = *value.Path
	}

	return out
}

const insert_temp_query string = `
    insert into photos.temp (id, placeholder, description, "path", mime, "date", modification_date, modification_source)
    values (:id, :placeholder, :description, :path, :mime, :date, :modification_date, :modification_source)
`

func (self *tempRepository) Create(
	photo models.TempPhoto,
) (models.TempPhoto, error) {
	var err error

	if nil != photo.Path {
		err = CheckTempExistsByPath(self.connection, *photo.Path)

		if nil == err {
			err = cmnerrors.Duplicate("photo_temp_path")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	if nil == err {
		photo.Id, err = gen_uuid.GenerateAvailable(
			self.connection,
			CheckTempExistsById,
		)
	}

	if nil == err {
		photo.Create = time.Now()
		mapped := unmapTemp(&photo)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(insert_temp_query, mapped)
	}

	return photo, err
}

const update_temp_query string = `
    update photos.temp
    set placeholder=:placeholder, description=:description, "path"=:path,
        mime=:mime, "date"=:date, modification_date=:modification_date,
        modification_source=:modification_source
    where id=:id
`

func (self *tempRepository) Update(photo models.TempPhoto) error {
	err := CheckTempExistsById(self.connection, photo.Id)

	if nil == err {
		mapped := unmapTemp(&photo)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(update_temp_query, mapped)
	}

	return err
}

const get_temp_by_id_query string = `
    select * from photos.temp where id = $1
`

func (self *tempRepository) GetById(
	photoId uuid.UUID,
) (models.TempPhoto, error) {
	var out TempPhoto
	err := CheckTempExistsById(self.connection, photoId)

	if nil == err {
		err = self.connection.Get(&out, get_temp_by_id_query, photoId)
	}

	return mapTemp(&out), err
}

const delete_by_id_query string = `
    delete from photos."temp" where id = $1
`

func (self *tempRepository) Remove(photoId uuid.UUID) error {
	err := CheckTempExistsById(self.connection, photoId)

	if nil == err {
		_, err = self.connection.Exec(delete_by_id_query, photoId)
	}

	return err
}

var count_temp_by_id_query string = exist.GenericCounter("photos.temp")

func CheckTempExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("photo_temp_id", db, count_temp_by_id_query, id)
}

const count_temp_by_path_query string = `
    select count(*) from photos.temp where path = $1
`

func CheckTempExistsByPath(db *sqlx.DB, path string) error {
	return exist.Check("photo_temp_path", db, count_temp_by_path_query, path)
}

