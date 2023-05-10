package metrics_test

import (
	"app/adapters/logger"
	"app/adapters/metrics"
	"app/config"
	. "app/test/matchers"
	"fmt"
	"io"
	"net/http"
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHealth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "/health suite")
}

var _ = Describe("/", Ordered, func() {
	var url string

	BeforeAll(func() {
		var cfg metrics.Config
		fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			config.Module,
			metrics.Module,
			fx.Populate(&cfg),
		).RequireStart()
		url = fmt.Sprintf("http://%s/metrics", cfg.Addr)
	})

	It("returns status OK", func() {
		res := Must2(http.Get(url))
		Expect(res.StatusCode).To(Equal(http.StatusOK))
	})

	It("returns how many requests were made", func() {
		res := Must2(http.Get(url))
		b := Must2(io.ReadAll(res.Body))
		Expect(b).To(ContainSubstring(`promhttp_metric_handler_requests_total{code="200"} 1`))
	})
})
