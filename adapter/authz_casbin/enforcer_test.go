package authz_casbin_test

import (
	"app/adapter/authz_casbin"
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

var _ = Describe("casbin enforcer", func() {
	user := ref.New("111", "user")
	post := ref.New("222", "post")
	anyPost := ref.New("*", "post")
	admin := ref.New("admin", "role")

	adminAnyPostDelete := authz.NewAppRequest(admin, anyPost, "delete")
	userPostDelete := authz.NewAppRequest(user, post, "delete")

	var app *fxtest.App
	var sut *authz_casbin.Enforcer
	var roles *authz_casbin.RoleManager

	BeforeEach(func() {
		app = fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			test.CasbinStringAdapter,
			config.Module,
			postgres.Module,
			authz_casbin.Module,
			fx.Populate(&sut, &roles),
		)
		app.RequireStart()
	})

	AfterEach(func() {
		app.RequireStop()
	})

	It("has permission after direct grant", func(ctx SpecContext) {
		Expect(sut.Check(userPostDelete)).To(BeFalse())

		Must(sut.Grant(userPostDelete))

		Expect(sut.Check(userPostDelete)).To(BeTrue())
	})

	It("has permission after role grant", func(ctx SpecContext) {
		Must(sut.Grant(adminAnyPostDelete))
		Expect(sut.Check(userPostDelete)).To(BeFalse())

		Must(roles.Assign(user, admin))

		Expect(sut.Check(userPostDelete)).To(BeTrue())
	})
})
