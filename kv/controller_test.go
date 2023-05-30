package kv_test

import (
	apphttp "app/adapter/http"
	"app/adapter/logger"
	"app/adapter/redis"
	"app/config"
	"app/kv"
	"app/test"
	. "app/test/matchers"
	"app/test/req"
	"fmt"
	"io"
	"net/http"

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
		var httpConfig apphttp.Config

		app = fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			config.Module,
			test.AvailablePortProvider,
			apphttp.FiberModule,
			redis.Module,
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
			Must2(req.Post(url+"/foo/bar", nil, nil))
			res := Must2(req.Get(url+"/foo", nil))

			b := Must2(io.ReadAll(res.Body))
			Expect(string(b)).To(Equal("bar"))
		})
	})

	Context("POST", func() {
		It("responds with status created", func() {
			res := Must2(req.Post(url+"/foo/bar", nil, nil))
			Expect(res.StatusCode).To(Equal(http.StatusCreated))
		})

		It("sets the value to the key", func() {
			Must2(req.Post(url+"/foo/bar", nil, nil))

			res := Must2(req.Get(url+"/foo", nil))
			b := Must2(io.ReadAll(res.Body))
			Expect(string(b)).To(Equal("bar"))
		})
	})

	Context("DELETE", func() {
		It("forgets the key", func() {
			res := Must2(req.Post(url+"/bar/foo", nil, nil))
			Expect(res.StatusCode).To(Equal(http.StatusCreated))

			res = Must2(req.Delete(url+"/bar", nil))
			Expect(res.StatusCode).To(Equal(http.StatusOK))

			res = Must2(req.Get(url+"/bar", nil))
			Expect(res.StatusCode).To(Equal(http.StatusNotFound))
		})
	})
})
