package user_test

import (
	"testing"
	"web/config"
	"web/logger"
	"web/mongo"
	"web/test"
	. "web/test/matchers"
	"web/user"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestUser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Service Test")
}

var _ = Describe("user service", func() {
	var app *fxtest.App
	var userService *user.UserService

	BeforeEach(func() {
		app = fxtest.New(
			GinkgoT(),
			fx.NopLogger,
			logger.Module,
			config.Module,
			fx.Decorate(test.RandomAppConfigPort),
			fx.Decorate(test.NewNopUserInstrumentation),
			mongo.Module,
			user.Module,
			fx.Populate(&userService),
		)
		app.RequireStart()
		Must(userService.DeleteAll())
	})

	AfterEach(func() {
		app.RequireStop()
	})

	Describe("DeleteAll", func() {
		It("removes all users", func() {
			Must2(userService.CreateUser("joao"))
			Expect(userService.List()).NotTo(BeEmpty())

			Must(userService.DeleteAll())
			Expect(userService.List()).To(BeEmpty())
		})
	})

	It("created users can be found by name", func() {
		user := Must2(userService.CreateUser("joao"))

		found := Must2(userService.FindByName(user.Name))
		Expect(found).To(Equal(user))
	})

	It("created users appear on users listing", func() {
		user := Must2(userService.CreateUser("joao"))
		Expect(userService.List()).To(ContainElement(user))
	})

	It("removed users do not appear on users listing", func() {
		user := Must2(userService.CreateUser("joao"))
		Must(userService.Remove(user))

		users := Must2(userService.List())
		Expect(users).NotTo(ContainElement(user))
	})
})
