package casdoor_test

import (
	"app/adapter/casdoor"
	"app/adapter/logger"
	"app/adapter/validation"
	"app/config"
	"time"

	. "app/test/matchers"
	"context"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestCasdoorUserProvider(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "casdoor")
}

var _ = Describe("providers", func() {
	var (
		fxApp         *fxtest.App
		userProvider  *casdoor.UserProvider
		tokenProvider *casdoor.TokenProvider
	)

	BeforeEach(func() {
		fxApp = fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			validation.Module,
			config.Module,
			casdoor.Module,
			fx.Populate(&userProvider, &tokenProvider),
		).RequireStart()

		// cleanup users to avoid conflicts
		for _, user := range Must2(userProvider.List(context.Background())) {
			if user.Email != "admin@example.com" {
				Must(userProvider.Delete(context.Background(), user.Email))
			}
		}
	})

	AfterEach(func() { fxApp.RequireStop() })

	Describe("user provider", func() {
		Describe("create", func() {
			It("works", func() {
				err := userProvider.Create(context.Background(), "alice@template.com", "p455w0rd")
				Expect(err).To(BeNil())
			})
		})

		Describe("list", func() {
			It("lists created users", func() {
				Must(userProvider.Create(context.Background(), "alice@template.com", "p455w0rd"))

				users, err := userProvider.List(context.Background())
				Expect(err).To(BeNil())
				Expect(users).To(HaveLen(2))
				Expect(users[0].Email).To(Equal("alice@template.com"))
				Expect(users[1].Email).To(Equal("admin@example.com"))
			})
		})

		Describe("delete", func() {
			It("deletes a user", func() {
				Expect(Must2(userProvider.List(context.Background()))).To(HaveLen(1))

				Must(userProvider.Create(context.Background(), "alice@template.com", "p455w0rd"))
				Expect(Must2(userProvider.List(context.Background()))).To(HaveLen(2))

				Must(userProvider.Delete(context.Background(), "alice@template.com"))
				Expect(Must2(userProvider.List(context.Background()))).To(HaveLen(1))
			})
		})
	})

	Describe("token provider", func() {
		Describe("get", func() {
			It("returns an oauth2 token", func() {
				Must(userProvider.Create(context.Background(), "alice@template.com", "p455w0rd"))

				token, err := tokenProvider.Get(context.Background(), "alice@template.com", "p455w0rd")
				Expect(err).To(BeNil())
				Expect(token.TokenType).To(Equal("Bearer"))
				Expect(token.AccessToken).NotTo(BeEmpty())
				Expect(token.RefreshToken).NotTo(BeEmpty())
				Expect(token.Expiry).To(BeTemporally("~", time.Now().Add(7*24*time.Hour), time.Second))
			})

			When("user does not exist", func() {
				It("returns an error", func() {
					_, err := tokenProvider.Get(context.Background(), "alice@template.com", "wrong-password")
					Expect(err.Error()).To(ContainSubstring("the user does not exist"))
				})
			})

			When("password is incorrect", func() {
				It("returns an error", func() {
					Must(userProvider.Create(context.Background(), "alice@template.com", "p455w0rd"))

					_, err := tokenProvider.Get(context.Background(), "alice@template.com", "wrong-password")
					Expect(err.Error()).To(ContainSubstring("invalid username or password"))
				})
			})
		})

		Describe("parse", func() {
			It("parses a token", func() {
				Must(userProvider.Create(context.Background(), "alice@template.com", "p455w0rd"))

				token := Must2(tokenProvider.Get(context.Background(), "alice@template.com", "p455w0rd"))

				claims, err := tokenProvider.Parse(token.AccessToken)
				Expect(err).To(BeNil())
				Expect(claims.Subject).To(HaveLen(36))
				Expect(claims.Email).To(Equal("alice@template.com"))
			})

			When("token is invalid", func() {
				It("returns an error", func() {
					claims := jwt.MapClaims{"sub": "alice", "email": "alice@template.com"}
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
					tokenString := Must2(token.SignedString([]byte("secret")))

					_, err := tokenProvider.Parse(tokenString)
					Expect(err.Error()).To(ContainSubstring("token signature is invalid"))
				})
			})
		})
	})
})
