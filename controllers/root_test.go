package controllers_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"web/config"
	"web/controllers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
)

var _ = Describe("/", func() {
	var app *fx.App

	BeforeEach(func() {
		app = fx.New(config.Module, controllers.Module)
		app.Start(context.Background())
	})

	AfterEach(func() {
		app.Stop(context.Background())
	})

	It("says hello world", func() {
		res, _ := http.Get("http://localhost:3000/")
		b, _ := ioutil.ReadAll(res.Body)
		Expect(string(b)).To(Equal("Hello, World!"))
	})
})
