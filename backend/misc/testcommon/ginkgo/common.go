package ginkgo

import (
	"fmt"

	"github.com/onsi/ginkgo/v2"
)

type GinkgoProviderWrap struct{}

func (GinkgoProviderWrap) Errorf(format string, args ...interface{}) {
	ginkgo.GinkgoHelper()
	ginkgo.Fail(fmt.Sprintf(format, args...))
}

func (GinkgoProviderWrap) Logf(fmt string, args ...interface{}) {
	ginkgo.GinkgoWriter.Printf(fmt, args...)
}

