package authentication_test

import (
	"net/http"
	models_om "rent_service/builders/mothers/domain/models"
	mserver "rent_service/builders/mothers/test/application/server"
	"rent_service/internal/domain/models"
	"rent_service/misc/nullable"
	"rent_service/misc/testcommon/defservices"
	ginkgocommon "rent_service/misc/testcommon/ginkgo"
	psqlcommon "rent_service/misc/testcommon/psql"
	"rent_service/misc/testcommon/server"
	"rent_service/server/authenticator"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestRunner(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Authentication Suite")
}

var (
	seContext server.Context
	sContext  defservices.Context
	rContext  psqlcommon.Context
)

var _ = ginkgo.BeforeSuite(func() {
	rContext = psqlcommon.New()
	sContext = defservices.New(rContext.Factory.ToFactories())
	seContext = server.New(
		sContext.Factory.ToFactories(),
		mserver.DefaultControllers()...,
	)
})

var _ = ginkgo.AfterSuite(func() {
	seContext.Close()
	sContext.Close()
	rContext.Close()
})

var _ = ginkgo.Describe("Authetication", func() {
	var (
		user   models.User
		client *httpexpect.Expect
		resp   *httpexpect.Response
	)

	ginkgo.BeforeEach(func() {
		user = models_om.UserDefault(nullable.None[string]()).Build()
		client = seContext.GetClient(ginkgocommon.GinkgoProviderWrap{})
		gomega.Expect(user).NotTo(gomega.BeZero())
		gomega.Expect(client).NotTo(gomega.BeNil())
		rContext.Inserter.InsertUser(&user)
	})

	ginkgo.When("I login into my account", func() {
		var (
			email    string
			password string
		)

		ginkgo.JustBeforeEach(func() {
			resp = client.POST("/sessions").WithJSON(map[string]interface{}{
				"email":    email,
				"password": password,
			}).Expect()
		})

		ginkgo.Context("and provide valid credentials (username and password)", func() {
			var renterToken authenticator.ApiToken

			ginkgo.BeforeEach(func() {
				email = user.Email
				password = user.Password
			})

			ginkgo.It("valid token with OK status is returned", func() {
				resp.Status(http.StatusOK).JSON().Object().Decode(&renterToken)
				gomega.Expect(renterToken.Access).NotTo(gomega.BeZero())
				gomega.Expect(renterToken.Renew).NotTo(gomega.BeZero())
			})
		})

		ginkgo.Context("and provide invalid credentials (username and password)", func() {
			checker := func() {
				resp.Status(http.StatusUnauthorized).Body().IsEmpty()
			}

			ginkgo.Describe("unknown email", func() {
				ginkgo.BeforeEach(func() {
					email = "unknown@some_mail.su"
					password = user.Password
				})

				ginkgo.It("no token with Unauthorized status is returned", checker)
			})

			ginkgo.Describe("wrong password", func() {
				ginkgo.BeforeEach(func() {
					email = user.Email
					password = "some_password_that_is_definetly_incorrect"
				})

				ginkgo.It("no token with Unauthorized status is returned", checker)
			})

			ginkgo.Describe("unknown email and wrong password", func() {
				ginkgo.BeforeEach(func() {
					email = "unknown@some_mail.su"
					password = "some_password_that_is_definetly_incorrect"
				})

				ginkgo.It("no token with Unauthorized status is returned", checker)
			})

			ginkgo.Describe("credentials from my other account", func() {
				ginkgo.BeforeEach(func() {
					otherUser := models_om.UserDefault(nullable.None[string]()).Build()
					email = user.Email
					password = otherUser.Password
					rContext.Inserter.InsertUser(&otherUser)
				})

				ginkgo.It("no token with Unauthorized status is returned", checker)
			})
		})
	})
})

