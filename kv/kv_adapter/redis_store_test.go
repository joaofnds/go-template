package kv_adapter_test

import (
	"app/adapter/logger"
	"app/adapter/redis"
	"app/adapter/validation"
	"app/config"
	"app/kv"
	"app/kv/kv_adapter"
	"app/kv/kv_module"
	. "app/test/matchers"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestKVRedisStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KV Redis Store Suite")
}

var _ = Describe("Redis Store", func() {
	var app *fxtest.App
	var store *kv_adapter.RedisStore

	BeforeEach(func() {
		app = fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			config.Module,
			validation.Module,
			redis.Module,
			kv_module.Module,
			fx.Populate(&store),
		)
		app.RequireStart()
	})

	AfterEach(func() {
		app.RequireStop()
	})

	It("retrieves values", func(ctx SpecContext) {
		Must(store.Set(ctx, "foo", "bar"))

		val := Must2(store.Get(ctx, "foo"))

		Expect(val).To(Equal("bar"))
	})

	It("cannot get values deleted", func(ctx SpecContext) {
		Must(store.Set(ctx, "foo", "bar"))
		Must(store.Del(ctx, "foo"))

		_, err := store.Get(ctx, "foo")
		Expect(err).To(Equal(kv.ErrNotFound))
	})
})
