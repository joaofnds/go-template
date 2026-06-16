package casbin_test

import (
	"app/adapter/casbin"
	"app/adapter/logger"
	"app/adapter/postgres"
	"app/adapter/validation"
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
	user := ref.New("user", "111")
	admin := ref.New("role", "admin")
	customer := ref.New("role", "customer")

	var (
		app    *fxtest.App
		txPool *test.TxPool
		sut    *casbin.RoleManager
	)

	BeforeEach(func() {
		app = fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			test.TxIsolation,
			validation.Module,
			config.Module,
			postgres.Module,
			casbin.Module,
			fx.Populate(&txPool, &sut),
		)
		app.RequireStart()
		Must(txPool.Begin())
	})

	AfterEach(func() {
		Must(txPool.Rollback())
		app.RequireStop()
	})

	It("grants roles", func() {
		Must(sut.Assign(user, admin))
	})

	It("retrieves roles", func() {
		Must(sut.Assign(user, admin))
		Must(sut.Assign(user, customer))

		Expect(sut.GetAll(user)).To(ConsistOf(admin, customer))
	})

	It("revokes roles", func() {
		Must(sut.Assign(user, admin))

		Must(sut.Revoke(user, admin))

		Expect(sut.GetAll(user)).To(BeEmpty())
	})

	It("removes all roles for a user", func() {
		Must(sut.Assign(user, admin))
		Must(sut.Assign(user, customer))

		Expect(sut.GetAll(user)).To(ConsistOf(admin, customer))

		Must(sut.RevokeAll(user))

		Expect(sut.GetAll(user)).To(BeEmpty())
	})
})
