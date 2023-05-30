package health_test

import (
	"app/adapter/health"
	apphttp "app/adapter/http"
	"app/adapter/logger"
	"app/adapter/postgres"
	"app/adapter/redis"
	"app/config"
	"app/test"
	. "app/test/matchers"

	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestHealth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "/health suite")
}

var _ = Describe("/health", Ordered, func() {
	var app *fxtest.App
	var url string

	Context("healthy", func() {
		BeforeAll(func() {
			var cfg apphttp.Config
			app = fxtest.New(
				GinkgoT(),
				logger.NopLoggerProvider,
				apphttp.NopProbeProvider,
				test.AvailablePortProvider,
				config.Module,
				redis.Module,
				postgres.Module,
				apphttp.FiberModule,
				health.Module,
				fx.Invoke(func(app *fiber.App, controller *health.Controller) {
					controller.Register(app)
				}),
				fx.Populate(&cfg),
			)
			url = fmt.Sprintf("http://localhost:%d/health", cfg.Port)
			app.RequireStart()
		})

		AfterAll(func() { app.RequireStop() })

		It("returns status OK", func() {
			res := Must2(http.Get(url))
			Expect(res.StatusCode).To(Equal(http.StatusOK))
		})

		It("checks postgres connection", func() {
			res := Must2(http.Get(url))
			body := Must2(io.ReadAll(res.Body))
			Expect(body).To(ContainSubstring(`"postgres":{"status":"up"}`))
		})

		It("checks redis connection", func() {
			res := Must2(http.Get(url))
			body := Must2(io.ReadAll(res.Body))
			Expect(body).To(ContainSubstring(`"redis":{"status":"up"}`))
		})
	})

	Context("unhealthy", func() {
		BeforeEach(func() {
			var cfg apphttp.Config
			app = fxtest.New(
				GinkgoT(),
				logger.NopLoggerProvider,
				apphttp.NopProbeProvider,
				test.AvailablePortProvider,
				config.Module,
				apphttp.FiberModule,
				health.Module,
				fx.Invoke(func(app *fiber.App, controller *health.Controller) {
					controller.Register(app)
				}),
				fx.Decorate(func() health.Checker {
					return health.NewUnhealthyHealthService()
				}),
				fx.Populate(&cfg),
			)
			url = fmt.Sprintf("http://localhost:%d/health", cfg.Port)
			app.RequireStart()
		})

		AfterEach(func() { app.RequireStop() })

		It("returns status OK", func() {
			res := Must2(http.Get(url))
			Expect(res.StatusCode).To(Equal(http.StatusServiceUnavailable))
		})

		It("checks postgres connection", func() {
			res := Must2(http.Get(url))
			body := Must2(io.ReadAll(res.Body))
			Expect(body).To(ContainSubstring(`"postgres":{"status":"down"}`))
		})

		It("checks redis connection", func() {
			res := Must2(http.Get(url))
			body := Must2(io.ReadAll(res.Body))
			Expect(body).To(ContainSubstring(`"redis":{"status":"down"}`))
		})
	})
})
