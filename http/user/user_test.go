package user_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"web/config"
	"web/mongo"
	"web/http/fiber"
	http_user "web/http/user"
	"web/logger"
	"web/test"
	"web/user"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("/users", Ordered, func() {
	var app *fxtest.App
	var userService *user.UserService
	var url string

	BeforeAll(func() {
		var appConfig config.AppConfig

		app = fxtest.New(
			GinkgoT(),
			logger.Module,
			config.Module,
			fx.Decorate(test.RandomAppConfigPort),
			fiber.Module,
			mongo.Module,
			user.Module,
			http_user.Providers,
			fx.Populate(&userService),
			fx.Populate(&appConfig),
		).RequireStart()

		url = fmt.Sprintf("http://localhost:%d/users", appConfig.Port)
	})

	AfterAll(func() {
		app.RequireStop()
	})

	BeforeEach(func() {
		userService.DeleteAll()
	})

	Context("GET", func() {
		It("concats all user's names", func() {
			userService.CreateUser("joao")
			userService.CreateUser("fernandes")

			res, err := http.Get(url)
			Expect(err).To(BeNil())
			Expect(res.StatusCode).To(Equal(http.StatusOK))

			b, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			Expect(string(b)).To(Equal("joaofernandes"))
		})
	})

	Context("POST", func() {
		It("adds the user", func() {
			body := bytes.NewBufferString(`{"name": "joao"}`)
			res, _ := http.Post(url, "application/json", body)
			Expect(res.StatusCode).To(Equal(http.StatusCreated))

			body = bytes.NewBufferString(`{"name": "vitor"}`)
			res, _ = http.Post(url, "application/json", body)
			Expect(res.StatusCode).To(Equal(http.StatusCreated))

			res, _ = http.Get(url)
			b, _ := ioutil.ReadAll(res.Body)

			Expect(string(b)).To(Equal("joaovitor"))
		})
	})
})
