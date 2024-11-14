package review

import (
	"errors"
	"fmt"
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/errors/cmnerrors"
	sqlCollection "rent_service/internal/repository/implementation/sql/collection"
	"rent_service/internal/repository/implementation/sql/exist"
	gen_uuid "rent_service/internal/repository/implementation/sql/generate/uuid"
	"rent_service/internal/repository/implementation/sql/repositories/instance"
	"rent_service/internal/repository/implementation/sql/repositories/user"
	"rent_service/internal/repository/implementation/sql/technical"
	"rent_service/internal/repository/implementation/sql/utctime"
	"rent_service/internal/repository/interfaces/review"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Review struct {
	Id         uuid.UUID       `db:"id"`
	InstanceId uuid.UUID       `db:"instance_id"`
	UserId     uuid.UUID       `db:"user_id"`
	Content    string          `db:"content"`
	Rating     float64         `db:"rating"`
	Date       utctime.UTCTime `db:"date"`
	technical.Info
}

type repository struct {
	connection *sqlx.DB
	setter     technical.ISetter
}

func New(connection *sqlx.DB, setter technical.ISetter) review.IRepository {
	return &repository{connection, setter}
}

func mapf(value *Review) models.Review {
	return models.Review{
		Id:         value.Id,
		InstanceId: value.InstanceId,
		UserId:     value.UserId,
		Content:    value.Content,
		Rating:     value.Rating,
		Date:       value.Date.Time,
	}
}

func unmapf(value *models.Review) Review {
	return Review{
		Id:         value.Id,
		InstanceId: value.InstanceId,
		UserId:     value.UserId,
		Content:    value.Content,
		Rating:     value.Rating,
		Date:       utctime.FromTime(value.Date),
	}
}

func mapSort(sort review.Sort) string {
	switch sort {
	case review.SORT_NONE:
		return "date desc"
	case review.SORT_DATE_ASC:
		return "date asc"
	case review.SORT_DATE_DSC:
		return "date desc"
	case review.SORT_RATING_ASC:
		return "rating asc"
	case review.SORT_RATING_DSC:
		return "rating desc"
	default:
		return "date desc"
	}
}

var ranges = map[review.Rating]struct{ min, max float64 }{
	0: {-0.5, 0.5},
	1: {0.5, 1.5},
	2: {1.5, 2.5},
	3: {2.5, 3.5},
	4: {3.5, 4.5},
	5: {4.5, 5.5},
}

func mapRatings(ratings []review.Rating) string {
	filter := "("
	first := true

	for _, rating := range ratings {
		if r, found := ranges[rating]; found {
			if !first {
				filter += " or "
			} else {
				first = false
			}

			filter += fmt.Sprintf(
				"(rating >= %v and rating <= %v)", r.min, r.max,
			)
		}
	}

	filter += ")"

	if first {
		return "(true)"
	} else {
		return fmt.Sprintf("(%v)", filter)
	}
}

const insert_query string = `
    insert into instances.reviews (id, instance_id, user_id, "content", rating, "date", modification_date, modification_source)
    values (:id, :instance_id, :user_id, :content, :rating, :date, :modification_date, :modification_source)
`

func (self *repository) Create(review models.Review) (models.Review, error) {
	err := instance.CheckExistsById(self.connection, review.InstanceId)

	if nil == err {
		err = user.CheckExistsById(self.connection, review.UserId)
	}

	if nil == err {
		review.Id, err = gen_uuid.GenerateAvailable(
			self.connection,
			CheckExistsById,
		)
	}

	if nil == err {
		err = CheckExistsByUserIdAndInstanceId(
			self.connection,
			review.UserId,
			review.InstanceId,
		)

		if nil == err {
			err = cmnerrors.Duplicate("review_user")
		} else if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = nil
		}
	}

	if nil == err {
		mapped := unmapf(&review)
		self.setter.Update(&mapped.Info)
		_, err = self.connection.NamedExec(insert_query, mapped)
	}

	return review, err
}

const get_with_filter_query string = `
    select *
    from instances.reviews
    where instance_id = $1 and %v
    order by %v
    offset $2
`

func (self *repository) GetWithFilter(
	filter review.Filter,
	sort review.Sort,
) (collection.Collection[models.Review], error) {
	if err := instance.CheckExistsById(self.connection, filter.InstanceId); nil != err {
		return nil, err
	}

	ratings := mapRatings(filter.Ratings)
	sortValue := mapSort(sort)
	query := fmt.Sprintf(get_with_filter_query, ratings, sortValue)

	return collection.MapCollection(
		mapf,
		sqlCollection.New[Review](func(offset uint) (*sqlx.Rows, error) {
			return self.connection.Queryx(query, filter.InstanceId, offset)
		}),
	), nil
}

var count_by_id_query string = exist.GenericCounter("instances.reviews")

func CheckExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("review_id", db, count_by_id_query, id)
}

var count_by_user_id_and_instance_id_query string = `
    select count(*)
    from instances.reviews
    where user_id = $1 and instance_id = $2
`

func CheckExistsByUserIdAndInstanceId(
	db *sqlx.DB,
	userId uuid.UUID,
	instanceId uuid.UUID,
) error {
	return exist.Check(
		"review_duplicate",
		db,
		count_by_user_id_and_instance_id_query,
		userId,
		instanceId,
	)
}

