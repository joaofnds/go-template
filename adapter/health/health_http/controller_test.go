package health_http_test

import (
	"app/adapter/health"
	"app/test/driver"
	"app/test/harness"

	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
)

func TestHealth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "/health suite")
}

var _ = Describe("/health", Ordered, func() {
	var app *harness.Harness
	var api *driver.Driver
	var checker *health.FakeHealthService

	BeforeAll(func() {
		app = harness.Setup(
			harness.WithFxOptions(
				fx.Decorate(func(original health.Checker) health.Checker {
					checker = health.NewFakeHealthService(original)
					return checker
				}),
			),
		)
		api = app.NewDriver()
	})
	AfterAll(func() { app.Teardown() })

	Context("healthy", func() {
		It("returns status OK", func() {
			resp := api.Health.MustGetReq()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("checks postgres connection", func() {
			Expect(api.Health.MustGet()).To(ContainSubstring(`"postgres":{"status":"up"}`))
		})

		It("checks redis connection", func() {
			Expect(api.Health.MustGet()).To(ContainSubstring(`"redis":{"status":"up"}`))
		})
	})

	Context("unhealthy", func() {
		BeforeAll(func() { checker.UseUnhealthy() })

		It("returns status service unavailable", func() {
			resp := api.Health.MustGetReq()
			Expect(resp.StatusCode).To(Equal(http.StatusServiceUnavailable))
		})

		It("checks postgres connection", func() {
			Expect(api.Health.MustGetFailed()).To(ContainSubstring(`"postgres":{"status":"down"}`))
		})

		It("checks redis connection", func() {
			Expect(api.Health.MustGetFailed()).To(ContainSubstring(`"redis":{"status":"down"}`))
		})
	})
})
