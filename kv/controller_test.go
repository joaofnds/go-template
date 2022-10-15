package kv_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"web/config"
	webhttp "web/http"
	"web/kv"
	"web/test"
	. "web/test/matchers"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("/kv", Ordered, func() {
	var app *fxtest.App
	var url string

	BeforeAll(func() {
		var httpConfig webhttp.Config

		app = fxtest.New(
			GinkgoT(),
			test.NopLogger,
			test.RandomAppConfigPort,
			config.Module,
			webhttp.FiberModule,
			kv.Module,
			fx.Invoke(func(app *fiber.App, controller *kv.Controller) {
				controller.Register(app)
			}),
			fx.Populate(&httpConfig),
		).RequireStart()

		url = fmt.Sprintf("http://localhost:%d/kv", httpConfig.Port)
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