package authn_http_test

import (
	"app/test/driver"
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
	var app *driver.Driver

	BeforeAll(func() { app = driver.Setup() })
	BeforeEach(func() { app.BeginTx() })
	AfterEach(func() { app.RollbackTx() })
	AfterAll(func() { app.Teardown() })

	Describe("register", func() {
		email := "me@template.com"

		BeforeEach(func() { app.Auth.MustDelete(email) })
		AfterEach(func() { app.Auth.MustDelete(email) })

		It("creates the user", func() {
			user := app.Auth.MustRegister(email, "p455w0rd")

			Expect(app.Users.Get(user.ID)).To(Equal(user))
		})

		It("logs in after registration", func() {
			app.Auth.MustRegister(email, "p455w0rd")

			token, err := app.Auth.Login(email, "p455w0rd")
			Expect(err).To(BeNil())
			Expect(token).NotTo(BeNil())
		})
	})

	Describe("login", func() {
		It("returns token", func() {
			token := app.Auth.MustLogin("admin", "123")

			Expect(token.TokenType).To(Equal("Bearer"))
			Expect(token.AccessToken).NotTo(BeEmpty())
			Expect(token.RefreshToken).NotTo(BeEmpty())
			Expect(token.Expiry).To(BeTemporally("~", time.Now().Add(7*24*time.Hour), time.Second))
		})
	})

	Describe("user info", func() {
		It("returns user info", func() {
			app.Auth.MustLogin("admin", "123")
			userInfo := app.Auth.MustUserInfo()

			Expect(userInfo).To(HaveKeyWithValue("email", "admin@example.com"))
		})
	})
})
