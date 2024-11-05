package requests

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/requests"
	"time"

	"github.com/google/uuid"
)

type ProvideBuilder struct {
	id               uuid.UUID
	productId        uuid.UUID
	renterId         uuid.UUID
	pickUpPointId    uuid.UUID
	payPlans         map[uuid.UUID]models.PayPlan
	name             string
	description      string
	condition        string
	verificationCode string
	createDate       time.Time
}

func NewProvide() *ProvideBuilder {
	return &ProvideBuilder{}
}

func (self *ProvideBuilder) WithId(id uuid.UUID) *ProvideBuilder {
	self.id = id
	return self
}

func (self *ProvideBuilder) WithProductId(productId uuid.UUID) *ProvideBuilder {
	self.productId = productId
	return self
}

func (self *ProvideBuilder) WithRenterId(renterId uuid.UUID) *ProvideBuilder {
	self.renterId = renterId
	return self
}

func (self *ProvideBuilder) WithPickUpPointId(pickUpPointId uuid.UUID) *ProvideBuilder {
	self.pickUpPointId = pickUpPointId
	return self
}

func (self *ProvideBuilder) WithPayPlans(plans ...models.PayPlan) *ProvideBuilder {
	self.payPlans = make(map[uuid.UUID]models.PayPlan, len(plans))

	for _, v := range plans {
		self.payPlans[v.PeriodId] = v
	}

	return self
}

func (self *ProvideBuilder) WithName(name string) *ProvideBuilder {
	self.name = name
	return self
}

func (self *ProvideBuilder) WithDescription(description string) *ProvideBuilder {
	self.description = description
	return self
}

func (self *ProvideBuilder) WithCondition(condition string) *ProvideBuilder {
	self.condition = condition
	return self
}

func (self *ProvideBuilder) WithVerificationCode(verificationCode string) *ProvideBuilder {
	self.verificationCode = verificationCode
	return self
}

func (self *ProvideBuilder) WithCreateDate(createDate time.Time) *ProvideBuilder {
	self.createDate = createDate
	return self
}

func (self *ProvideBuilder) Build() requests.Provide {
	return requests.Provide{
		Id:               self.id,
		ProductId:        self.productId,
		RenterId:         self.renterId,
		PickUpPointId:    self.pickUpPointId,
		PayPlans:         self.payPlans,
		Name:             self.name,
		Description:      self.description,
		Condition:        self.condition,
		VerificationCode: self.verificationCode,
		CreateDate:       self.createDate,
	}
}

