package authzcasbin_test

import (
	"app/adapter/authzcasbin"
	"app/adapter/logger"
	"app/adapter/postgres"
	"app/config"
	"app/internal/ref"
	"app/test"
	. "app/test/matchers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("casbin role manager", func() {
	user := ref.New("111", "user")
	admin := ref.New("admin", "role")

	var app *fxtest.App
	var sut *authzcasbin.RoleManager

	BeforeEach(func() {
		app = fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			test.CasbinStringAdapter,
			config.Module,
			postgres.Module,
			authzcasbin.Module,
			fx.Populate(&sut),
		)
		app.RequireStart()
	})

	AfterEach(func() {
		app.RequireStop()
	})

	Describe("roles", func() {
		It("grants roles", func() {
			Must(sut.Assign(user, admin))
		})

		It("retrieves roles", func() {
			Must(sut.Assign(user, admin))

			Expect(sut.GetAll(user)).To(ConsistOf(admin))
		})

		It("revokes roles", func() {
			Must(sut.Assign(user, admin))

			Must(sut.Revoke(user, admin))

			Expect(sut.GetAll(user)).To(BeEmpty())
		})
	})
})
