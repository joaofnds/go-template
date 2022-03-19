package kv_test

import (
	"testing"
	"web/kv"

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
		err := store.Set("foo", "bar")
		Expect(err).To(BeNil())

		val, err := store.Get("foo")
		Expect(err).To(BeNil())

		Expect(val).To(Equal("bar"))
	})

	It("cannot get values deleted", func() {
		store.Set("foo", "bar")
		store.Del("foo")

		_, err := store.Get("foo")
		Expect(err).To(Equal(kv.ErrNotFound))
	})
})
