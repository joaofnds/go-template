package user_test

import (
	"web/config"
	"web/logger"
	"web/mongo"
	"web/test"
	"web/user"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

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
			mongo.Module,
			user.Module,
			fx.Populate(&userService),
		)
		app.RequireStart()
		userService.DeleteAll()
	})

	AfterEach(func() {
		app.RequireStop()
	})

	Describe("DeleteAll", func() {
		It("removes all users", func() {
			userService.CreateUser("joao")
			Expect(userService.List()).NotTo(BeEmpty())

			userService.DeleteAll()
			Expect(userService.List()).To(BeEmpty())
		})
	})

	It("created users can be found by name", func() {
		user, err := userService.CreateUser("joao")
		Expect(err).To(BeNil())

		found, err := userService.FindByName(user.Name)
		Expect(err).To(BeNil())
		Expect(found).To(Equal(user))
	})

	It("created users appear on users listing", func() {
		user, err := userService.CreateUser("joao")
		Expect(err).To(BeNil())
		Expect(userService.List()).To(ContainElement(user))
	})

	It("removed users do not appear on users listing", func() {
		user, err := userService.CreateUser("joao")
		Expect(err).To(BeNil())
		userService.Remove(user)

		users, err := userService.List()
		Expect(err).To(BeNil())
		Expect(users).NotTo(ContainElement(user))
	})
})
