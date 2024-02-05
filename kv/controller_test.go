package kv_test

import (
	apphttp "app/adapter/http"
	"app/adapter/logger"
	"app/adapter/redis"
	"app/config"
	"app/kv"
	"app/test"
	"app/test/driver"
	"net/http"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("/kv", Ordered, func() {
	var fxApp *fxtest.App
	var app *driver.Driver

	BeforeAll(func() {
		var httpConfig apphttp.Config

		fxApp = fxtest.New(
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
			driver.Provider,
			fx.Populate(&app, &httpConfig),
		).RequireStart()
	})

	AfterAll(func() {
		fxApp.RequireStop()
	})

	Context("GET", func() {
		It("returns the value under the key", func() {
			app.KV.MustSet("foo", "bar")

			val := app.KV.MustGet("foo")

			Expect(val).To(Equal("bar"))
		})
	})

	Context("POST", func() {
		It("responds with status created", func() {
			res, _ := app.KV.SetReq("foo", "bar")
			Expect(res.StatusCode).To(Equal(http.StatusCreated))
		})

		It("sets the value to the key", func() {
			app.KV.MustSet("foo", "bar")

			val := app.KV.MustGet("foo")

			Expect(val).To(Equal("bar"))
		})
	})

	Context("DELETE", func() {
		It("forgets the key", func() {
			app.KV.MustSet("bar", "foo")
			app.KV.MustDel("bar")

			res, _ := app.KV.GetReq("bar")
			Expect(res.StatusCode).To(Equal(http.StatusNotFound))
		})
	})
})
