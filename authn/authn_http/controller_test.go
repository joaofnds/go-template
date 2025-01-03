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

	BeforeAll(func() { app = harness.Setup(); api = app.NewDriver() })
	BeforeEach(func() { app.BeforeEach() })
	AfterEach(func() { app.AfterEach() })
	AfterAll(func() { app.Teardown() })

	Describe("register", func() {
		email := "me@template.com"

		It("creates the user", func() {
			user := api.Auth.MustRegister(email, "p455w0rd")

			Expect(api.Users.Get(user.ID)).To(Equal(user))
		})

		It("logs in after registration", func() {
			api.Auth.MustRegister(email, "p455w0rd")

			token, err := api.Auth.Login(email, "p455w0rd")
			Expect(err).To(BeNil())
			Expect(token).NotTo(BeNil())
		})
	})

	Describe("login", func() {
		It("returns token", func() {
			token := api.Auth.MustLogin("admin", "123")

			Expect(token.TokenType).To(Equal("Bearer"))
			Expect(token.AccessToken).NotTo(BeEmpty())
			Expect(token.RefreshToken).NotTo(BeEmpty())
			Expect(token.Expiry).To(BeTemporally("~", time.Now().Add(7*24*time.Hour), time.Second))
		})
	})

	Describe("user info", func() {
		It("returns user info", func() {
			userDriver := app.DriverFor("admin", "123")
			userInfo := userDriver.Auth.MustUserInfo()

			Expect(userInfo).To(HaveKeyWithValue("email", "admin@example.com"))
		})
	})
})
