package testcommon

import (
	"math"
	"rent_service/internal/misc/types/currency"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

var EPSILON float64 = math.Nextafter(1, 2) - 1

func SetBase(t provider.T, parent string, epic string, feature string) {
	t.AddParentSuite(parent)
	t.Epic(epic)
	t.Feature(feature)
}

func MethodDescriptor(subSuite string, story string) func(t provider.T, title string, description string) {
	return func(t provider.T, title string, description string) {
		t.AddSubSuite(subSuite)
		t.Story(story)
		t.Title(title)
		t.Description(description)
	}
}

func AssignParameter[T any](ctx provider.StepCtx, name string, value T) T {
	ctx.WithParameters(allure.NewParameter(name, value))
	return value
}

type IProvider interface {
	WithNewStep(string, func(sCtx provider.StepCtx), ...*allure.Parameter)
}

type Asserter[T any] interface {
	EqualFunc(cmp func(T, T) bool, expected T, actual T, name string)
	ElementsMatchFunc(cmp func(T, T) bool, expected []T, actual []T, name string)
}

type asserter[T any] struct {
	provider IProvider
	getter   func(provider.StepCtx) provider.Asserts
	prefix   string
}

func Assert[T any](p IProvider) Asserter[T] {
	return &asserter[T]{
		p,
		func(sCtx provider.StepCtx) provider.Asserts {
			return sCtx.Assert()
		},
		"ASSERT: ",
	}
}

func Require[T any](p IProvider) Asserter[T] {
	return &asserter[T]{
		p,
		func(sCtx provider.StepCtx) provider.Asserts {
			return sCtx.Require()
		},
		"REQUIRE: ",
	}
}

// Adopted from testify
func diffLists[T any, F any](cmp func(T, F) bool, expected []T, actual []F) ([]T, []F) {
	eEx := make([]T, 0)
	aEx := make([]F, 0)
	eLen := len(expected)
	aLen := len(actual)
	aVisited := make([]bool, aLen)

	for i := 0; i < eLen; i++ {
		found := false

		for j := 0; !found && j < aLen; j++ {
			if !aVisited[j] && cmp(expected[i], actual[j]) {
				aVisited[j] = true
				found = true
			}
		}

		if !found {
			eEx = append(eEx, expected[i])
		}
	}

	for j := 0; j < aLen; j++ {
		if !aVisited[j] {
			aEx = append(aEx, actual[j])
		}
	}

	return eEx, aEx
}

func ElementsMatchFunc[T any, F any](cmp func(T, F) bool, expected []T, actual []F) bool {
	if len(expected) != len(actual) {
		return false
	}

	extraA, extraB := diffLists(cmp, expected, actual)

	return len(extraA) == 0 && len(extraB) == 0
}

func (self *asserter[T]) ElementsMatchFunc(cmp func(T, T) bool, expected []T, actual []T, name string) {
	wrap := func(e, a []T) bool { return ElementsMatchFunc(cmp, e, a) }

	generalCustomComparator(
		self.provider.WithNewStep,
		self.getter,
		wrap,
		expected,
		actual,
		self.prefix+name,
	)
}

func (self *asserter[T]) EqualFunc(cmp func(T, T) bool, expected T, actual T, name string) {
	generalCustomComparator(
		self.provider.WithNewStep,
		self.getter,
		cmp,
		expected,
		actual,
		self.prefix+name,
	)
}

func generalCustomComparator[T any](
	s func(string, func(sCtx provider.StepCtx), ...*allure.Parameter),
	a func(provider.StepCtx) provider.Asserts,
	cmp func(T, T) bool,
	expected T,
	actual T,
	name string,
) {
	s(name, func(sCtx provider.StepCtx) {
		sCtx.WithParameters(
			allure.NewParameter("Expected", expected),
			allure.NewParameter("Actual", actual),
		)
		a(sCtx).True(cmp(expected, actual))
	})
}

func CompareCurrency(expected currency.Currency, actual currency.Currency) bool {
	return expected.Name == actual.Name &&
		EPSILON > math.Abs(expected.Value-actual.Value)
}

