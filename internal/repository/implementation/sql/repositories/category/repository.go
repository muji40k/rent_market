package category

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	sqlCollection "rent_service/internal/repository/implementation/sql/collection"
	"rent_service/internal/repository/implementation/sql/exist"
	"rent_service/internal/repository/implementation/sql/technical"
	"rent_service/internal/repository/interfaces/category"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Category struct {
	Id       uuid.UUID     `db:"id"`
	ParentId uuid.NullUUID `db:"parent_id"`
	Name     string        `db:"name"`
	technical.Info
}

type repository struct {
	connection *sqlx.DB
}

func New(connection *sqlx.DB) category.IRepository {
	return &repository{connection}
}

func mapCategory(value *Category) models.Category {
	out := models.Category{
		Id:   value.Id,
		Name: value.Name,
	}

	if value.ParentId.Valid {
		out.ParentId = new(uuid.UUID)
		*out.ParentId = value.ParentId.UUID
	}

	return out
}

const get_all_query string = "select * from categories.categories offset $1"

func (self *repository) GetAll() (collection.Collection[models.Category], error) {
	return collection.MapCollection(
		mapCategory,
		sqlCollection.New[Category](func(offset uint) (*sqlx.Rows, error) {
			return self.connection.Queryx(get_all_query, offset)
		}),
	), nil
}

const get_path_query string = `
    with recursive path (id, name, parent_id, modification_date, modification_source, depth) as (
        select *, 0 from categories.categories where id = $1
        UNION ALL
        select categories.id, categories.name, categories.parent_id,
               categories.modification_date, categories.modification_source,
               depth+1
        from categories.categories join path
             on categories.id = path.parent_id
    )
    select id, name, parent_id, modification_date, modification_source
    from path
    order by depth desc
    offset $2
`

func (self *repository) GetPath(
	leaf uuid.UUID,
) (collection.Collection[models.Category], error) {
	if err := CheckExistsById(self.connection, leaf); nil != err {
		return nil, err
	}

	return collection.MapCollection(
		mapCategory,
		sqlCollection.New[Category](func(offset uint) (*sqlx.Rows, error) {
			return self.connection.Queryx(get_path_query, leaf, offset)
		}),
	), nil
}

var count_by_id_query string = exist.GenericCounter("categories.categories")

func CheckExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("category_id", db, count_by_id_query, id)
}

