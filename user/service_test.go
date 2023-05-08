package user_test

import (
	"app/config"
	"app/storage/mongo"
	"app/test"
	. "app/test/matchers"
	"app/user"
	"testing"

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
	var userService *user.Service

	BeforeEach(func() {
		app = fxtest.New(
			GinkgoT(),
			test.NopLogger,
			test.RandomAppConfigPort,
			test.NopUserInstrumentation,
			config.Module,
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
