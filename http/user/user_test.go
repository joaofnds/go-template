package user_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"web/config"
	"web/http/fiber"
	http_user "web/http/user"
	"web/user"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("/users", func() {
	var app *fxtest.App
	var us *user.UserService

	BeforeEach(func() {
		app = fxtest.New(
			GinkgoT(),
			fx.NopLogger,
			config.Module,
			fiber.Module,
			http_user.Providers,
			user.Module,
			fx.Populate(&us),
		)
		app.RequireStart()
		us.DeleteAll()
	})

	AfterEach(func() {
		app.RequireStop()
	})

	Context("GET", func() {
		It("concats all user's names", func() {
			us.CreateUser("joao")
			us.CreateUser("fernandes")

			res, _ := http.Get("http://localhost:3000/users")
			b, _ := ioutil.ReadAll(res.Body)

			Expect(string(b)).To(Equal("joaofernandes"))
		})
	})

	Context("POST", func() {
		It("adds the user", func() {
			body := bytes.NewBufferString(`{"name": "joao"}`)
			res, _ := http.Post("http://localhost:3000/users", "application/json", body)
			Expect(res.StatusCode).To(Equal(http.StatusCreated))

			body = bytes.NewBufferString(`{"name": "vitor"}`)
			res, _ = http.Post("http://localhost:3000/users", "application/json", body)
			Expect(res.StatusCode).To(Equal(http.StatusCreated))

			res, _ = http.Get("http://localhost:3000/users")
			b, _ := ioutil.ReadAll(res.Body)

			Expect(string(b)).To(Equal("joaovitor"))
		})
	})
})
