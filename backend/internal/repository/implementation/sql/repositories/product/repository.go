package product

import (
	"fmt"
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	sqlCollection "rent_service/internal/repository/implementation/sql/collection"
	"rent_service/internal/repository/implementation/sql/exist"
	"rent_service/internal/repository/implementation/sql/repositories/category"
	"rent_service/internal/repository/implementation/sql/technical"
	"rent_service/internal/repository/interfaces/product"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Product struct {
	Id          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	CategoryId  uuid.UUID `db:"category_id"`
	Description string    `db:"description"`
	technical.Info
}

type Characteristic struct {
	Id        uuid.UUID `db:"id"`
	ProductId uuid.UUID `db:"product_id"`
	Name      string    `db:"name"`
	Value     string    `db:"value"`
	technical.Info
}

type repository struct {
	connection *sqlx.DB
}

func New(connection *sqlx.DB) product.IRepository {
	return &repository{connection}
}

func mapf(value *Product) models.Product {
	return models.Product{
		Id:          value.Id,
		Name:        value.Name,
		CategoryId:  value.CategoryId,
		Description: value.Description,
	}
}

const get_by_id_query string = `
    select id, name, category_id, description, modification_date,
           modification_source
    from products.products
    where id = $1
`

func (self *repository) GetById(productId uuid.UUID) (models.Product, error) {
	var product Product
	err := CheckExistsById(self.connection, productId)

	if nil == err {
		err = self.connection.Get(&product, get_by_id_query, productId)
	}

	return mapf(&product), err
}

const count_instances_by_product string = `
    select count(*)
    from (
        select instances.id
        from instances.instances
        where instances.product_id = products.id
    ) as i join (
        select storage.instance_id
        from records.pick_up_points_instances as storage
        where out_date is null
    ) as s on i.id = s.instance_id
`

const exists_characteristics_by_filter_query string = `
    select %v = (
        select count(*)
        from products.characteristics
        where product_id = products.id and (%v)
    )
`

const get_with_filter_query string = `
    select id, name, category_id, description, modification_date,
           modification_source
    from (
        select id, name, category_id, description, modification_date,
               modification_source, (%v) as cnt,
               (%v) as rank
        from products.products
        where category_id = $1 and (%v)
    ) where rank > 0
    order by %v
    offset $2
`

func getCharacteristics(ranges []product.Range, selectors []product.Selector) string {
	if nil == ranges && nil == selectors {
		return "true"
	}

	var iter collection.Iterator[string]

	if nil != ranges {
		iter = collection.MapIterator(
			func(rng *product.Range) string {
				return fmt.Sprintf(
					"(name = '%v' and value >= '%v'::float8 and value <= '%v'::float8)",
					rng.Key, rng.Min, rng.Max,
				)
			},
			collection.SliceIterator(ranges),
		)
	}

	if nil != selectors {
		iter = collection.ChainIterator(
			iter,
			collection.MapIterator(
				func(selector *product.Selector) string {
					filter, any := collection.Reduce(
						collection.MapIterator(
							func(s *string) string {
								return fmt.Sprintf("'%v'", *s)
							},
							collection.SliceIterator(selector.Values),
						),
						func(a *string, b *string) string {
							return *a + ", " + *b
						},
					)

					if !any {
						panic("Iterator was filtered from empty selectors...")
					}

					return fmt.Sprintf(
						"(name = '%v' and value in (%v))",
						selector.Key, filter,
					)
				},
				collection.FilterIterator(
					func(item *product.Selector) bool {
						return nil != item.Values && 0 != len(item.Values)
					},
					collection.SliceIterator(selectors),
				),
			),
		)
	}

	filter, any := collection.Reduce(iter, func(a *string, b *string) string {
		return *a + " or " + *b
	})

	if !any {
		return "true"
	} else {
		return fmt.Sprintf(
			exists_characteristics_by_filter_query,
			len(selectors)+len(ranges), filter,
		)
	}
}

func getSearch(query *string) string {
	if nil == query {
		return "1"
	}

	return fmt.Sprintf(
		"ts_rank_cd(ts, websearch_to_tsquery('russian', '%v'))", *query,
	)
}

func getSort(sort product.Sort) (counter string, sortBy string) {
	switch sort {
	case product.SORT_NONE:
		return "1", "rank desc"
	case product.SORT_OFFERS_ASC:
		return count_instances_by_product, "cnt asc, rank desc"
	case product.SORT_OFFERS_DSC:
		return count_instances_by_product, "cnt desc, rank desc"
	default:
		return "1", "rank desc"
	}
}

func getQuery(filter product.Filter, sort product.Sort) string {
	search := getSearch(filter.Query)
	chars := getCharacteristics(filter.Ranges, filter.Selectors)
	cnt, sortBy := getSort(sort)

	return fmt.Sprintf(
		get_with_filter_query, cnt, search, chars, sortBy,
	)
}

func (self *repository) GetWithFilter(
	filter product.Filter,
	sort product.Sort,
) (collection.Collection[models.Product], error) {
	if err := category.CheckExistsById(self.connection, filter.CategoryId); nil != err {
		return nil, err
	}

	query := getQuery(filter, sort)

	return collection.MapCollection(
		mapf,
		sqlCollection.New[Product](func(offset uint) (*sqlx.Rows, error) {
			return self.connection.Queryx(query, filter.CategoryId, offset)
		}),
	), nil
}

var count_by_id_query string = exist.GenericCounter("products.products")

func CheckExistsById(db *sqlx.DB, id uuid.UUID) error {
	return exist.Check("product_id", db, count_by_id_query, id)
}

type characteristicsRepository struct {
	connection *sqlx.DB
}

func NewCharacteristics(connection *sqlx.DB) product.ICharacteristicsRepository {
	return &characteristicsRepository{connection}
}

func mapCharacteristic(value *Characteristic) models.Charachteristic {
	return models.Charachteristic{
		Id:    value.Id,
		Name:  value.Name,
		Value: value.Value,
	}
}

const get_characteristic_by_id string = `
    select * from products.characteristics where product_id = $1
`

func (self *characteristicsRepository) GetByProductId(
	productId uuid.UUID,
) (models.ProductCharacteristics, error) {
	var rows *sqlx.Rows
	var out = models.ProductCharacteristics{
		ProductId: productId,
		Map:       make(map[string]models.Charachteristic),
	}
	err := CheckExistsById(self.connection, productId)

	if nil == err {
		rows, err = self.connection.Queryx(get_characteristic_by_id, productId)
	}

	for mapped := (Characteristic{}); nil == err && rows.Next(); {
		if err = rows.StructScan(&mapped); nil == err {
			out.Map[mapped.Name] = mapCharacteristic(&mapped)
		}
	}

	if nil != err {
		out = models.ProductCharacteristics{}
	}

	return out, err
}

type photoRepository struct {
	connection *sqlx.DB
}

func NewPhoto(connection *sqlx.DB) product.IPhotoRepository {
	return &photoRepository{connection}
}

const get_photo_by_id_query string = `
    select photo_id
    from products.photos
    where product_id = $1
    offset $2
`

func (self *photoRepository) GetByProductId(
	productId uuid.UUID,
) (collection.Collection[uuid.UUID], error) {
	if err := CheckExistsById(self.connection, productId); nil != err {
		return nil, err
	}

	return sqlCollection.NewDirect[uuid.UUID](func(offset uint) (*sqlx.Rows, error) {
		return self.connection.Queryx(get_photo_by_id_query, productId, offset)
	}), nil
}

