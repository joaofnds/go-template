package http_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	apphttp "app/adapter/http"
	"app/adapter/logger"
	"app/adapter/postgres"
	"app/adapter/validation"
	"app/config"
	"app/test"
	"app/test/driver"
	. "app/test/matchers"
	"app/user"
	userhttp "app/user/http"
	usermodule "app/user/module"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestUserHTTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "user http suite")
}

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
			test.Queue,
			test.Transaction,
			test.AvailablePortProvider,
			config.Module,
			apphttp.FiberProvider,
			validation.Module,
			postgres.Module,
			usermodule.Module,
			fx.Invoke(func(app *fiber.App, controller *userhttp.Controller) {
				controller.Register(app)
			}),
			fx.Populate(&httpConfig, &userService),
		).RequireStart()

		app = driver.NewDriver(fmt.Sprintf("http://localhost:%d", httpConfig.Port))
	})

	BeforeEach(func() { Must(userService.DeleteAll(context.Background())) })

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
