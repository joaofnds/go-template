package kv_http_test

import (
	"app/test/driver"
	"app/test/harness"
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
	var app *harness.Harness
	var api *driver.Driver

	BeforeAll(func() { app = harness.Setup(); api = app.NewDriver() })
	BeforeEach(func() { app.BeforeEach() })
	AfterEach(func() { app.AfterEach() })
	AfterAll(func() { app.Teardown() })

	Context("GET", func() {
		It("returns the value under the key", func() {
			api.KV.MustSet("foo", "bar")

			val := api.KV.MustGet("foo")

			Expect(val).To(Equal("bar"))
		})
	})

	Context("POST", func() {
		It("responds with status created", func() {
			resp, _ := api.KV.SetReq("foo", "bar")
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))
		})

		It("sets the value to the key", func() {
			api.KV.MustSet("foo", "bar")

			val := api.KV.MustGet("foo")

			Expect(val).To(Equal("bar"))
		})
	})

	Context("DELETE", func() {
		It("forgets the key", func() {
			api.KV.MustSet("bar", "foo")
			api.KV.MustDel("bar")

			resp := api.KV.MustGetReq("bar")
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})
	})
})
