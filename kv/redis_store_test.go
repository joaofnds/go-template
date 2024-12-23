package kv_test

import (
	"app/adapter/logger"
	"app/adapter/redis"
	"app/adapter/validation"
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

var _ = Describe("Redis Store", func() {
	var app *fxtest.App
	var store *kv.RedisStore

	BeforeEach(func() {
		app = fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			validation.Module,
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
