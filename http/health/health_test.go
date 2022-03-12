package health_test

import (
	"net/http"
	"testing"
	"web/config"
	"web/http/fiber"
	"web/http/health"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestHealth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "/health suite")
}

var _ = Describe("/", func() {
	var app *fxtest.App

	BeforeEach(func() {
		app = fxtest.New(
			GinkgoT(),
			fx.NopLogger,
			config.Module,
			fiber.Module,
			health.Providers,
		)
		app.RequireStart()
	})

	AfterEach(func() {
		app.RequireStop()
	})

	It("returns status OK", func() {
		res, _ := http.Get("http://localhost:3000/health")
		Expect(res.StatusCode).To(Equal(http.StatusOK))
	})
})
