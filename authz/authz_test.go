package authz_test

import (
	"testing"

	"app/adapter/logger"
	"app/adapter/postgres"
	"app/authz"
	"app/config"
	"app/internal/ref"
	"app/test"
	. "app/test/matchers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestAuthz(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "authz suite")
}

var _ = Describe("authz service", func() {
	user := ref.New("111", "user")
	post := ref.New("222", "post")
	anyPost := ref.New("*", "post")
	admin := ref.New("admin", "role")

	adminAnyPostDelete := authz.NewAppRequest(admin, anyPost, "delete")
	userPostDelete := authz.NewAppRequest(user, post, "delete")

	var app *fxtest.App
	var sut *authz.Service

	BeforeEach(func() {
		app = fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			test.Transaction,
			test.CasbinStringAdapter,
			config.Module,
			postgres.Module,
			authz.Module,
			fx.Populate(&sut),
		)
		app.RequireStart()
	})

	AfterEach(func() {
		app.RequireStop()
	})

	Describe("roles", func() {
		It("grants roles", func() {
			Must(sut.GrantRole(user, admin))
		})

		It("retrieves roles", func() {
			Must(sut.GrantRole(user, admin))

			Expect(sut.GetRoles(user)).To(ConsistOf(admin))
		})

		It("revokes roles", func() {
			Must(sut.GrantRole(user, admin))

			Must(sut.RevokeRole(user, admin))

			Expect(sut.GetRoles(user)).To(BeEmpty())
		})
	})

	It("has permission after direct grant", func(ctx SpecContext) {
		Expect(sut.Check(userPostDelete)).To(BeFalse())

		Must(sut.Grant(userPostDelete))

		Expect(sut.Check(userPostDelete)).To(BeTrue())
	})

	It("has permission after role grant", func(ctx SpecContext) {
		Must(sut.Grant(adminAnyPostDelete))
		Expect(sut.Check(userPostDelete)).To(BeFalse())

		Must(sut.GrantRole(user, admin))

		Expect(sut.Check(userPostDelete)).To(BeTrue())
	})
})
