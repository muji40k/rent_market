package currency

import "rent_service/internal/misc/types/currency"

type CurrencyBuilder struct {
	name  string
	value float64
}

func NewCurrency() *CurrencyBuilder {
	return &CurrencyBuilder{}
}

func (self *CurrencyBuilder) WithName(name string) *CurrencyBuilder {
	self.name = name
	return self
}

func (self *CurrencyBuilder) WithValue(value float64) *CurrencyBuilder {
	self.value = value
	return self
}

func (self *CurrencyBuilder) Build() currency.Currency {
	return currency.Currency{
		Name:  self.name,
		Value: self.value,
	}
}

