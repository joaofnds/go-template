package http_test

import (
	"context"
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
	"app/test/matchers"
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
		fxApp = fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			test.Queue,
			test.Transaction,
			test.AvailablePortProvider,
			driver.Provider,
			config.Module,
			featureflags.Module,
			apphttp.FiberModule,
			validation.Module,
			postgres.Module,
			usermodule.Module,
			fx.Invoke(func(fiberapp *fiber.App, controller *userhttp.Controller) {
				controller.Register(fiberapp)
			}),
			fx.Populate(&app, &userService),
		).RequireStart()
	})

	BeforeEach(func() { matchers.Must(userService.DeleteAll(context.Background())) })

	AfterAll(func() { fxApp.RequireStop() })

	It("creates and gets user", func() {
		bob := app.User.MustCreateUser("bob")
		found := app.User.MustGetUser(bob.Name)

		Expect(found).To(Equal(bob))
	})

	It("lists users", func() {
		bob := app.User.MustCreateUser("bob")
		dave := app.User.MustCreateUser("dave")

		users := app.User.MustListUsers()
		Expect(users).To(ConsistOf(bob, dave))
	})

	It("deletes users", func() {
		bob := app.User.MustCreateUser("bob")
		dave := app.User.MustCreateUser("dave")

		app.User.MustDeleteUser(dave.Name)

		_, err := app.User.GetUser(dave.Name)
		Expect(err).To(Equal(driver.RequestFailure{
			Status: http.StatusNotFound,
			Body:   "Not Found",
		}))

		users := app.User.MustListUsers()
		Expect(users).To(ConsistOf(bob))
	})

	It("switches feature flag", func() {
		bob := app.User.MustCreateUser("bob")
		bobFeatures := app.User.MustGetFeature(bob.Name)
		Expect(bobFeatures).To(HaveKeyWithValue("cool-feature", "on"))

		frank := app.User.MustCreateUser("frank")
		frankFeatures := app.User.MustGetFeature(frank.Name)
		Expect(frankFeatures).To(HaveKeyWithValue("cool-feature", "off"))
	})
})
