package authn_http_test

import (
	"app/test/driver"
	"app/test/harness"
	"net/http"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Test")
}

var _ = Describe("/auth", Ordered, func() {
	var app *harness.Harness
	var api *driver.Driver
	email := "user@template.com"
	password := "p455w0rd"

	BeforeAll(func() { app = harness.Setup() })
	BeforeEach(func() { app.BeforeEach(); api = app.NewDriver() })
	AfterEach(func() { app.AfterEach() })
	AfterAll(func() { app.Teardown() })

	Describe("register", func() {
		It("creates the user", func() {
			api = app.NewDriver()
			user := api.Auth.MustRegister(email, password)

			Expect(api.Users.MustList()).To(ConsistOf(user))
		})

		It("logs in after registration", func() {
			api.Auth.MustRegister(email, password)

			token, err := api.Auth.Login(email, password)
			Expect(err).To(BeNil())
			Expect(token).NotTo(BeNil())
		})
	})

	Describe("login", func() {
		It("returns token", func() {
			api.Auth.MustRegister(email, password)
			token := api.Auth.MustLogin(email, password)

			Expect(token.TokenType).To(Equal("Bearer"))
			Expect(token.AccessToken).NotTo(BeEmpty())
			Expect(token.RefreshToken).NotTo(BeEmpty())
			Expect(token.Expiry).To(BeTemporally("~", time.Now().Add(7*24*time.Hour), time.Second))
		})
	})

	Describe("user info", func() {
		It("returns user info", func() {
			createdUser := api.Auth.MustRegister(email, password)

			userDriver := app.NewDriver()
			userDriver.Login(email, password)
			Expect(userDriver.Auth.MustUserInfo()).To(Equal(createdUser))
		})
	})

	When("email is invalid", func() {
		It("returns 400", func() {
			_, err := api.Auth.Register("bad_email.com", password)
			Expect(err).To(Equal(driver.RequestFailure{
				Status: http.StatusBadRequest,
				Body:   `{"errors":["Field validation for 'Email' failed on the 'email' tag"]}`,
			}))
		})
	})
})
