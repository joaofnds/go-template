package authn_http_test

import (
	"app/test/driver"
	"app/test/harness"
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

	BeforeAll(func() { app = harness.Setup(); api = app.NewDriver() })
	BeforeEach(func() { app.BeforeEach() })
	AfterEach(func() { app.AfterEach() })
	AfterAll(func() { app.Teardown() })

	Describe("register", func() {
		It("creates the user", func() {
			user := api.Auth.MustRegister(email, password)

			Expect(api.Users.Get(user.ID)).To(Equal(user))
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
			api.Auth.MustRegister(email, password)
			userDriver := app.DriverFor(email, password)
			userInfo := userDriver.Auth.MustUserInfo()

			Expect(userInfo).To(HaveKeyWithValue("email", "user@template.com"))
		})
	})
})
