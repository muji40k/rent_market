package user_test

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	models_b "rent_service/builders/domain/models"
	models_om "rent_service/builders/mothers/domain/models"
	mserver "rent_service/builders/mothers/test/application/server"
	"rent_service/builders/mothers/test/repository/psql"
	"rent_service/builders/mothers/test/service/defservices"
	"rent_service/internal/domain/models"
	sv1 "rent_service/internal/logic/context/v1"
	simplecodegen "rent_service/internal/logic/services/implementations/defservices/codegen/simple"
	"rent_service/internal/logic/services/implementations/defservices/codenc"
	servuser "rent_service/internal/logic/services/implementations/defservices/services/user"
	siuser "rent_service/internal/logic/services/interfaces/user"
	rv1 "rent_service/internal/repository/context/v1"
	repouser "rent_service/internal/repository/implementation/sql/repositories/user"
	simplesetter "rent_service/internal/repository/implementation/sql/technical/implementations/simple"
	riuser "rent_service/internal/repository/interfaces/user"
	"rent_service/misc/nullable"
	defservicescommon "rent_service/misc/testcommon/defservices"
	ginkgocommon "rent_service/misc/testcommon/ginkgo"
	psqlcommon "rent_service/misc/testcommon/psql"
	servercommon "rent_service/misc/testcommon/server"
	"rent_service/server"
	"rent_service/server/authenticator"
	"rent_service/server/controllers/users/self/update-requests"
	"rent_service/server/headers"
	"strings"
	"testing"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/mail"
	"github.com/gavv/httpexpect/v2"
	"github.com/jmoiron/sqlx"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

const (
	IMAP_PORT           string = "TEST_IMAP_PORT"
	RECEPTIENT_EMAIL    string = "TEST_RECEPIENT_EMAIL"
	RECEPTIENT_NAME     string = "TEST_RECEPIENT_USERNAME"
	RECEPTIENT_PASSWORD string = "TEST_RECEPIENT_PASSWORD"
)

func UserRecepient() *models_b.UserBuilder {
	return models_om.UserDefault(nullable.None[string]()).
		WithEmail(os.Getenv(RECEPTIENT_EMAIL)).
		WithName(os.Getenv(RECEPTIENT_NAME))
}

func TestRunner(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "User password change with 2FA")
}

var (
	seContext servercommon.Context
	sContext  defservicescommon.Context
	rContext  psqlcommon.Context
	db        *sqlx.DB
)

var _ = ginkgo.BeforeSuite(func() {
	rContext = psqlcommon.New()
	sContext = defservicescommon.New(rContext.Factory.ToFactories())

	dbconfig := psql.Parse()
	db = nullable.UnwrapF(dbconfig.GetConnection())

	seContext = servercommon.New(
		sContext.Factory.ToFactories(),
		append(mserver.DefaultControllers(), func(
			a authenticator.IAuthenticator,
			_ *sv1.Context,
		) server.IController {
			return updaterequests.New(a,
				defservices.UserPasswordUpdateRequestProvider(
					func() siuser.IPasswordUpdateService {
						return servuser.NewPasswordUpdate(
							sContext.Factory.CreateAuthenticator(),
							defservices.Email2FA(),
							simplecodegen.New(6),
							codenc.MD5CodeEncoder{},
							time.Minute,
							rv1.New(rContext.Factory.ToFactories()),
							psql.UserPasswordUpdateRequestProvider(
								func() riuser.IPasswordUpdateRepository {
									return repouser.NewPasswordUpdate(
										db,
										simplesetter.New("test"),
									)
								},
							),
						)
					},
				),
			)
		})...,
	)
})

var _ = ginkgo.AfterSuite(func() {
	seContext.Close()
	sContext.Close()
	rContext.Close()
	if nil != db {
		db.Close()
	}
})

var _ = ginkgo.Describe("Changing password with email 2FA", func() {
	var (
		user        models.User
		client      *httpexpect.Expect
		newpassword string = "some_new_password_for_test_purposes"
		resp        *httpexpect.Response
	)

	ginkgo.BeforeEach(func() {
		user = UserRecepient().Build()
		client = seContext.GetClient(
			ginkgocommon.GinkgoProviderWrap{},
			func(t servercommon.Provider) httpexpect.Printer {
				return httpexpect.NewDebugPrinter(t, true)
			},
		)
		gomega.Expect(user).NotTo(gomega.BeZero())
		gomega.Expect(client).NotTo(gomega.BeNil())
		rContext.Inserter.InsertUser(&user)
	})

	ginkgo.AfterEach(func() {
		ResetSeen()
		nullable.UnwrapF(db.Exec("delete from users.password_update_requests where user_id = $1", user.Id))
		nullable.UnwrapF(db.Exec("delete from users.users where id = $1", user.Id))
	})

	checkPsswd := func(old bool) {
		resp := client.POST("/sessions").WithJSON(map[string]interface{}{
			"email":    user.Email,
			"password": newpassword,
		}).Expect()

		if old {
			resp.Status(http.StatusUnauthorized)
		} else {
			resp.Status(http.StatusOK)
		}

		resp = client.POST("/sessions").WithJSON(map[string]interface{}{
			"email":    user.Email,
			"password": user.Password,
		}).Expect()

		if old {
			resp.Status(http.StatusOK)
		} else {
			resp.Status(http.StatusUnauthorized)
		}
	}

	ginkgo.When("I try to change password for my account", func() {
		var token authenticator.ApiToken
		var oldPassword string

		ginkgo.BeforeEach(func() {
			resp := client.POST("/sessions").WithJSON(map[string]interface{}{
				"email":    user.Email,
				"password": user.Password,
			}).Expect()

			resp.Status(http.StatusOK).JSON().Decode(&token)
			gomega.Expect(token.Access).NotTo(gomega.BeZero())
			gomega.Expect(token.Renew).NotTo(gomega.BeZero())
		})

		ginkgo.JustBeforeEach(func() {
			resp = client.POST("/users/self/update-requests/passwords").
				WithHeader(headers.API_KEY, token.Access).
				WithJSON(map[string]interface{}{
					"old_password": oldPassword,
					"new_password": newpassword,
				}).Expect()
		})

		ginkgo.Context("enter correct old password", func() {
			var updateRequest siuser.PasswordUpdateRequest

			ginkgo.BeforeEach(func() {
				oldPassword = user.Password
			})

			ginkgo.JustBeforeEach(func() {
				resp.Status(http.StatusOK).JSON().Decode(&updateRequest)
				gomega.Expect(updateRequest.Required).To(gomega.BeTrue())
				gomega.Expect(updateRequest.Id).NotTo(gomega.BeZero())
				gomega.Expect(updateRequest.ValidTo).To(gomega.BeTemporally(">", time.Now()))
			})

			ginkgo.Context("and provide valid 2FA code from email", func() {
				ginkgo.It("no error is retuned and my password is changed", func() {
					resp = client.DELETE(
						"/users/self/update-requests/passwords/{id}", updateRequest.Id,
					).WithHeader(
						headers.API_KEY, token.Access,
					).WithJSON(map[string]interface{}{
						"code": codenc.MD5CodeEncoder{}.Encode(
							GetCodeFromEmail(),
						),
					}).Expect()

					resp.Status(http.StatusOK).Body().IsEmpty()
					checkPsswd(false)
				})
			})

			ginkgo.Context("and provide invalid 2FA code from email", func() {
				ginkgo.It("status Forbidden is returned and password stays the same", func() {
					resp = client.DELETE(
						"/users/self/update-requests/passwords/{id}", updateRequest.Id,
					).WithHeader(
						headers.API_KEY, token.Access,
					).WithJSON(map[string]interface{}{
						"code": codenc.MD5CodeEncoder{}.Encode(
							"definetly_invalid_code",
						),
					}).Expect()

					resp.Status(http.StatusForbidden).Body().IsEmpty()
					checkPsswd(true)
				})
			})
		})

		ginkgo.Context("enter incorrect old password", func() {
			ginkgo.BeforeEach(func() {
				oldPassword = "password_that_is_definetly_incorrect"
			})

			ginkgo.It("status Unauthorized is returned and password stays the same", func() {
				resp.Status(http.StatusUnauthorized).Body().IsEmpty()
				checkPsswd(true)
			})
		})
	})
})

func PipeError[T any](_ T, err ...error) error {
	return err[0]
}

type SecretENV string

func (SecretENV) String() string {
	return "[hidden]"
}

func FromENV(name string) SecretENV {
	return SecretENV(os.Getenv(name))
}

func GetClient() (*imapclient.Client, error) {
	tlconf := tls.Config{
		ServerName: os.Getenv(defservices.MAIL_SERVER),
	}
	opts1 := imapclient.Options{
		TLSConfig: &tlconf,
	}

	server := FromENV(defservices.MAIL_SERVER)
	port := FromENV(IMAP_PORT)
	user := FromENV(RECEPTIENT_NAME)
	password := FromENV(RECEPTIENT_PASSWORD)

	file, _ := os.Create("secrets.txt")

	fmt.Fprintf(file, "server: %v", server)
	fmt.Fprintf(file, "port: %v", port)
	fmt.Fprintf(file, "user: %v", user)
	fmt.Fprintf(file, "password: %v", password)

	file.Close()

	c, err := imapclient.DialStartTLS(
		string(server)+":"+string(port),
		&opts1,
	)

	if nil == err {
		err = c.Login(
			string(user),
			string(password),
		).Wait()
	}

	return c, err
}

func GetCodeFromEmail() string {
	c, err := GetClient()
	gomega.Expect(c).NotTo(gomega.BeNil())
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	gomega.Expect(
		PipeError(c.Select("INBOX", nil).Wait()),
	).NotTo(gomega.HaveOccurred())

	var data *imap.SearchData
	var criteria = imap.SearchCriteria{
		NotFlag: []imap.Flag{imap.FlagSeen},
		Header: []imap.SearchCriteriaHeaderField{
			{Key: "From", Value: os.Getenv(defservices.SENDER_EMAIL)},
		},
	}
	data, err = c.Search(&criteria, nil).Wait()

	gomega.Expect(data).NotTo(gomega.BeNil())
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(len(data.AllSeqNums())).NotTo(gomega.Equal(0))

	opts2 := imap.FetchOptions{
		BodySection: []*imap.FetchItemBodySection{{}},
	}

	cmd := c.Fetch(data.All, &opts2)
	gomega.Expect(cmd).NotTo(gomega.BeNil())

	md := cmd.Next()
	gomega.Expect(md).NotTo(gomega.BeNil())

	var body imapclient.FetchItemDataBodySection
	var ok bool = false
	for item := md.Next(); !ok && nil != item; {
		body, ok = item.(imapclient.FetchItemDataBodySection)

		if !ok {
			item = md.Next()
		}
	}
	gomega.Expect(body).NotTo(gomega.BeZero())
	gomega.Expect(ok).To(gomega.BeTrue())

	var reader *mail.Reader
	reader, err = mail.CreateReader(body.Literal)
	gomega.Expect(reader).NotTo(gomega.BeNil())
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	var prt *mail.Part
	prt, err = reader.NextPart()
	gomega.Expect(prt).NotTo(gomega.BeNil())
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	var b []byte
	b, err = io.ReadAll(prt.Body)
	gomega.Expect(b).NotTo(gomega.BeNil())
	gomega.Expect(len(b)).NotTo(gomega.Equal(0))
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	gomega.Expect(reader.Close()).NotTo(gomega.HaveOccurred())
	gomega.Expect(cmd.Close()).NotTo(gomega.HaveOccurred())
	gomega.Expect(c.Logout().Wait()).NotTo(gomega.HaveOccurred())

	text := strings.TrimSpace(string(b))
	lines := strings.Split(text, "\n")
	gomega.Expect(len(lines)).To(gomega.Equal(4))
	fields := strings.Fields(lines[1])
	gomega.Expect(len(fields)).To(gomega.Equal(4))

	return fields[3][:6]
}

func ResetSeen() {
	c, err := GetClient()

	if nil == err {
		_, err = c.Select("INBOX", nil).Wait()
	}

	var data *imap.SearchData
	if nil == err {
		var criteria = imap.SearchCriteria{
			NotFlag: []imap.Flag{imap.FlagSeen},
		}
		data, err = c.Search(&criteria, nil).Wait()
	}

	if nil == err && 0 != len(data.AllSeqNums()) {
		opts := imap.FetchOptions{
			BodySection: []*imap.FetchItemBodySection{{}},
		}
		cmd := c.Fetch(data.All, &opts)

		if nil != cmd {
			_ = cmd.Close()
		}
	}

	if nil != c {
		_ = c.Logout().Wait()
	}
}

