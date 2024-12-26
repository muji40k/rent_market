package server_test

import (
	"context"
	"fmt"
	"net/http"
	models_b "rent_service/builders/domain/models"
	"rent_service/builders/misc/collect"
	"rent_service/builders/misc/generator"
	currency_om "rent_service/builders/mothers/currency"
	models_om "rent_service/builders/mothers/domain/models"
	requests_om "rent_service/builders/mothers/domain/requests"
	mserver "rent_service/builders/mothers/test/application/server"
	"rent_service/builders/mothers/test/repository/psql"
	"rent_service/builders/mothers/test/tracer"
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/requests"
	rv1t "rent_service/internal/factory/repositories/v1/tracer"
	sv1t "rent_service/internal/factory/services/v1/tracer"
	sinstance "rent_service/internal/logic/services/interfaces/instance"
	sphoto "rent_service/internal/logic/services/interfaces/photo"
	sprovide "rent_service/internal/logic/services/interfaces/provide"
	sstorage "rent_service/internal/logic/services/interfaces/storage"
	suser "rent_service/internal/logic/services/interfaces/user"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/misc/tracer/cleanstack"
	"rent_service/internal/misc/types/collection"
	"rent_service/misc/contextholder"
	"rent_service/misc/nullable"
	"rent_service/misc/testcommon"
	"rent_service/misc/testcommon/defservices"
	psqlcommon "rent_service/misc/testcommon/psql"
	"rent_service/misc/testcommon/server"
	"rent_service/server/authenticator"
	"rent_service/server/headers"
	"time"

	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	// "go.opentelemetry.io/otel"
	// "go.opentelemetry.io/otel/trace"
)

func MapProvisionRequest(value *requests.Provide, userId uuid.UUID) sprovide.ProvideRequest {
	return sprovide.ProvideRequest{
		Id:            value.Id,
		ProductId:     value.ProductId,
		UserId:        userId,
		PickUpPointId: value.PickUpPointId,
		Name:          value.Name,
		Description:   value.Description,
		Condition:     value.Condition,
		PayPlans: collection.Collect(collection.MapIterator(
			func(kv *collection.KV[uuid.UUID, models.PayPlan]) sprovide.PayPlan {
				return sprovide.PayPlan{
					Id:       kv.Value.Id,
					PeriodId: kv.Key,
					Price:    kv.Value.Price,
				}
			},
			collection.HashMapIterator(value.PayPlans),
		)),
		VerificationCode: value.VerificationCode,
		CreateDate:       date.New(value.CreateDate),
	}
}

type ServerE2ETestSuite struct {
	suite.Suite
	clener    *cleanstack.Cleaner
	seContext server.Context
	sContext  defservices.Context
	rContext  psqlcommon.Context
}

func (self *ServerE2ETestSuite) BeforeAll(t provider.T) {
	self.clener = cleanstack.New()
	holder := contextholder.New()
	provider := tracer.JaegerTracer(self.clener)
	tracer := provider.Tracer("rent_service")

	self.rContext.SetUp(t)

	rcw := rv1t.New(self.rContext.Factory, holder, tracer)
	self.sContext.SetUp(t, rcw.ToFactories())

	scw := sv1t.New(self.sContext.Factory, holder, tracer)
	self.seContext.SetUp(t, scw.ToFactories(),
		append(
			[]mserver.ServerExtender{
				mserver.TracerExtender(provider, holder),
			},
			mserver.DefaultControllers()...,
		)...,
	)
}

func (self *ServerE2ETestSuite) AfterAll(t provider.T) {
	self.seContext.TearDown(t)
	self.sContext.TearDown(t)
	self.rContext.TearDown(t)
	self.clener.Clean(context.Background())
}

func (self *ServerE2ETestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"End to end tests",
		"Server api tests",
		"MVP",
	)
}

func (self *ServerE2ETestSuite) TestStartProvision(t provider.T) {
	var (
		pup      models.PickUpPoint
		skUser   models.User
		category models.Category
		product  models.Product
		periods  []models.Period
		photos   [][]byte
	)

	t.Title("Start provision")
	t.Description(`
        New user registers as renter and create new provision request,
        then storekeeper accepts it. As a result new provision appears in
        renter's list
    `)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create photos", func(sCtx provider.StepCtx) {
			photos = collection.Collect(
				collection.MapIterator(
					func(_ *int) []byte {
						return models_om.ImagePNGContent(nullable.None[int]())
					},
					collection.RangeIterator(collection.RangeEnd(5)),
				),
			)
		})

		gg := generator.NewGeneratorGroup()

		pergen := psql.GeneratorStepList(t, "periods", &periods,
			func(i uint) (models.Period, uuid.UUID) {
				var v models.Period
				switch i % 6 {
				case 0:
					v = models_om.PeriodDay().Build()
				case 1:
					v = models_om.PeriodWeek().Build()
				case 2:
					v = models_om.PeriodMonth().Build()
				case 3:
					v = models_om.PeriodQuarter().Build()
				case 4:
					v = models_om.PeriodHalf().Build()
				case 5:
					v = models_om.PeriodYear().Build()
				}

				return v, v.Id
			},
			self.rContext.Inserter.InsertPeriod,
		)
		gg.Add(pergen, 6)

		pupgen := psql.GeneratorStepValue(t, "pick up point", &pup,
			func() (models.PickUpPoint, uuid.UUID) {
				v := models_om.PickUpPointExample("test").Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(pupgen)

		sugen := psql.GeneratorStepValue(t, "storekeeper user", &skUser,
			func() (models.User, uuid.UUID) {
				v := models_om.UserStorekeeper(nullable.None[string]()).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(sugen)

		sgen := psql.GeneratorStepNewValue(t, "storekeeper",
			func() (models.Storekeeper, uuid.UUID) {
				v := models_om.StorekeeperWithUserId(
					sugen.Generate(),
					pupgen.Generate(),
				).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertStorekeeper,
		)
		gg.Add(sgen, 1)

		cgen := psql.GeneratorStepValue(t, "category", &category,
			func() (models.Category, uuid.UUID) {
				v := models_om.CategoryRandomId().WithName("test").Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertCategory,
		)
		gg.AddFinish(cgen)

		pgen := psql.GeneratorStepValue(t, "product", &product,
			func() (models.Product, uuid.UUID) {
				v := models_om.ProductExmaple("test", cgen.Generate()).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertProduct,
		)
		gg.AddFinish(pgen)

		pcgen := psql.GeneratorStepNewValue(t, "product characteristics",
			func() (models.ProductCharacteristics, uuid.UUID) {
				v := models_om.ProductCharacteristics(
					pgen.Generate(),
					collect.Do(
						models_om.CharacteristicExample("key1", "value1"),
						models_om.CharacteristicExampleNumeric("key2", 1234),
					)...,
				).Build()
				return v, v.ProductId
			},
			self.rContext.Inserter.InsertProductCharacteristics,
		)
		gg.Add(pcgen, 1)

		gg.Generate()
		gg.Finish()
	})

	// Act -- Register renter user
	var user models.User
	t.WithNewStep("Register renter user", func(sCtx provider.StepCtx) {
		// Arrange
		sCtx.WithNewStep("Create reference user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserRenter(nullable.None[string]()).Build(),
			)
		})

		// Act
		sCtx.WithNewStep("Act", func(sCtx provider.StepCtx) {
			client := self.seContext.GetClient(sCtx)
			resp := client.POST("/users").WithJSON(map[string]interface{}{
				"name":     user.Name,
				"email":    user.Email,
				"password": user.Password,
			}).Expect()

			// Assert
			resp.Status(http.StatusCreated)
		})
	})

	// Act -- Login renter user
	var renterToken authenticator.ApiToken
	t.WithNewStep("Login renter user", func(sCtx provider.StepCtx) {
		// Act
		client := self.seContext.GetClient(sCtx)
		resp := client.POST("/sessions").WithJSON(map[string]interface{}{
			"email":    user.Email,
			"password": user.Password,
		}).Expect()

		// Assert
		resp.Status(http.StatusOK).
			JSON().Object().Decode(&renterToken)
	})

	// Act -- Get renter user info
	t.WithNewStep("Get renter user info", func(sCtx provider.StepCtx) {
		client := self.seContext.GetClient(sCtx)
		resp := client.GET("/users/self").
			WithHeader(headers.API_KEY, renterToken.Access).Expect()

		// Assert
		var info suser.Info
		resp.Status(http.StatusOK).
			JSON().Object().Decode(&info)
		sCtx.WithNewStep("Same user", func(sCtx provider.StepCtx) {
			sCtx.Assert().NotZero(info.Id, "Non null id")
			sCtx.Assert().Equal(user.Email, info.Email, "Same email")
			sCtx.Assert().Equal(user.Name, info.Name, "Same name")
		})
		user.Id = info.Id
	})

	// Act -- Register As renter
	t.WithNewStep("Register As renter", func(sCtx provider.StepCtx) {
		client := self.seContext.GetClient(sCtx)
		resp := client.POST("/renters/self").
			WithHeader(headers.API_KEY, renterToken.Access).Expect()

		// Assert
		resp.Status(http.StatusCreated)
	})

	// Act -- Check renter's record
	t.WithNewStep("Check renter's record", func(sCtx provider.StepCtx) {
		client := self.seContext.GetClient(sCtx)
		resp := client.GET("/renters/self").
			WithHeader(headers.API_KEY, renterToken.Access).Expect()

		// Assert
		resp.Status(http.StatusOK)
	})

	// Act -- Create provision request
	var rq sprovide.ProvideRequest
	t.WithNewStep("Create provision request", func(sCtx provider.StepCtx) {
		// Arrange
		var createForm map[string]interface{}
		var refRQ sprovide.ProvideRequest

		sCtx.WithNewStep("Create form", func(sCtx provider.StepCtx) {
			mrq := requests_om.ProvideExample(
				"e2e_test", product.Id, uuid.UUID{}, pup.Id,
				nullable.Some(""), nullable.Some(time.Time{}),
				collect.DoN(6, func(i uint) *models_b.PayPlanBuilder {
					return models_om.PayPlanWithPeriodIdAndPrice(
						periods[i].Id,
						currency_om.RUB(float64(100*(i+1))).Build(),
					)
				})...,
			).Build()
			refRQ = MapProvisionRequest(&mrq, user.Id)
			createForm = map[string]interface{}{
				"product":       product.Id,
				"pick_up_point": pup.Id,
				"name":          mrq.Name,
				"description":   mrq.Description,
				"condition":     mrq.Condition,
				"pay_plans": collection.Collect(
					collection.MapIterator(
						func(kv *collection.KV[uuid.UUID, models.PayPlan]) map[string]interface{} {
							return map[string]interface{}{
								"period": kv.Key,
								"price":  kv.Value.Price,
							}
						},
						collection.HashMapIterator(mrq.PayPlans),
					),
				),
			}
		})

		// Act
		sCtx.WithNewStep("Act", func(sCtx provider.StepCtx) {
			client := self.seContext.GetClient(sCtx)
			resp := client.POST("/provision-requests").
				WithHeader(headers.API_KEY, renterToken.Access).
				WithJSON(createForm).
				Expect()

			// Assert
			resp.Status(http.StatusCreated).JSON().Object().Decode(&rq)
			testcommon.Assert[sprovide.ProvideRequest](sCtx).EqualFunc(
				func(e sprovide.ProvideRequest, a sprovide.ProvideRequest) bool {
					return uuid.UUID{} != a.Id &&
						e.ProductId == a.ProductId &&
						e.UserId == a.UserId &&
						e.PickUpPointId == a.PickUpPointId &&
						e.Name == a.Name &&
						e.Description == a.Description &&
						e.Condition == a.Condition &&
						"" != a.VerificationCode &&
						date.Date{} != a.CreateDate &&
						testcommon.ElementsMatchFunc(
							func(e sprovide.PayPlan, a sprovide.PayPlan) bool {
								return uuid.UUID{} != a.Id &&
									e.PeriodId == a.PeriodId &&
									testcommon.CompareCurrency(e.Price, a.Price)
							}, e.PayPlans, a.PayPlans,
						)
				}, refRQ, rq, "Right request")
		})
	})

	// Act -- Login storekeeper user
	var storekeeperToken authenticator.ApiToken
	t.WithNewStep("Login storekeeper user", func(sCtx provider.StepCtx) {
		client := self.seContext.GetClient(sCtx)
		resp := client.POST("/sessions").WithJSON(map[string]interface{}{
			"email":    skUser.Email,
			"password": skUser.Password,
		}).Expect()

		// Assert
		resp.Status(http.StatusOK).
			JSON().Object().Decode(&storekeeperToken)
	})

	// Act -- Display provision request
	t.WithNewStep("Display provision request", func(sCtx provider.StepCtx) {
		client := self.seContext.GetClient(sCtx)
		resp := client.GET("/provision-requests").
			WithHeader(headers.API_KEY, storekeeperToken.Access).
			WithQuery("pickUpPointId", pup.Id).
			Expect()

		// Assert
		var content []sprovide.ProvideRequest
		resp.Status(http.StatusOK).JSON().Decode(&content)
		sCtx.Assert().Equal(1, len(content), "One request in collection")
		testcommon.Assert[sprovide.ProvideRequest](sCtx).EqualFunc(
			func(e sprovide.ProvideRequest, a sprovide.ProvideRequest) bool {
				return e.Id == a.Id &&
					e.ProductId == a.ProductId &&
					e.UserId == a.UserId &&
					e.PickUpPointId == a.PickUpPointId &&
					e.Name == a.Name &&
					e.Description == a.Description &&
					e.Condition == a.Condition &&
					e.VerificationCode == a.VerificationCode &&
					psqlcommon.CompareTimeMicro(e.CreateDate.Time, a.CreateDate.Time) &&
					testcommon.ElementsMatchFunc(
						func(e sprovide.PayPlan, a sprovide.PayPlan) bool {
							return e.Id == a.Id &&
								e.PeriodId == a.PeriodId &&
								testcommon.CompareCurrency(e.Price, a.Price)
						}, e.PayPlans, a.PayPlans,
					)
			}, rq, content[0], "Same provide request")
	})

	// Act -- Create and load temp photos of new instatnce
	var photoIds []uuid.UUID
	t.WithNewStep("Create and load temp photos of new instatnce", func(sCtx provider.StepCtx) {
		client := self.seContext.GetClient(sCtx)

		photoIds = collection.Collect(
			collection.MapIterator(
				func(pair *collection.Pair[uint, []byte]) uuid.UUID {
					resp := client.POST("/photos/temp").
						WithHeader(headers.API_KEY, storekeeperToken.Access).
						WithJSON(map[string]interface{}{
							"mime":        "image/png",
							"placeholder": fmt.Sprintf("placeholder for %v", pair.A),
							"description": fmt.Sprintf("placeholder for %v", pair.A),
						}).
						Expect()

					// Assert
					var id uuid.UUID
					resp.Status(http.StatusCreated).JSON().Decode(&id)

					resp = client.POST("/photos/temp/{id}", id).
						WithHeader(headers.API_KEY, storekeeperToken.Access).
						WithHeader("Content-Type", "image/png").
						WithBytes(pair.B).
						Expect()

					// Assert
					resp.Status(http.StatusCreated)

					return id
				},
				collection.EnumerateIterator(collection.SliceIterator(photos)),
			),
		)
	})

	// Act -- Accept provision request
	var updatedDescription = "Updated description from accept stage"
	t.WithNewStep("Accept provision request", func(sCtx provider.StepCtx) {
		client := self.seContext.GetClient(sCtx)
		resp := client.PUT("/provision-requests/{id}", rq.Id).
			WithHeader(headers.API_KEY, storekeeperToken.Access).
			WithJSON(map[string]interface{}{
				"action": "satisfy",
				"overrides": map[string]interface{}{
					"description": updatedDescription,
				},
				"state_photos":      photoIds,
				"verification_code": rq.VerificationCode,
			}).Expect()

		// Assert
		resp.Status(http.StatusOK)
	})

	// Act -- Display storages for pick up point
	var instanceId uuid.UUID
	t.WithNewStep("Display storages for pick up point", func(sCtx provider.StepCtx) {
		client := self.seContext.GetClient(sCtx)
		resp := client.GET("/stored-instances").
			WithHeader(headers.API_KEY, storekeeperToken.Access).
			WithQuery("pickUpPointId", pup.Id).
			Expect()

		// Assert
		var content []sstorage.Storage
		resp.Status(http.StatusOK).JSON().Decode(&content)
		sCtx.Assert().Equal(1, len(content), "Single element in collection")
		sCtx.WithNewStep("Correct storage", func(sCtx provider.StepCtx) {
			sCtx.Assert().NotEmpty(content[0].Id, "Non null id")
			sCtx.Assert().Equal(content[0].PickUpPointId, pup.Id, "Same pick up point")
			sCtx.Assert().NotEmpty(content[0].InstanceId, "Non null instance id")
			sCtx.Assert().NotEmpty(content[0].InDate, "Non null in date")
			sCtx.Assert().Nil(content[0].OutDate, "Nil out date")
		})
		instanceId = content[0].InstanceId
	})

	// Act -- Display provision requests for renter
	t.WithNewStep("Display provision requests for renter", func(sCtx provider.StepCtx) {
		client := self.seContext.GetClient(sCtx)
		resp := client.GET("/provision-requests").
			WithHeader(headers.API_KEY, renterToken.Access).
			WithQuery("userId", user.Id).
			Expect()

		// Assert
		var content []sprovide.ProvideRequest
		resp.Status(http.StatusOK).JSON().Decode(&content)
		sCtx.Assert().Equal(0, len(content), "Collection is empty")
	})

	// Act -- Display provisions for renter
	t.WithNewStep("Display provisions for renter", func(sCtx provider.StepCtx) {
		client := self.seContext.GetClient(sCtx)
		resp := client.GET("/provisions").
			WithHeader(headers.API_KEY, renterToken.Access).
			WithQuery("userId", user.Id).
			Expect()

		// Assert
		var content []sprovide.Provision
		resp.Status(http.StatusOK).JSON().Decode(&content)
		sCtx.Assert().Equal(1, len(content), "Single element in collection")
		sCtx.WithNewStep("Correct provision", func(sCtx provider.StepCtx) {
			sCtx.Assert().NotEmpty(content[0].Id, "Non null id")
			sCtx.Assert().Equal(content[0].UserId, user.Id, "Same pick up point")
			sCtx.Assert().Equal(instanceId, content[0].InstanceId, "Same instance id as in storage")
			sCtx.Assert().NotEmpty(content[0].StartDate, "Non null start date")
			sCtx.Assert().Nil(content[0].EndDate, "Nil end date")
		})
	})

	// Act -- Display intstance
	t.WithNewStep("Display instance", func(sCtx provider.StepCtx) {
		client := self.seContext.GetClient(sCtx)
		resp := client.GET("/instances/{id}", instanceId).Expect()

		// Assert
		var content sinstance.Instance
		resp.Status(http.StatusOK).JSON().Decode(&content)
		sCtx.WithNewStep("Correct instance with updated description", func(sCtx provider.StepCtx) {
			sCtx.Assert().Equal(instanceId, content.Id, "Same id")
			sCtx.Assert().Equal(product.Id, content.ProductId, "Same product id")
			sCtx.Assert().Equal(rq.Name, content.Name, "Same name")
			sCtx.Assert().Equal(updatedDescription, content.Description, "Updated description")
			sCtx.Assert().Equal(rq.Condition, content.Condition, "Same condition")
		})
	})

	// Act -- Display intstance pay plans
	t.WithNewStep("Display instance pay plans", func(sCtx provider.StepCtx) {
		client := self.seContext.GetClient(sCtx)
		resp := client.GET("/instances/{id}/pay-plans", instanceId).Expect()

		// Assert
		var content []sinstance.PayPlan
		resp.Status(http.StatusOK).JSON().Decode(&content)
		sCtx.Assert().Equal(len(rq.PayPlans), len(content), "Same length of collection")
		sCtx.Assert().True(testcommon.ElementsMatchFunc(
			func(e sprovide.PayPlan, a sinstance.PayPlan) bool {
				return e.PeriodId == a.PeriodId &&
					instanceId == a.InstanceId &&
					testcommon.CompareCurrency(e.Price, a.Price)
			}, rq.PayPlans, content,
		), "Correct pay plans")
	})

	// Act -- Display intstance photos
	t.WithNewStep("Display instance photos", func(sCtx provider.StepCtx) {
		client := self.seContext.GetClient(sCtx)
		resp := client.GET("/instances/{id}/photos", instanceId).Expect()

		// Assert
		var content []uuid.UUID
		resp.Status(http.StatusOK).JSON().Decode(&content)
		sCtx.Assert().Equal(len(photoIds), len(content), "Same length of collection")

		sCtx.WithNewStep("Correct photos", func(sCtx provider.StepCtx) {
			collection.ForEach(
				collection.SliceIterator(content),
				func(id *uuid.UUID) {
					// Act
					resp := client.GET("/photos/{id}", *id).Expect()

					// Assert
					var content sphoto.Photo
					resp.Status(http.StatusOK).JSON().Decode(&content)

					if *id != content.Id {
						sCtx.Errorf("Photo id not match:\nExpected: %v\nGot: %v",
							*id, content.Id,
						)
					}

					if "image/png" != content.Mime {
						sCtx.Errorf("Photo mime isn't 'image/png'. Got: %v", content.Mime)
					}

					if "" == content.Placeholder {
						sCtx.Errorf("Got empty placeholder")
					}

					if "" == content.Description {
						sCtx.Errorf("Got empty description")
					}

					if "" == content.Href {
						sCtx.Errorf("Got empty href")
					}

					if d := (date.Date{}); d == content.Date {
						sCtx.Errorf("Got empty date")
					}
				},
			)
		})
	})

	// Act -- Logout renter user
	t.WithNewStep("Logout renter user", func(sCtx provider.StepCtx) {
		client := self.seContext.GetClient(sCtx)
		resp := client.DELETE("/sessions").
			WithHeader(headers.API_KEY, renterToken.Access).Expect()

		// Assert
		resp.Status(http.StatusOK)
	})

	// Act -- Logout storekeepr user
	t.WithNewStep("Logout storekeeper user", func(sCtx provider.StepCtx) {
		client := self.seContext.GetClient(sCtx)
		resp := client.DELETE("/sessions").
			WithHeader(headers.API_KEY, storekeeperToken.Access).Expect()

		// Assert
		resp.Status(http.StatusOK)
	})
}

func TestServerE2ETestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ServerE2ETestSuite))
}

