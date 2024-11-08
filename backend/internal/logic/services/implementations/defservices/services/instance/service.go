package instance

import (
	"errors"
	"math"
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/records"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/access"
	"rent_service/internal/logic/services/implementations/defservices/authenticator"
	"rent_service/internal/logic/services/implementations/defservices/emptymathcer"
	"rent_service/internal/logic/services/implementations/defservices/photoregistry"
	"rent_service/internal/logic/services/interfaces/instance"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/logic/services/types/token"
	. "rent_service/internal/misc/types/collection"
	instance_provider "rent_service/internal/repository/context/providers/instance"
	rent_provider "rent_service/internal/repository/context/providers/rent"
	review_provider "rent_service/internal/repository/context/providers/review"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
	instance_repository "rent_service/internal/repository/interfaces/instance"
	review_repository "rent_service/internal/repository/interfaces/review"
	"rent_service/misc/mapfuncs"
	"time"

	"github.com/google/uuid"
)

type repoproviders struct {
	instance instance_provider.IProvider
}

type accessor struct {
	instance access.IInstance
}

type service struct {
	repos         repoproviders
	access        accessor
	authenticator authenticator.IAuthenticator
}

func New(
	authenticator authenticator.IAuthenticator,
	instance instance_provider.IProvider,
	instanceacc access.IInstance,
) instance.IService {
	return &service{
		repoproviders{instance},
		accessor{instanceacc},
		authenticator,
	}
}

func mapf(value *models.Instance) instance.Instance {
	return instance.Instance{
		Id:          value.Id,
		ProductId:   value.ProductId,
		Name:        value.Name,
		Description: value.Description,
		Condition:   value.Condition,
	}
}

func unmapf(value *instance.Instance) models.Instance {
	return models.Instance{
		Id:          value.Id,
		ProductId:   value.ProductId,
		Name:        value.Name,
		Description: value.Description,
		Condition:   value.Condition,
	}
}

func mapFilter(value *instance.Filter) instance_repository.Filter {
	return instance_repository.Filter{
		ProductId: value.ProductId,
	}
}

func mapSort(value *instance.Sort) (instance_repository.Sort, error) {
	switch *value {
	case instance.SORT_NONE:
		return instance_repository.SORT_NONE, nil
	case instance.SORT_RATING_ASC:
		return instance_repository.SORT_RATING_ASC, nil
	case instance.SORT_RATING_DSC:
		return instance_repository.SORT_RATING_DSC, nil
	case instance.SORT_DATE_ASC:
		return instance_repository.SORT_DATE_ASC, nil
	case instance.SORT_DATE_DSC:
		return instance_repository.SORT_DATE_DSC, nil
	case instance.SORT_PRICE_ASC:
		return instance_repository.SORT_PRICE_ASC, nil
	case instance.SORT_PRICE_DSC:
		return instance_repository.SORT_PRICE_DSC, nil
	case instance.SORT_USAGE_ASC:
		return instance_repository.SORT_USAGE_ASC, nil
	case instance.SORT_USAGE_DSC:
		return instance_repository.SORT_USAGE_DSC, nil
	default:
		return instance_repository.SORT_NONE, cmnerrors.Unknown("sort")
	}

}

func (self *service) ListInstances(
	filter instance.Filter,
	sort instance.Sort,
) (Collection[instance.Instance], error) {
	var instance Collection[models.Instance]
	sortr, err := mapSort(&sort)

	if nil == err {
		repo := self.repos.instance.GetInstanceRepository()
		instance, err = repo.GetWithFilter(mapFilter(&filter), sortr)
	}

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return MapCollection(mapf, instance), err
}

func (self *service) GetInstanceById(
	instanceId uuid.UUID,
) (instance.Instance, error) {
	repo := self.repos.instance.GetInstanceRepository()
	instance, err := repo.GetById(instanceId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return mapf(&instance), err
}

func (self *service) UpdateInstance(
	token token.Token,
	instance instance.Instance,
) error {
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.instance.Access(user.Id, instance.Id)
	}

	if nil == err {
		repo := self.repos.instance.GetInstanceRepository()
		err = repo.Update(unmapf(&instance))

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return err
}

type payPlansRepoProviders struct {
	payPlans instance_provider.IPayPlansProvider
}

type payPlansRepoAccessors struct {
	instance access.IInstance
}

type payPlansService struct {
	repos         payPlansRepoProviders
	access        payPlansRepoAccessors
	authenticator authenticator.IAuthenticator
}

func NewPayPlans(
	authenticator authenticator.IAuthenticator,
	payPlans instance_provider.IPayPlansProvider,
	instance access.IInstance,
) instance.IPayPlansService {
	return &payPlansService{
		payPlansRepoProviders{payPlans},
		payPlansRepoAccessors{instance},
		authenticator,
	}
}

func mapPayPlans(value *models.InstancePayPlans) Collection[instance.PayPlan] {
	out := make([]instance.PayPlan, len(value.Map))
	i := 0

	for _, payPlan := range value.Map {
		out[i] = instance.PayPlan{
			InstanceId: payPlan.Id,
			PeriodId:   payPlan.PeriodId,
			Price:      payPlan.Price,
		}
		i++
	}

	return SliceCollection(out)
}

func (self *payPlansService) GetInstancePayPlans(
	instanceId uuid.UUID,
) (Collection[instance.PayPlan], error) {
	repo := self.repos.payPlans.GetInstancePayPlansRepository()
	payPlans, err := repo.GetByInstanceId(instanceId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return mapPayPlans(&payPlans), err
}

func unmapPayPlan(value *instance.PayPlanUpdateForm) models.PayPlan {
	return models.PayPlan{
		PeriodId: value.PeriodId,
		Price:    value.Price,
	}
}

func (self *payPlansService) UpdateInstancePayPlans(
	token token.Token,
	instanceId uuid.UUID,
	payPlans instance.PayPlansUpdateForm,
) error {
	var reference models.InstancePayPlans
	var notFound []instance.PayPlanUpdateForm
	repo := self.repos.payPlans.GetInstancePayPlansRepository()
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.instance.Access(user.Id, instanceId)
	}

	if nil == err {
		reference, err = repo.GetByInstanceId(instanceId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	if nil == err {
		updated := make(map[uuid.UUID]models.PayPlan)

		for i := 0; nil == err && len(payPlans) > i; i++ {
			plan := payPlans[i]

			if id, ok := mapfuncs.FindByValueF(
				reference.Map,
				func(v *models.PayPlan) bool {
					return v.PeriodId == plan.PeriodId
				},
			); !ok {
				notFound = append(notFound, plan)
			} else {
				if _, found := updated[id]; found {
					err = cmnerrors.Incorrect("Duplicate id")
				} else {
					v := reference.Map[id]
					v.Price = plan.Price
					updated[id] = v
				}
			}
		}

		if nil == err {
			reference.Map = updated
		}
	}

	if nil == err {
		err = repo.Update(reference)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	for i := 0; nil == err && len(notFound) > i; i++ {
		_, err = repo.AddPayPlan(instanceId, unmapPayPlan(&notFound[i]))

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return err
}

type photoRepoProviders struct {
	photo instance_provider.IPhotoProvider
}

type photoAccessors struct {
	instance access.IInstance
}

type photoService struct {
	repos         photoRepoProviders
	access        photoAccessors
	authenticator authenticator.IAuthenticator
	registry      photoregistry.IRegistry
}

func NewPhoto(
	authenticator authenticator.IAuthenticator,
	registry photoregistry.IRegistry,
	photo instance_provider.IPhotoProvider,
	instance access.IInstance,
) instance.IPhotoService {
	return &photoService{
		photoRepoProviders{photo},
		photoAccessors{instance},
		authenticator,
		registry,
	}
}

func (self *photoService) ListInstancePhotos(
	instanceId uuid.UUID,
) (Collection[uuid.UUID], error) {
	repo := self.repos.photo.GetInstancePhotoRepository()
	photos, err := repo.GetByInstanceId(instanceId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return photos, err
}

func (self *photoService) AddInstancePhotos(
	token token.Token,
	instanceId uuid.UUID,
	tempPhotos []uuid.UUID,
) error {
	repo := self.repos.photo.GetInstancePhotoRepository()
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.instance.Access(user.Id, instanceId)
	}

	for i := 0; err == nil && len(tempPhotos) > i; i++ {
		var id uuid.UUID
		id, err = self.registry.MoveFromTemp(tempPhotos[i])

		if nil == err {
			err = repo.Create(instanceId, id)

			if cerr := (repo_errors.ErrorDuplicate{}); errors.As(err, &cerr) {
				err = cmnerrors.AlreadyExists(cerr.What...)
			} else if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
				err = cmnerrors.NotFound(cerr.What...)
			} else if nil != err {
				err = cmnerrors.Internal(cmnerrors.DataAccess(err))
			}
		}
	}

	return err
}

type reviewRepoProviders struct {
	review review_provider.IProvider
	rent   rent_provider.IProvider
}

type reviewService struct {
	repos         reviewRepoProviders
	authenticator authenticator.IAuthenticator
}

func NewReview(
	authenticator authenticator.IAuthenticator,
	review review_provider.IProvider,
	rent rent_provider.IProvider,
) instance.IReviewService {
	return &reviewService{reviewRepoProviders{review, rent}, authenticator}
}

func mapReview(value *models.Review) instance.Review {
	return instance.Review{
		Id:         value.Id,
		InstanceId: value.InstanceId,
		UserId:     value.UserId,
		Content:    value.Content,
		Rating:     value.Rating,
		Date:       date.New(value.Date),
	}
}

func checkRating(rating instance.ReviewRating) error {
	if 0 > rating || 5 < rating {
		return instance.ErrorRatingIncorrectValue{Value: rating}
	}

	return nil
}

func checkRatingRaw(rating float64) error {
	epsilon := math.Nextafter(1, 2) - 1
	if 0-epsilon > rating || 5+epsilon < rating {
		return cmnerrors.Incorrect("rating")
	}

	return nil
}

func mapReviewFilter(value *instance.ReviewFilter) (review_repository.Filter, error) {
	out := review_repository.Filter{
		InstanceId: value.InstanceId,
		Ratings:    make([]review_repository.Rating, len(value.Ratings)),
	}

	for i, v := range value.Ratings {
		if err := checkRating(v); nil != err {
			return review_repository.Filter{}, err
		}

		out.Ratings[i] = review_repository.Rating(v)
	}

	return out, nil
}

func mapReviewSort(value *instance.ReviewSort) (review_repository.Sort, error) {
	switch *value {
	case instance.REVIEW_SORT_NONE:
		return review_repository.SORT_NONE, nil
	case instance.REVIEW_SORT_DATE_ASC:
		return review_repository.SORT_DATE_ASC, nil
	case instance.REVIEW_SORT_DATE_DSC:
		return review_repository.SORT_DATE_DSC, nil
	case instance.REVIEW_SORT_RATING_ASC:
		return review_repository.SORT_RATING_ASC, nil
	case instance.REVIEW_SORT_RATING_DSC:
		return review_repository.SORT_RATING_DSC, nil
	default:
		return review_repository.SORT_NONE, cmnerrors.Unknown("sort")
	}

}

func (self *reviewService) ListInstanceReviews(
	filter instance.ReviewFilter,
	sort instance.ReviewSort,
) (Collection[instance.Review], error) {
	var reviews Collection[models.Review]
	var filterr review_repository.Filter
	sortr, err := mapReviewSort(&sort)

	if nil == err {
		filterr, err = mapReviewFilter(&filter)
	}

	if nil == err {
		repo := self.repos.review.GetReviewRepository()
		reviews, err = repo.GetWithFilter(filterr, sortr)
	}

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return MapCollection(mapReview, reviews), err
}

func (self *reviewService) PostInstanceReview(
	token token.Token,
	instanceId uuid.UUID,
	review instance.ReviewPostForm,
) error {
	var rents Collection[records.Rent]
	var user models.User

	err := emptymathcer.Match(emptymathcer.Item("content", review.Content))

	if nil == err {
		err = checkRatingRaw(review.Rating)
	}

	if nil == err {
		user, err = self.authenticator.LoginWithToken(token)
	}

	if nil == err {
		repo := self.repos.rent.GetRentRepository()
		rents, err = repo.GetPastByUserId(user.Id)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	if nil == err {
		if _, found := Find(rents.Iter(), func(rent *records.Rent) bool {
			return rent.InstanceId == instanceId
		}); !found {
			err = cmnerrors.Authorization(cmnerrors.NoAccess("instance"))
		}
	}

	if nil == err {
		review := models.Review{
			InstanceId: instanceId,
			UserId:     user.Id,
			Content:    review.Content,
			Rating:     review.Rating,
			Date:       time.Now(),
		}
		repo := self.repos.review.GetReviewRepository()
		_, err = repo.Create(review)

		if cerr := (repo_errors.ErrorDuplicate{}); errors.As(err, &cerr) {
			err = cmnerrors.AlreadyExists(cerr.What...)
		} else if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return err
}

