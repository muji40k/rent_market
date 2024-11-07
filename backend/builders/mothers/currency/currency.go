package currency

import (
	builder "rent_service/builders/currency"
)

func Currency(name string, value float64) *builder.CurrencyBuilder {
	return builder.NewCurrency().
		WithName(name).
		WithValue(value)
}

func RUB(value float64) *builder.CurrencyBuilder {
	return Currency("rub", value)
}

func USD(value float64) *builder.CurrencyBuilder {
	return Currency("usd", value)
}

