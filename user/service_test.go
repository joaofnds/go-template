package user_test

import (
	"context"
	"testing"

	"app/adapter/logger"
	"app/adapter/postgres"
	"app/config"
	"app/test"
	. "app/test/matchers"
	"app/user"
	usermodule "app/user/module"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestUserService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "user service suite")
}

var _ = Describe("user service", func() {
	var app *fxtest.App
	var userService *user.Service

	BeforeEach(func() {
		app = fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			test.Queue,
			test.Transaction,
			config.Module,
			postgres.Module,
			usermodule.Module,
			fx.Populate(&userService),
		)
		app.RequireStart()
		Must(userService.DeleteAll(context.Background()))
	})

	AfterEach(func() {
		app.RequireStop()
	})

	Describe("DeleteAll", func() {
		It("removes all users", func(ctx SpecContext) {
			Must2(userService.CreateUser(ctx, "joao"))
			Expect(userService.List(ctx)).NotTo(BeEmpty())

			Must(userService.DeleteAll(ctx))
			Expect(userService.List(ctx)).To(BeEmpty())
		})
	})

	It("created users can be found by name", func(ctx SpecContext) {
		user := Must2(userService.CreateUser(ctx, "joao"))

		found := Must2(userService.FindByName(ctx, user.Name))
		Expect(found).To(Equal(user))
	})

	It("created users appear on users listing", func(ctx SpecContext) {
		user := Must2(userService.CreateUser(ctx, "joao"))
		Expect(userService.List(ctx)).To(ContainElement(user))
	})

	It("removed users do not appear on users listing", func(ctx SpecContext) {
		user := Must2(userService.CreateUser(ctx, "joao"))
		Must(userService.Remove(ctx, user))

		users := Must2(userService.List(ctx))
		Expect(users).NotTo(ContainElement(user))
	})
})
