package user_test

import (
	"fmt"
	"net/http"

	apphttp "app/adapters/http"
	"app/adapters/logger"
	"app/adapters/postgres"
	"app/config"
	"app/test"
	"app/test/driver"
	"app/user"

	. "app/test/matchers"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("/users", Ordered, func() {
	var (
		app         *driver.Driver
		fxApp       *fxtest.App
		userService *user.Service
	)

	BeforeAll(func() {
		var httpConfig apphttp.Config

		fxApp = fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			test.QueueProvider,
			config.Module,
			apphttp.FiberProvider,
			postgres.Module,
			user.Module,
			fx.Invoke(func(app *fiber.App, controller *user.Controller) {
				controller.Register(app)
			}),
			fx.Populate(&httpConfig, &userService),
		).RequireStart()

		app = driver.NewDriver(fmt.Sprintf("http://localhost:%d", httpConfig.Port))
	})

	BeforeEach(func() { Must(userService.DeleteAll()) })

	AfterAll(func() { fxApp.RequireStop() })

	It("creates and gets user", func() {
		bob := Must2(app.CreateUser("bob"))
		found := Must2(app.GetUser(bob.Name))

		Expect(found).To(Equal(bob))
	})

	It("lists users", func() {
		bob := Must2(app.CreateUser("bob"))
		dave := Must2(app.CreateUser("dave"))

		users := Must2(app.ListUsers())
		Expect(users).To(Equal([]user.User{bob, dave}))
	})

	It("deletes users", func() {
		bob := Must2(app.CreateUser("bob"))
		dave := Must2(app.CreateUser("dave"))

		Must(app.DeleteUser(dave.Name))

		res := Must2(app.API.DeleteUser(dave.Name))
		Expect(res.StatusCode).To(Equal(http.StatusNotFound))

		users := Must2(app.ListUsers())
		Expect(users).To(Equal([]user.User{bob}))
	})
})
