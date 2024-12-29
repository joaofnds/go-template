package kv_http_test

import (
	"app/test/driver"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestKVHTTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KV HTTP Suite")
}

var _ = Describe("/kv", Ordered, func() {
	var app *driver.Driver

	BeforeAll(func() { app = driver.Setup() })
	BeforeEach(func() { app.BeginTx() })
	AfterEach(func() { app.RollbackTx() })
	AfterAll(func() { app.Teardown() })

	Context("GET", func() {
		It("returns the value under the key", func() {
			app.KV.MustSet("foo", "bar")

			val := app.KV.MustGet("foo")

			Expect(val).To(Equal("bar"))
		})
	})

	Context("POST", func() {
		It("responds with status created", func() {
			res, _ := app.KV.SetReq("foo", "bar")
			Expect(res.StatusCode).To(Equal(http.StatusCreated))
		})

		It("sets the value to the key", func() {
			app.KV.MustSet("foo", "bar")

			val := app.KV.MustGet("foo")

			Expect(val).To(Equal("bar"))
		})
	})

	Context("DELETE", func() {
		It("forgets the key", func() {
			app.KV.MustSet("bar", "foo")
			app.KV.MustDel("bar")

			res := app.KV.MustGetReq("bar")
			Expect(res.StatusCode).To(Equal(http.StatusNotFound))
		})
	})
})
