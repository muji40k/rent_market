package product

import (
	"rent_service/cmd/benchmark/utils/metrics"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/interfaces/product"
	"testing"
	"time"

	"github.com/google/uuid"
)

type TestProductRepository struct {
	ctor  func() (product.IRepository, func())
	repo  product.IRepository
	clean func()
}

func New(ctor func() (product.IRepository, func())) *TestProductRepository {
	return &TestProductRepository{ctor, nil, nil}
}

func (self *TestProductRepository) Before() {
	self.repo, self.clean = self.ctor()
}

var query = "product"

var filter = product.Filter{
	CategoryId: uuid.MustParse("923297f7-3dc3-4d02-a4ef-f82dd4ce1fb9"),
	Query:      &query,
	Selectors: []product.Selector{
		{
			Characteristic: product.Characteristic{Key: "key1"},
			Values:         []string{"value1", "value2", "value3"},
		},
	},
	Ranges: []product.Range{
		{
			Characteristic: product.Characteristic{Key: "key2"},
			Min:            1,
			Max:            2,
		},
	},
}

func SingleCall(repo product.IRepository) (time.Duration, time.Duration) {
	start := time.Now()
	res, err := repo.GetWithFilter(filter, product.SORT_NONE)
	call := time.Since(start)

	if nil != err {
		panic(err)
	}

	start = time.Now()
	collection.Collect(res.Iter())
	serialization := time.Since(start)

	return call, serialization
}

func (self *TestProductRepository) Func(b *testing.B) {
	var totalCall time.Duration
	var totalSerialization time.Duration

	for range b.N {
		call, ser := SingleCall(self.repo)
		totalCall += call
		totalSerialization += ser
	}

	b.ReportMetric(
		float64(totalCall)/float64(b.N),
		metrics.CALL_PER_OP,
	)
	b.ReportMetric(
		float64(totalSerialization)/float64(b.N),
		metrics.SERIALIZATION_PER_OP,
	)
}

func (self *TestProductRepository) After() {
	if nil != self.clean {
		self.clean()
	}
}

