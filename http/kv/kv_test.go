package kv_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"web/config"
	"web/http/fiber"
	http_kv "web/http/kv"
	"web/kv"
	"web/logger"
	"web/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestUser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User HTTP Suite")
}

var _ = Describe("/kv", Ordered, func() {
	var app *fxtest.App
	var url string

	BeforeAll(func() {
		var appConfig config.AppConfig

		app = fxtest.New(
			GinkgoT(),
			fx.NopLogger,
			logger.Module,
			config.Module,
			fx.Decorate(test.RandomAppConfigPort),
			fiber.Module,
			kv.Module,
			http_kv.Providers,
			fx.Populate(&appConfig),
		).RequireStart()

		url = fmt.Sprintf("http://localhost:%d/kv", appConfig.Port)
	})

	AfterAll(func() {
		app.RequireStop()
	})

	Context("GET", func() {
		It("returns the value under the key", func() {
			http.Post(url+"/foo/bar", "", strings.NewReader(""))
			res, _ := http.Get(url + "/foo")

			b, _ := io.ReadAll(res.Body)
			Expect(string(b)).To(Equal("bar"))
		})
	})

	Context("POST", func() {
		It("responds with status created", func() {
			res, _ := http.Post(url+"/foo/bar", "", strings.NewReader(""))
			Expect(res.StatusCode).To(Equal(http.StatusCreated))
		})

		It("sets the value to the key", func() {
			http.Post(url+"/foo/bar", "", strings.NewReader(""))

			res, _ := http.Get(url + "/foo")
			b, _ := io.ReadAll(res.Body)
			Expect(string(b)).To(Equal("bar"))
		})
	})

	Context("DELETE", func() {
		It("adds the user", func() {
			res, _ := http.Post(url+"/foo/bar", "", strings.NewReader(""))
			Expect(res.StatusCode).To(Equal(http.StatusCreated))

			res, _ = http.Get(url)
			Expect(res.StatusCode).To(Equal(http.StatusNotFound))

			req, _ := http.NewRequest(http.MethodDelete, url+"/foo/bar", strings.NewReader(""))
			http.DefaultClient.Do(req)

			res, _ = http.Get(url)
			Expect(res.StatusCode).To(Equal(http.StatusNotFound))
		})
	})
})
