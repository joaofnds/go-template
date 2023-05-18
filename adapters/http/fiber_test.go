package http_test

import (
	apphttp "app/adapters/http"
	"app/adapters/logger"
	"app/config"
	"app/test"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	. "app/test/matchers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var PanicHandler = fx.Invoke(func(app *fiber.App) {
	app.All("panic", func(c *fiber.Ctx) error {
		panic("panic handler")
	})
})

func TestFiber(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "fiber suite")
}

var _ = Describe("fiber middlewares", func() {
	var (
		fxApp      *fxtest.App
		url        string
		httpConfig apphttp.Config
	)

	BeforeEach(func() {

		fxApp = fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			test.RandomAppConfigPort,
			apphttp.NopProbeProvider,
			config.Module,
			apphttp.Module,
			PanicHandler,
			fx.Populate(&httpConfig),
		).RequireStart()

		url = fmt.Sprintf("http://localhost:%d", httpConfig.Port)
	})

	AfterEach(func() {
		fxApp.RequireStop()
	})

	It("recovers from panic", func() {
		req := Must2(http.Get(url + "/panic"))
		Expect(req.StatusCode).To(Equal(http.StatusInternalServerError))

		req = Must2(http.Get(url + "/somethingelse"))
		Expect(req.StatusCode).To(Equal(http.StatusNotFound))
	})

	It("limits requests", func() {
		for i := 0; i < httpConfig.Limiter.Requests; i++ {
			req := Must2(http.Get(url + "/somethingelse"))
			Expect(req.StatusCode).To(Equal(http.StatusNotFound))
		}

		req := Must2(http.Get(url + "/somethingelse"))
		Expect(req.StatusCode).To(Equal(http.StatusTooManyRequests))
	})
})
