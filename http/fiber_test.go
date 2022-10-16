package http_test

import (
	"app/config"
	apphttp "app/http"
	"app/test"
	"fmt"
	"net/http"
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	. "app/test/matchers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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
			test.NopLogger,
			test.RandomAppConfigPort,
			test.NopHTTPInstrumentation,
			test.PanicHandler,
			config.Module,
			apphttp.FiberModule,
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
