package health_test

import (
	"web/config"
	"web/health"
	webhttp "web/http"
	"web/kv"
	"web/mongo"
	"web/test"
	testhealth "web/test/health"
	. "web/test/matchers"

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

var _ = Describe("/health", func() {
	var app *fxtest.App
	var url string

	Context("healthy", func() {
		BeforeEach(func() {
			var cfg webhttp.Config
			app = fxtest.New(
				GinkgoT(),
				test.NopLogger,
				test.RandomAppConfigPort,
				test.NopHTTPInstrumentation,
				config.Module,
				webhttp.FiberModule,
				health.Module,
				mongo.Module,
				kv.Module,
				fx.Invoke(func(app *fiber.App, controller *health.Controller) {
					controller.Register(app)
				}),
				fx.Populate(&cfg),
			)
			url = fmt.Sprintf("http://localhost:%d/health", cfg.Port)
			app.RequireStart()
		})

		AfterEach(func() { app.RequireStop() })

		It("returns status OK", func() {
			res := Must2(http.Get(url))
			Expect(res.StatusCode).To(Equal(http.StatusOK))
		})

		It("checks mongo connection", func() {
			res := Must2(http.Get(url))
			b := Must2(io.ReadAll(res.Body))
			Expect(b).To(ContainSubstring(`"mongo":{"status":"up"}`))
		})

		It("checks redis connection", func() {
			res := Must2(http.Get(url))
			b := Must2(io.ReadAll(res.Body))
			Expect(b).To(ContainSubstring(`"redis":{"status":"up"}`))
		})
	})

	Context("unhealthy", func() {
		BeforeEach(func() {
			var cfg webhttp.Config
			app = fxtest.New(
				GinkgoT(),
				test.NopLogger,
				test.RandomAppConfigPort,
				test.NopHTTPInstrumentation,
				testhealth.UnhealthyHealthService,
				config.Module,
				mongo.Module,
				kv.Module,
				health.Module,
				webhttp.FiberModule,
				fx.Invoke(func(app *fiber.App, controller *health.Controller) {
					controller.Register(app)
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

		It("checks mongo connection", func() {
			res := Must2(http.Get(url))
			b := Must2(io.ReadAll(res.Body))
			Expect(b).To(ContainSubstring(`"mongo":{"status":"down"}`))
		})

		It("checks redis connection", func() {
			res := Must2(http.Get(url))
			b := Must2(io.ReadAll(res.Body))
			Expect(b).To(ContainSubstring(`"redis":{"status":"down"}`))
		})
	})
})
