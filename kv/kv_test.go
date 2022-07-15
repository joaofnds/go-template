package kv_test

import (
	"testing"
	"web/kv"
	. "web/test/matchers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestUser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KV Test")
}

var _ = Describe("kv service", func() {
	var app *fxtest.App
	var store *kv.KeyValStore

	BeforeEach(func() {
		app = fxtest.New(
			GinkgoT(),
			fx.NopLogger,
			kv.Module,
			fx.Populate(&store),
		)
		app.RequireStart()
	})

	AfterEach(func() {
		app.RequireStop()
	})

	It("can retrieve values", func() {
		Must(store.Set("foo", "bar"))

		val := Must2(store.Get("foo"))

		Expect(val).To(Equal("bar"))
	})

	It("cannot get values deleted", func() {
		Must(store.Set("foo", "bar"))
		Must(store.Del("foo"))

		_, err := store.Get("foo")
		Expect(err).To(Equal(kv.ErrNotFound))
	})
})
