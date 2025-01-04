package casdoor_test

import (
	"app/adapter/casdoor"
	"app/adapter/logger"
	"app/adapter/validation"
	"app/authn"
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
		email         = "alice@template.com"
		password      = "p455w0rd"
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

		Must(userProvider.Delete(context.Background(), email))
	})

	AfterEach(func() { fxApp.RequireStop() })

	Describe("user provider", func() {
		Describe("create", func() {
			It("works", func() {
				err := userProvider.Create(context.Background(), email, password)
				Expect(err).To(BeNil())
			})
		})

		Describe("delete", func() {
			It("deletes a user", func() {
				_, err := tokenProvider.Get(context.Background(), email, password)
				Expect(err).To(MatchError(authn.ErrUserNotFound))

				Must(userProvider.Create(context.Background(), email, password))

				_, err = tokenProvider.Get(context.Background(), email, password)
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("token provider", func() {
		Describe("get", func() {
			It("returns an oauth2 token", func() {
				Must(userProvider.Create(context.Background(), email, password))

				token, err := tokenProvider.Get(context.Background(), email, password)
				Expect(err).To(BeNil())
				Expect(token.TokenType).To(Equal("Bearer"))
				Expect(token.AccessToken).NotTo(BeEmpty())
				Expect(token.RefreshToken).NotTo(BeEmpty())
				Expect(token.Expiry).To(BeTemporally("~", time.Now().Add(7*24*time.Hour), time.Second))
			})

			When("user does not exist", func() {
				It("returns an error", func() {
					_, err := tokenProvider.Get(context.Background(), email, password)
					Expect(err).To(MatchError(authn.ErrUserNotFound))
				})
			})

			When("password is incorrect", func() {
				It("returns an error", func() {
					Must(userProvider.Create(context.Background(), email, password))

					_, err := tokenProvider.Get(context.Background(), email, "wrong-password")
					Expect(err).To(MatchError(authn.ErrWrongPassword))
				})
			})
		})

		Describe("parse", func() {
			It("parses a token", func() {
				Must(userProvider.Create(context.Background(), email, password))

				token := Must2(tokenProvider.Get(context.Background(), email, password))

				claims, err := tokenProvider.Parse(token.AccessToken)
				Expect(err).To(BeNil())
				Expect(claims.Subject).To(HaveLen(36))
				Expect(claims.Email).To(Equal(email))
			})

			When("token is invalid", func() {
				It("returns an error", func() {
					claims := jwt.MapClaims{"sub": "alice", "email": email}
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
					tokenString := Must2(token.SignedString([]byte("secret")))

					_, err := tokenProvider.Parse(tokenString)
					Expect(err).To(MatchError(authn.ErrParseToken))
				})
			})
		})
	})
})
