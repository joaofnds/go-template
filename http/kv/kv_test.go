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
	"web/test"
	. "web/test/matchers"

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
			test.NopLogger,
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
			Must2(http.Post(url+"/foo/bar", "", strings.NewReader("")))
			res := Must2(http.Get(url + "/foo"))

			b := Must2(io.ReadAll(res.Body))
			Expect(string(b)).To(Equal("bar"))
		})
	})

	Context("POST", func() {
		It("responds with status created", func() {
			res := Must2(http.Post(url+"/foo/bar", "", strings.NewReader("")))
			Expect(res.StatusCode).To(Equal(http.StatusCreated))
		})

		It("sets the value to the key", func() {
			Must2(http.Post(url+"/foo/bar", "", strings.NewReader("")))

			res := Must2(http.Get(url + "/foo"))
			b := Must2(io.ReadAll(res.Body))
			Expect(string(b)).To(Equal("bar"))
		})
	})

	Context("DELETE", func() {
		It("adds the user", func() {
			res := Must2(http.Post(url+"/foo/bar", "", strings.NewReader("")))
			Expect(res.StatusCode).To(Equal(http.StatusCreated))

			res = Must2(http.Get(url))
			Expect(res.StatusCode).To(Equal(http.StatusNotFound))

			req := Must2(http.NewRequest(http.MethodDelete, url+"/foo/bar", strings.NewReader("")))
			Must2(http.DefaultClient.Do(req))

			res = Must2(http.Get(url))
			Expect(res.StatusCode).To(Equal(http.StatusNotFound))
		})
	})
})
