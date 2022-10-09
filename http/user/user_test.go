package user_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"web/config"
	"web/http/fiber"
	http_user "web/http/user"
	"web/mongo"
	"web/test"
	. "web/test/matchers"
	"web/user"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestUser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User HTTP Suite")
}

var _ = Describe("/users", Ordered, func() {
	var app *fxtest.App
	var userService *user.UserService
	var url string

	BeforeAll(func() {
		var appConfig config.AppConfig

		app = fxtest.New(
			GinkgoT(),
			test.NopLogger,
			test.RandomAppConfigPort,
			config.Module,
			fiber.Module,
			mongo.Module,
			user.Module,
			http_user.Providers,
			fx.Populate(&appConfig, &userService),
		).RequireStart()

		url = fmt.Sprintf("http://localhost:%d/users", appConfig.Port)
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
