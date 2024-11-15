package review_test

import (
	"math"
	"math/rand/v2"
	"rent_service/builders/misc/collect"
	"slices"

	models_b "rent_service/builders/domain/models"
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"
	"rent_service/builders/mothers/test/repository/psql"
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/errors/cmnerrors"
	"rent_service/internal/repository/interfaces/review"
	"rent_service/misc/nullable"
	"rent_service/misc/testcommon"
	psqlcommon "rent_service/misc/testcommon/psql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func CompareReview(expected models.Review, actual models.Review) bool {
	return expected.Id == actual.Id &&
		expected.InstanceId == actual.InstanceId &&
		expected.UserId == actual.UserId &&
		expected.Content == actual.Content &&
		testcommon.EPSILON > math.Abs(expected.Rating-actual.Rating) &&
		psqlcommon.CompareTimeMicro(expected.Date, actual.Date)
}

type ReviewRepositoryTestSuite struct {
	suite.Suite
	repo review.IRepository
	psqlcommon.Context
}

func (self *ReviewRepositoryTestSuite) BeforeAll(t provider.T) {
	self.Context.SetUp(t)
}

func (self *ReviewRepositoryTestSuite) AfterAll(t provider.T) {
	self.Context.TearDown(t)
}

func (self *ReviewRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Review repository",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.Inserter.ClearDB()
	})

	t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
		self.repo = self.Factory.CreateReviewRepository()
	})
}

var describeCreate = testcommon.MethodDescriptor(
	"Create",
	"Create review",
)

var describeGetWithFilter = testcommon.MethodDescriptor(
	"GetWithFilter",
	"Get reviews with filter",
)

func CheckCreated(expected models.Review, actual models.Review) bool {
	return uuid.UUID{} != actual.Id &&
		expected.InstanceId == actual.InstanceId &&
		expected.UserId == actual.UserId &&
		expected.Content == actual.Content &&
		testcommon.EPSILON >= math.Abs(expected.Rating-actual.Rating) &&
		psqlcommon.CompareTimeMicro(expected.Date, actual.Date)
}

func (self *ReviewRepositoryTestSuite) TestCreatePositive(t provider.T) {
	var (
		user      models.User
		category  models.Category
		product   models.Product
		instance  models.Instance
		reference models.Review
	)

	describeCreate(t,
		"Simple create test",
		"Checks that no error is returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserDefault(nullable.None[string]()).Build(),
			)
			self.Inserter.InsertUser(&user)
		})

		t.WithNewStep("Create and insert category", func(sCtx provider.StepCtx) {
			category = testcommon.AssignParameter(sCtx, "category",
				models_om.CategoryRandomId().WithName("test").Build(),
			)
			self.Inserter.InsertCategory(&category)
		})

		t.WithNewStep("Create and insert product", func(sCtx provider.StepCtx) {
			product = testcommon.AssignParameter(sCtx, "product",
				models_om.ProductExmaple("1", category.Id).Build(),
			)
			self.Inserter.InsertProduct(&product)
		})

		t.WithNewStep("Create and insert instance", func(sCtx provider.StepCtx) {
			instance = testcommon.AssignParameter(sCtx, "instance",
				models_om.InstanceExample("1", product.Id).Build(),
			)
			self.Inserter.InsertInstance(&instance)
		})

		t.WithNewStep("Create review", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "review",
				models_om.ReviewExample(
					"Test",
					instance.Id,
					user.Id,
					nullable.None[float64](),
					nullable.None[time.Time](),
				).WithId(uuid.UUID{}).Build(),
			)
		})
	})

	// Act
	var result models.Review
	var err error

	t.WithNewStep("Create review", func(sCtx provider.StepCtx) {
		result, err = self.repo.Create(reference)
	}, allure.NewParameter("review", reference))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Assert[models.Review](t).EqualFunc(
		CheckCreated, reference, result, "Same review with non null uuid",
	)
}

func (self *ReviewRepositoryTestSuite) TestCreateDuplicate(t provider.T) {
	var (
		user      models.User
		category  models.Category
		product   models.Product
		instance  models.Instance
		reference models.Review
	)

	describeCreate(t,
		"Can't add second review",
		"Checks that error returned and mapped to Duplicate",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserDefault(nullable.None[string]()).Build(),
			)
			self.Inserter.InsertUser(&user)
		})

		t.WithNewStep("Create and insert category", func(sCtx provider.StepCtx) {
			category = testcommon.AssignParameter(sCtx, "category",
				models_om.CategoryRandomId().WithName("test").Build(),
			)
			self.Inserter.InsertCategory(&category)
		})

		t.WithNewStep("Create and insert product", func(sCtx provider.StepCtx) {
			product = testcommon.AssignParameter(sCtx, "product",
				models_om.ProductExmaple("1", category.Id).Build(),
			)
			self.Inserter.InsertProduct(&product)
		})

		t.WithNewStep("Create and insert instance", func(sCtx provider.StepCtx) {
			instance = testcommon.AssignParameter(sCtx, "instance",
				models_om.InstanceExample("1", product.Id).Build(),
			)
			self.Inserter.InsertInstance(&instance)
		})

		t.WithNewStep("Create review", func(sCtx provider.StepCtx) {
			builder := models_om.ReviewExample(
				"Test",
				instance.Id,
				user.Id,
				nullable.None[float64](),
				nullable.None[time.Time](),
			)
			created := testcommon.AssignParameter(sCtx, "created",
				builder.Build(),
			)
			reference = testcommon.AssignParameter(sCtx, "review",
				builder.WithId(uuid.UUID{}).Build(),
			)
			self.Inserter.InsertReview(&created)
		})
	})

	// Act
	var err error

	t.WithNewStep("Create delivery", func(sCtx provider.StepCtx) {
		_, err = self.repo.Create(reference)
	}, allure.NewParameter("delivery", reference))

	// Assert
	var nferr cmnerrors.ErrorDuplicate

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is Duplicate")
}

func (self *ReviewRepositoryTestSuite) TestGetWithFilterPositive(t provider.T) {
	var (
		users     []models.User
		category  models.Category
		product   models.Product
		instance  models.Instance
		all       []models.Review
		reference []models.Review
		filter    review.Filter
		sort      review.Sort
	)

	describeGetWithFilter(t,
		"Get middle reviews",
		"Checks that requested collection, sorted by rating is returned with no error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert users", func(sCtx provider.StepCtx) {
			users = testcommon.AssignParameter(sCtx, "users",
				collect.DoN(10, collect.FmtWrap(func(p string) *models_b.UserBuilder {
					return models_om.UserDefault(nullable.Some(p))
				})),
			)
			psql.BulkInsert(self.Inserter.InsertUser, users...)
		})

		t.WithNewStep("Create and insert category", func(sCtx provider.StepCtx) {
			category = testcommon.AssignParameter(sCtx, "category",
				models_om.CategoryRandomId().WithName("test").Build(),
			)
			self.Inserter.InsertCategory(&category)
		})

		t.WithNewStep("Create and insert product", func(sCtx provider.StepCtx) {
			product = testcommon.AssignParameter(sCtx, "product",
				models_om.ProductExmaple("1", category.Id).Build(),
			)
			self.Inserter.InsertProduct(&product)
		})

		t.WithNewStep("Create and insert instance", func(sCtx provider.StepCtx) {
			instance = testcommon.AssignParameter(sCtx, "instance",
				models_om.InstanceExample("1", product.Id).Build(),
			)
			self.Inserter.InsertInstance(&instance)
		})

		t.WithNewStep("Create and insert reviews", func(sCtx provider.StepCtx) {
			rnd := func(base float64) float64 {
				if v := base + rand.Float64() - 0.5; 0 > v {
					return 0
				} else if 5 < v {
					return 5
				} else {
					return v
				}
			}

			collection.ForEach(
				collection.EnumerateIterator(collection.SliceIterator(users)),
				func(pair *collection.Pair[uint, models.User]) {
					i := pair.A / 2
					review := models_om.ReviewExample(
						pair.B.Name,
						instance.Id,
						pair.B.Id,
						nullable.Some(rnd(float64(i))),
						nullable.None[time.Time](),
					).Build()

					if uint(1) <= i && i <= uint(4) {
						reference = append(reference, review)
					}

					all = append(all, review)
				},
			)

			slices.SortFunc(reference, func(a models.Review, b models.Review) int {
				if testcommon.EPSILON > math.Abs(a.Rating-b.Rating) {
					return 0
				} else if a.Rating > b.Rating {
					return -1
				} else {
					return 1
				}
			})

			sCtx.WithParameters(allure.NewParameter("reviews", all))
			psql.BulkInsert(self.Inserter.InsertReview, all...)
		})

		t.WithNewStep("Create filter", func(sCtx provider.StepCtx) {
			filter = review.Filter{
				InstanceId: instance.Id,
				Ratings:    []review.Rating{1, 2, 3, 4},
			}
			sort = review.SORT_RATING_DSC
		})
	})

	// Act
	var result collection.Collection[models.Review]
	var err error

	t.WithNewStep("Get reviews with filter", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetWithFilter(filter, sort)
	}, allure.NewParameter("form", filter), allure.NewParameter("sort", sort))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Assert[[]models.Review](t).EqualFunc(
		func(e []models.Review, a []models.Review) bool {
			return collection.All(
				collection.ZipIterator(
					collection.SliceIterator(e),
					collection.SliceIterator(a),
				),
				func(pair *collection.Pair[models.Review, models.Review]) bool {
					return CompareReview(pair.A, pair.B)
				},
			)
		},
		reference, collection.Collect(result.Iter()),
		"Same reviews",
	)
}

func (self *ReviewRepositoryTestSuite) TestGetWithFilterNotFound(t provider.T) {
	var (
		id   uuid.UUID
		form review.Filter
		sort review.Sort
	)

	describeGetWithFilter(t,
		"Instance not found",
		"Checks that method return error NotFound",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Generate unknonw id", func(sCtx provider.StepCtx) {
			id = testcommon.AssignParameter(sCtx, "id", uuidgen.Generate())
		})

		t.WithNewStep("Create form", func(sCtx provider.StepCtx) {
			form = review.Filter{
				InstanceId: id,
			}
			sort = review.SORT_NONE
		})
	})

	// Act
	var err error

	t.WithNewStep("Get delivery company by id", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetWithFilter(form, sort)
	}, allure.NewParameter("form", form), allure.NewParameter("sort", sort))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func TestReviewRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ReviewRepositoryTestSuite))
}

