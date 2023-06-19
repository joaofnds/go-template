package http_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"app/adapter/featureflags"
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
			featureflags.Module,
			apphttp.FiberModule,
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
		bob := Must2(app.User.CreateUser("bob"))
		found := Must2(app.User.GetUser(bob.Name))

		Expect(found).To(Equal(bob))
	})

	It("lists users", func() {
		bob := Must2(app.User.CreateUser("bob"))
		dave := Must2(app.User.CreateUser("dave"))

		users := Must2(app.User.ListUsers())
		Expect(users).To(Equal([]user.User{bob, dave}))
	})

	It("deletes users", func() {
		bob := Must2(app.User.CreateUser("bob"))
		dave := Must2(app.User.CreateUser("dave"))

		Must(app.User.DeleteUser(dave.Name))

		_, err := app.User.GetUser(dave.Name)
		Expect(err).To(Equal(driver.RequestFailure{
			Status: http.StatusNotFound,
			Body:   "Not Found",
		}))

		users := Must2(app.User.ListUsers())
		Expect(users).To(Equal([]user.User{bob}))
	})

	It("switches feature flag", func() {
		bob := Must2(app.User.CreateUser("bob"))
		bobFeatures := Must2(app.User.GetFeature(bob.Name))
		Expect(bobFeatures).To(Equal(map[string]any{"cool-feature": "on"}))

		frank := Must2(app.User.CreateUser("frank"))
		frankFeatures := Must2(app.User.GetFeature(frank.Name))
		Expect(frankFeatures).To(Equal(map[string]any{"cool-feature": "off"}))
	})
})
