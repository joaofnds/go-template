package fiber_test

import (
	"fmt"
	"net/http"

	"web/config"
	"web/http/fiber"
	. "web/http/req"
	"web/test"
	. "web/test/matchers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("fiber", func() {
	var (
		fxApp *fxtest.App
		url   string
	)

	BeforeEach(func() {
		var cfg config.AppConfig

		fxApp = fxtest.New(
			GinkgoT(),
			test.NopLogger,
			test.RandomAppConfigPort,
			test.NopHTTPInstrumentation,
			test.PanicHandler,
			config.Module,
			fiber.Module,
			fx.Populate(&cfg),
		).RequireStart()

		url = fmt.Sprintf("http://localhost:%d", cfg.Port)
	})

	AfterEach(func() {
		fxApp.RequireStop()
	})

	It("recovers from panic", func() {
		req := Must2(Get(url+"/panic", nil))
		Expect(req.StatusCode).To(Equal(http.StatusInternalServerError))

		req = Must2(Get(url+"/somethingelse", nil))
		Expect(req.StatusCode).To(Equal(http.StatusNotFound))
	})

	It("limits requests", func() {
		for i := 0; i < 30; i++ {
			req := Must2(Get(url+"/somethingelse", nil))
			Expect(req.StatusCode).To(Equal(http.StatusNotFound))
		}

		req := Must2(Get(url+"/somethingelse", nil))
		Expect(req.StatusCode).To(Equal(http.StatusTooManyRequests))
	})
})
