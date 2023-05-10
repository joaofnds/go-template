package user_test

import (
	apphttp "app/adapters/http"
	"app/adapters/logger"
	"app/adapters/mongo"
	"app/config"
	"app/test"
	. "app/test/matchers"
	"app/user"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("/users", Ordered, func() {
	var app *fxtest.App
	var userService *user.Service
	var url string

	BeforeAll(func() {
		var httpConfig apphttp.Config

		app = fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			test.RandomAppConfigPort,
			config.Module,
			apphttp.FiberModule,
			mongo.Module,
			user.Module,
			fx.Invoke(func(app *fiber.App, controller *user.Controller) {
				controller.Register(app)
			}),
			fx.Populate(&httpConfig, &userService),
		).RequireStart()

		url = fmt.Sprintf("http://localhost:%d/users", httpConfig.Port)
	})

	AfterAll(func() {
		app.RequireStop()
	})

	BeforeEach(func() {
		Must(userService.DeleteAll())
	})

	Context("GET", func() {
		It("concats all user's names", func() {
			Must2(userService.CreateUser("joao"))
			Must2(userService.CreateUser("fernandes"))

			res := Must2(http.Get(url))
			Expect(res.StatusCode).To(Equal(http.StatusOK))

			b := Must2(io.ReadAll(res.Body))

			Expect(string(b)).To(Equal("joaofernandes"))
		})
	})

	Context("POST", func() {
		It("adds the user", func() {
			body := bytes.NewBufferString(`{"name": "joao"}`)
			res := Must2(http.Post(url, "application/json", body))
			Expect(res.StatusCode).To(Equal(http.StatusCreated))

			body = bytes.NewBufferString(`{"name": "vitor"}`)
			res = Must2(http.Post(url, "application/json", body))
			Expect(res.StatusCode).To(Equal(http.StatusCreated))

			res, _ = http.Get(url)
			b := Must2(io.ReadAll(res.Body))

			Expect(string(b)).To(Equal("joaovitor"))
		})
	})
})
