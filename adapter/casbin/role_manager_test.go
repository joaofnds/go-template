package casbin_test

import (
	"app/adapter/casbin"
	"app/adapter/logger"
	"app/adapter/postgres"
	"app/adapter/validation"
	"app/config"
	"app/internal/ref"
	. "app/test/matchers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("casbin role manager", func() {
	user := ref.New("user", "111")
	admin := ref.New("role", "admin")

	var app *fxtest.App
	var sut *casbin.RoleManager

	BeforeEach(func() {
		app = fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			validation.Module,
			config.Module,
			postgres.Module,
			casbin.Module,
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
