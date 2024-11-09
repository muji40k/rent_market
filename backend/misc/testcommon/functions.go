package testcommon

import (
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

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

