package kv_test

import (
	"app/adapters/logger"
	"app/adapters/redis"
	"app/config"
	"app/kv"
	. "app/test/matchers"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestUser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KV Test")
}

var _ = Describe("kv store", func() {
	var app *fxtest.App
	var store *kv.RedisStore

	BeforeEach(func() {
		app = fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			config.Module,
			redis.Module,
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
