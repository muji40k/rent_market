package pickuppoint

import (
	"fmt"
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	sqlCollection "rent_service/internal/repository/implementation/sql/collection"
	"rent_service/internal/repository/implementation/sql/exist"
	"rent_service/internal/repository/implementation/sql/repositories/address"
	"rent_service/internal/repository/implementation/sql/technical"
	"rent_service/internal/repository/interfaces/pickuppoint"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PickUpPoint struct {
	Id        uuid.UUID `db:"id"`
	AddressId uuid.UUID `db:"address_id"`
	Capacity  uint64    `db:"capacity"`
	technical.Info
}

type WorkingHours struct {
	Id            uuid.UUID     `db:"id"`
	PickUpPointId uuid.UUID     `db:"pick_up_point_id"`
	Day           time.Weekday  `db:"day"`
	StartTime     time.Duration `db:"start_time"`
	EndTime       time.Duration `db:"end_time"`
	technical.Info
}

type repository struct {
	connection *sqlx.DB
	address    *address.Repository
}

func New(
	connection *sqlx.DB,
	address *address.Repository,
) pickuppoint.IRepository {
	return &repository{connection, address}
}

func (self *repository) mapf(value *PickUpPoint) models.PickUpPoint {
	out := models.PickUpPoint{
		Id:       value.Id,
		Capacity: value.Capacity,
	}

	addr, err := self.address.GetById(value.AddressId)

	if nil == err {
		out.Address = addr
	}

	return out
}

const get_by_id_query string = `
    select * from pick_up_points.pick_up_points where id = $1
`

func (self *repository) GetById(
	pickUpPointId uuid.UUID,
) (models.PickUpPoint, error) {
	var out PickUpPoint
	err := CheckExistsById(self.connection, pickUpPointId)

	if nil == err {
		err = self.connection.Get(&out, get_by_id_query, pickUpPointId)
	}

	return self.mapf(&out), err
}

const get_all_query string = `
    select * from pick_up_points.pick_up_points offset $1
`

func (self *repository) GetAll() (collection.Collection[models.PickUpPoint], error) {
	return collection.MapCollection(
		self.mapf,
		sqlCollection.New[PickUpPoint](func(offset uint) (*sqlx.Rows, error) {
			return self.connection.Queryx(get_all_query, offset)
		}),
	), nil
}

var count_by_id_query string = exist.GenericCounter("pick_up_points.pick_up_points")

func CheckExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("pick_up_point_id", db, count_by_id_query, id)
}

type photoRepository struct {
	connection *sqlx.DB
}

func NewPhoto(connection *sqlx.DB) pickuppoint.IPhotoRepository {
	return &photoRepository{connection}
}

const get_photo_by_id_query string = `
    select photo_id
    from pick_up_points.photos
    where pick_up_point_id = $1
    offset $2
`

func (self *photoRepository) GetById(
	pickUpPointId uuid.UUID,
) (collection.Collection[uuid.UUID], error) {
	if err := CheckExistsById(self.connection, pickUpPointId); nil != err {
		return nil, err
	}

	return sqlCollection.NewDirect[uuid.UUID](func(offset uint) (*sqlx.Rows, error) {
		return self.connection.Queryx(get_photo_by_id_query, pickUpPointId, offset)
	}), nil
}

type workingHoursRepository struct {
	connection *sqlx.DB
}

func NewWorkingHours(connection *sqlx.DB) pickuppoint.IWorkingHoursRepository {
	return &workingHoursRepository{connection}
}

const get_wh_by_id_query string = `
    select *
    from pick_up_points.working_hours
    where pick_up_point_id = $1
`

func (self *workingHoursRepository) GetById(
	pickUpPointId uuid.UUID,
) (models.PickUpPointWorkingHours, error) {
	var wh models.PickUpPointWorkingHours
	var rows *sqlx.Rows
	err := CheckExistsById(self.connection, pickUpPointId)

	if nil == err {
		rows, err = self.connection.Queryx(get_wh_by_id_query, pickUpPointId)
	}

	if nil == err {
		wh.PickUpPointId = pickUpPointId
		wh.Map = make(map[time.Weekday]models.WorkingHours)

		for cwh := (WorkingHours{}); nil == err && rows.Next(); {
			if err = rows.StructScan(&cwh); nil == err {
				if _, found := wh.Map[cwh.Day]; found {
					err = fmt.Errorf(
						"Duplicate day '%v' was found for pick up point '%v'",
						cwh.Day, pickUpPointId,
					)
				} else {
					wh.Map[cwh.Day] = models.WorkingHours{
						Id:    cwh.Id,
						Day:   cwh.Day,
						Begin: cwh.StartTime,
						End:   cwh.EndTime,
					}
				}
			}
		}

		if nil != err {
			wh = models.PickUpPointWorkingHours{}
		}
	}

	return wh, err
}

