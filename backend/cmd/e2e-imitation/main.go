package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	models_b "rent_service/builders/domain/models"
	"rent_service/builders/misc/collect"
	currency_om "rent_service/builders/mothers/currency"
	models_om "rent_service/builders/mothers/domain/models"
	requests_om "rent_service/builders/mothers/domain/requests"
	"rent_service/builders/mothers/test/application/server"
	"rent_service/builders/mothers/test/repository/psql"
	"rent_service/builders/mothers/test/service/defservices"
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/requests"
	sinstance "rent_service/internal/logic/services/interfaces/instance"
	sphoto "rent_service/internal/logic/services/interfaces/photo"
	sprovide "rent_service/internal/logic/services/interfaces/provide"
	sstorage "rent_service/internal/logic/services/interfaces/storage"
	suser "rent_service/internal/logic/services/interfaces/user"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/misc/types/collection"
	"rent_service/misc/nullable"
	"rent_service/server/authenticator"
	"rent_service/server/headers"
	"time"

	"github.com/google/uuid"
)

func main() {
	ctx := Context{}
	ctx.BeforeAll()
	ctx.Before()

	host := os.Getenv("TARGET_HOST")

	if "" == host {
		host = "localhost:12345"
	}

	curl := url.URL{
		Scheme: "http",
		Host:   host,
	}

	ctx.TestStartProvision(&curl)

	ctx.AfterAll()
}

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

type Context struct {
	AuthInserter  *server.Inserter
	PhotoRegistry *defservices.PhotoRegistry
	RepoInserter  *psql.Inserter
}

func (self *Context) BeforeAll() {
	self.RepoInserter = psql.NewInserter()
	self.PhotoRegistry = defservices.NewPhotoRegistry()
	self.AuthInserter = server.NewInserter()
}

func (self *Context) AfterAll() {
	if nil != self.AuthInserter {
		self.AuthInserter.ClearDB()
		self.AuthInserter.Close()
	}

	if nil != self.PhotoRegistry {
		self.PhotoRegistry.Clear()
	}

	if nil != self.RepoInserter {
		self.RepoInserter.ClearDB()
		self.RepoInserter.Close()
	}
}

func (self *Context) Before() {
	self.RepoInserter.ClearDB()
	self.AuthInserter.ClearDB()
	self.PhotoRegistry.Clear()
}

func Hider[T any](value T) *T {
	return &value
}

func CallWrap(f func() (*http.Response, error)) *http.Response {
	if r, err := f(); nil == err {
		return r
	} else {
		panic(err)
	}
}

func DefaultCallWrap(req *http.Request) *http.Response {
	return CallWrap(func() (*http.Response, error) {
		return http.DefaultClient.Do(req)
	})
}

func DefaultNewRequest(method string, url *url.URL, body io.Reader) *http.Request {
	if rq, err := http.NewRequest(method, url.String(), body); nil == err {
		if nil != body {
			rq.Header.Add("Content-Type", "application/json")
		}

		return rq
	} else {
		panic(err)
	}
}

func PanicOnWrongStatus(resp *http.Response, target int) {
	if resp.StatusCode != target {
		panic(fmt.Errorf("Unexpected status code[%v]: %v", target, resp.StatusCode))
	}
}

func Marshal(value any) io.Reader {
	content, err := json.Marshal(value)

	if nil == err {
		return bytes.NewBuffer(content)
	} else {
		panic(err)
	}
}

func Unmarshal(resp *http.Response, value any) {
	content, err := io.ReadAll(resp.Body)

	if nil == err {
		err = json.Unmarshal(content, value)
	}

	if nil != err {
		panic(err)
	}
}

func InserterWrap[T any](inserter func(*T), iterator collection.Iterator[T]) collection.Iterator[T] {
	return collection.MapIterator(
		func(value *T) T {
			inserter(value)
			return *value
		},
		iterator,
	)
}

type bytesWrap struct {
	buffer *bytes.Buffer
}

func (self *bytesWrap) Read(p []byte) (n int, err error) {
	return self.buffer.Read(p)
}

func (self *bytesWrap) Close() error {
	return nil
}

func Step(builder func() *http.Request, do func(*http.Response)) {
	rq := builder()
	fmt.Printf("REQUEST: %v %v\n", rq.Method, rq.URL.String())
	resp := DefaultCallWrap(rq)
	fmt.Printf("RESPONSE: %v\n", resp.Status)

	if nil != resp.Body {
		content, err := io.ReadAll(resp.Body)

		if nil != err {
			panic(err)
		}

		resp.Body.Close()
		resp.Body = &bytesWrap{bytes.NewBuffer(content)}

		if 0 != len(content) {
			pp, err := json.MarshalIndent(json.RawMessage(content), "", "\t")

			if nil != err {
				panic(err)
			}

			fmt.Printf("CONTENT:\n%v\n", string(pp))
		}
	}

	do(resp)
	resp.Body.Close()
}

func SetQuery(url *url.URL, f func(*url.Values)) {
	values := url.Query()
	f(&values)
	url.RawQuery = values.Encode()
}

func MapPeriod(i int) *models_b.PeriodBuilder {
	switch i % 6 {
	case 1:
		return models_om.PeriodWeek()
	case 2:
		return models_om.PeriodMonth()
	case 3:
		return models_om.PeriodQuarter()
	case 4:
		return models_om.PeriodHalf()
	case 5:
		return models_om.PeriodYear()
	default:
		return models_om.PeriodDay()
	}
}

func (self *Context) TestStartProvision(base *url.URL) {
	// Arrange
	photos := collection.Collect(collection.MapIterator(
		func(_ *int) []byte {
			return models_om.ImagePNGContent(nullable.None[int]())
		},
		collection.RangeIterator(collection.RangeEnd(5)),
	))

	periods := collection.Collect(InserterWrap(
		self.RepoInserter.InsertPeriod,
		collection.MapIterator(
			func(i *int) models.Period { return MapPeriod(*i).Build() },
			collection.RangeIterator(collection.RangeEnd(6)),
		),
	))

	pup := models_om.PickUpPointExample("test").Build()
	self.RepoInserter.InsertPickUpPoint(&pup)

	skUser := models_om.UserStorekeeper(nullable.None[string]()).Build()
	self.RepoInserter.InsertUser(&skUser)

	self.RepoInserter.InsertStorekeeper(Hider(
		models_om.StorekeeperWithUserId(skUser.Id, pup.Id).Build(),
	))

	category := models_om.CategoryRandomId().WithName("test").Build()
	self.RepoInserter.InsertCategory(&category)

	product := models_om.ProductExmaple("test", category.Id).Build()
	self.RepoInserter.InsertProduct(&product)

	self.RepoInserter.InsertProductCharacteristics(Hider(
		models_om.ProductCharacteristics(
			product.Id,
			collect.Do(
				models_om.CharacteristicExample("key1", "value1"),
				models_om.CharacteristicExampleNumeric("key2", 1234),
			)...,
		).Build(),
	))

	// Act -- Register renter user
	user := models_om.UserRenter(nullable.None[string]()).Build()
	Step(
		func() *http.Request {
			url := base.JoinPath("users")
			return DefaultNewRequest("POST", url,
				Marshal(map[string]interface{}{
					"name":     user.Name,
					"email":    user.Email,
					"password": user.Password,
				}),
			)
		},
		func(resp *http.Response) {
			PanicOnWrongStatus(resp, http.StatusCreated)
		},
	)

	// Act -- Login renter user
	var renterToken authenticator.ApiToken
	Step(
		func() *http.Request {
			url := base.JoinPath("sessions")
			return DefaultNewRequest("POST", url,
				Marshal(map[string]interface{}{
					"email":    user.Email,
					"password": user.Password,
				}),
			)
		},
		func(resp *http.Response) {
			PanicOnWrongStatus(resp, http.StatusOK)
			Unmarshal(resp, &renterToken)
		},
	)

	// Act -- Get renter user info
	Step(
		func() *http.Request {
			url := base.JoinPath("users", "self")
			rq := DefaultNewRequest("GET", url, nil)
			rq.Header.Add(headers.API_KEY, renterToken.Access)
			return rq
		},
		func(resp *http.Response) {
			var info suser.Info
			PanicOnWrongStatus(resp, http.StatusOK)
			Unmarshal(resp, &info)
			user.Id = info.Id
		},
	)

	// Act -- Register As renter
	Step(
		func() *http.Request {
			url := base.JoinPath("renters", "self")
			rq := DefaultNewRequest("POST", url, nil)
			rq.Header.Add(headers.API_KEY, renterToken.Access)
			return rq
		},
		func(resp *http.Response) {
			PanicOnWrongStatus(resp, http.StatusCreated)
		},
	)

	// Act -- Check renter's record
	Step(
		func() *http.Request {
			url := base.JoinPath("renters", "self")
			rq := DefaultNewRequest("GET", url, nil)
			rq.Header.Add(headers.API_KEY, renterToken.Access)
			return rq
		},
		func(resp *http.Response) {
			PanicOnWrongStatus(resp, http.StatusOK)
		},
	)

	// Act -- Create provision request
	var rq sprovide.ProvideRequest

	Step(
		func() *http.Request {
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
			createForm := map[string]interface{}{
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

			url := base.JoinPath("provision-requests")
			rq := DefaultNewRequest("POST", url, Marshal(createForm))
			rq.Header.Add(headers.API_KEY, renterToken.Access)
			return rq
		},
		func(resp *http.Response) {
			PanicOnWrongStatus(resp, http.StatusCreated)
			Unmarshal(resp, &rq)
		},
	)

	// Act -- Login storekeeper user
	var storekeeperToken authenticator.ApiToken

	Step(
		func() *http.Request {
			url := base.JoinPath("sessions")
			return DefaultNewRequest("POST", url,
				Marshal(map[string]interface{}{
					"email":    skUser.Email,
					"password": skUser.Password,
				}),
			)
		},
		func(resp *http.Response) {
			PanicOnWrongStatus(resp, http.StatusOK)
			Unmarshal(resp, &storekeeperToken)
		},
	)

	// Act -- Display provision request
	Step(
		func() *http.Request {
			curl := base.JoinPath("provision-requests")
			SetQuery(curl, func(v *url.Values) {
				v.Add("pickUpPointId", pup.Id.String())
			})
			rq := DefaultNewRequest("GET", curl, nil)
			rq.Header.Add(headers.API_KEY, storekeeperToken.Access)
			return rq
		},
		func(resp *http.Response) {
			var content []sprovide.ProvideRequest
			PanicOnWrongStatus(resp, http.StatusOK)
			Unmarshal(resp, &content)
		},
	)

	// Act -- Create and load temp photos of new instatnce
	photoIds := collection.Collect(
		collection.MapIterator(
			func(pair *collection.Pair[uint, []byte]) uuid.UUID {
				var id uuid.UUID

				Step(
					func() *http.Request {
						curl := base.JoinPath("photos", "temp")
						rq := DefaultNewRequest("POST", curl, Marshal(
							map[string]interface{}{
								"mime":        "image/png",
								"placeholder": fmt.Sprintf("placeholder for %v", pair.A),
								"description": fmt.Sprintf("placeholder for %v", pair.A),
							},
						))
						rq.Header.Add(headers.API_KEY, storekeeperToken.Access)
						return rq
					},
					func(resp *http.Response) {
						PanicOnWrongStatus(resp, http.StatusCreated)
						Unmarshal(resp, &id)
					},
				)

				Step(
					func() *http.Request {
						url := base.JoinPath("photos", "temp", id.String())
						rq, err := http.NewRequest(
							"POST",
							url.String(),
							bytes.NewBuffer(pair.B),
						)

						if nil != err {
							panic(err)
						}

						rq.Header.Add("Content-Type", "image/png")
						rq.Header.Add(headers.API_KEY, storekeeperToken.Access)

						return rq
					},
					func(resp *http.Response) {
						PanicOnWrongStatus(resp, http.StatusCreated)
					},
				)

				return id
			},
			collection.EnumerateIterator(collection.SliceIterator(photos)),
		),
	)

	// Act -- Accept provision request
	const updatedDescription = "Updated description from accept stage"

	Step(
		func() *http.Request {
			url := base.JoinPath("provision-requests", rq.Id.String())
			rq := DefaultNewRequest("PUT", url, Marshal(map[string]interface{}{
				"action": "satisfy",
				"overrides": map[string]interface{}{
					"description": updatedDescription,
				},
				"state_photos":      photoIds,
				"verification_code": rq.VerificationCode,
			}))
			rq.Header.Add(headers.API_KEY, storekeeperToken.Access)
			return rq
		},
		func(resp *http.Response) {
			PanicOnWrongStatus(resp, http.StatusOK)
		},
	)

	// Act -- Display storages for pick up point
	var instanceId uuid.UUID

	Step(
		func() *http.Request {
			curl := base.JoinPath("stored-instances")
			SetQuery(curl, func(v *url.Values) {
				v.Add("pickUpPointId", pup.Id.String())
			})
			rq := DefaultNewRequest("GET", curl, nil)
			rq.Header.Add(headers.API_KEY, storekeeperToken.Access)
			return rq
		},
		func(resp *http.Response) {
			var content []sstorage.Storage
			PanicOnWrongStatus(resp, http.StatusOK)
			Unmarshal(resp, &content)
			instanceId = content[0].InstanceId
		},
	)

	// Act -- Display provision requests for renter
	Step(
		func() *http.Request {
			curl := base.JoinPath("provision-requests")
			SetQuery(curl, func(v *url.Values) {
				v.Add("userId", user.Id.String())
			})
			rq := DefaultNewRequest("GET", curl, nil)
			rq.Header.Add(headers.API_KEY, renterToken.Access)
			return rq
		},
		func(resp *http.Response) {
			var content []sprovide.ProvideRequest
			PanicOnWrongStatus(resp, http.StatusOK)
			Unmarshal(resp, &content)
			if 0 != len(content) {
				panic("Non empty collection")
			}
		},
	)

	// Act -- Display provisions for renter
	Step(
		func() *http.Request {
			curl := base.JoinPath("provisions")
			SetQuery(curl, func(v *url.Values) {
				v.Add("userId", user.Id.String())
			})
			rq := DefaultNewRequest("GET", curl, nil)
			rq.Header.Add(headers.API_KEY, renterToken.Access)
			return rq
		},
		func(resp *http.Response) {
			var content []sprovide.Provision
			PanicOnWrongStatus(resp, http.StatusOK)
			Unmarshal(resp, &content)
			if 1 != len(content) {
				panic("Incorrect collection (expected length is 1)")
			}
		},
	)

	// Act -- Display intstance
	Step(
		func() *http.Request {
			curl := base.JoinPath("instances", instanceId.String())
			return DefaultNewRequest("GET", curl, nil)
		},
		func(resp *http.Response) {
			var content sinstance.Instance
			PanicOnWrongStatus(resp, http.StatusOK)
			Unmarshal(resp, &content)
		},
	)

	// Act -- Display intstance pay plans
	Step(
		func() *http.Request {
			curl := base.JoinPath("instances", instanceId.String(), "pay-plans")
			return DefaultNewRequest("GET", curl, nil)
		},
		func(resp *http.Response) {
			var content []sinstance.PayPlan
			PanicOnWrongStatus(resp, http.StatusOK)
			Unmarshal(resp, &content)
			if 6 != len(content) {
				panic("Incorrect collection (expected length is 6)")
			}
		},
	)

	// Act -- Display intstance photos
	var ids []uuid.UUID

	Step(
		func() *http.Request {
			curl := base.JoinPath("instances", instanceId.String(), "photos")
			return DefaultNewRequest("GET", curl, nil)
		},
		func(resp *http.Response) {
			PanicOnWrongStatus(resp, http.StatusOK)
			Unmarshal(resp, &ids)
			if 5 != len(ids) {
				panic("Incorrect collection (expected length is 5)")
			}
		},
	)

	collection.ForEach(
		collection.SliceIterator(ids),
		func(id *uuid.UUID) {
			Step(
				func() *http.Request {
					curl := base.JoinPath("photos", id.String())
					return DefaultNewRequest("GET", curl, nil)
				},
				func(resp *http.Response) {
					var content sphoto.Photo
					PanicOnWrongStatus(resp, http.StatusOK)
					Unmarshal(resp, &content)
				},
			)
		},
	)

	// Act -- Logout renter user
	Step(
		func() *http.Request {
			curl := base.JoinPath("sessions")
			rq := DefaultNewRequest("DELETE", curl, nil)
			rq.Header.Add(headers.API_KEY, renterToken.Access)
			return rq
		},
		func(resp *http.Response) {
			PanicOnWrongStatus(resp, http.StatusOK)
		},
	)

	// Act -- Logout storekeepr user
	Step(
		func() *http.Request {
			curl := base.JoinPath("sessions")
			rq := DefaultNewRequest("DELETE", curl, nil)
			rq.Header.Add(headers.API_KEY, storekeeperToken.Access)
			return rq
		},
		func(resp *http.Response) {
			PanicOnWrongStatus(resp, http.StatusOK)
		},
	)
}

