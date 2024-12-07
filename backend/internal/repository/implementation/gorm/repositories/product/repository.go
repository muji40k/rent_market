package product

import (
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
	gormCollection "rent_service/internal/repository/implementation/gorm/collection"
	orm "rent_service/internal/repository/implementation/gorm/models"
	"rent_service/internal/repository/interfaces/product"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) product.IRepository {
	return &repository{db}
}

func (self *repository) GetById(productId uuid.UUID) (models.Product, error) {
	var out orm.Product
	err := self.db.Model(orm.Product{}).
		Where("id = ?", productId.String()).
		First(&out).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = repo_errors.NotFound("product_id")
	}

	if nil == err {
		return orm.MapProduct(&out), nil
	} else {
		return models.Product{}, err
	}
}

func getCnt(db *gorm.DB, sort product.Sort, filter product.Filter) interface{} {
	if product.SORT_NONE == sort {
		return "1"
	}

	// ToDo: Add instance
	return "1"
}

func getRank(db *gorm.DB, query *string) interface{} {
	if nil == query {
		return "1"
	}

	return db.Raw(
		"ts_rank_cd(ts, websearch_to_tsquery('russian', ?))", *query,
	)
}

func getChars(db *gorm.DB, filter product.Filter) interface{} {
	cnt := len(filter.Ranges) + len(filter.Selectors)

	if 0 == cnt {
		return true
	}

	base := db.Model(orm.ProductCharacteristic{}).
		Select("count(*)").
		Where("product_id = products.id")

	chars := db

	for _, rng := range filter.Ranges {
		chars = chars.Or(
			db.Where("name = ?", rng.Key).
				Where("value::float8 >= ?", rng.Min).
				Where("value::float8 <= ?", rng.Max),
		)
	}

	for _, selector := range filter.Selectors {
		chars = chars.Or(
			db.Where("name = ?", selector.Key).
				Where("value in ?", selector.Values),
		)
	}

	return db.Raw("select ? = (?)", cnt, base.Where(chars))
}

func getSort(sort product.Sort) string {
	switch sort {
	case product.SORT_NONE:
		return "rank desc"
	case product.SORT_OFFERS_ASC:
		return "cnt asc, rank desc"
	case product.SORT_OFFERS_DSC:
		return "cnt desc, rank desc"
	default:
		return "rank desc"
	}
}

func (self *repository) GetWithFilter(filter product.Filter, sort product.Sort) (collection.Collection[models.Product], error) {
	tx := self.db.Table("(?) as k",
		self.db.Model(orm.Product{}).
			Select(
				`id, name, category_id, description, modification_date,
                 modification_source, (?) as cnt, (?) as rank`,
				getCnt(self.db, sort, filter), getRank(self.db, filter.Query),
			).
			Where("category_id = ?", filter.CategoryId.String()).
			Where("(?)", getChars(self.db, filter)),
	).
		Select(
			`id, name, category_id, description, modification_date,
             modification_source`,
		).
		Where("rank > 0").
		Order(getSort(sort))

	return collection.MapCollection(
		orm.MapProduct,
		gormCollection.New[orm.Product](self.db, tx),
	), nil
}

