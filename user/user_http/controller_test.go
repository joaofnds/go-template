package user_http_test

import (
	"net/http"
	"testing"

	"app/test/driver"
	"app/test/harness"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUserHTTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "user http suite")
}

var _ = Describe("/users", Ordered, func() {
	var app *harness.Harness
	var api *driver.Driver

	BeforeAll(func() { app = harness.Setup(); api = app.NewDriver() })
	BeforeEach(func() { app.BeforeEach() })
	AfterEach(func() { app.AfterEach() })
	AfterAll(func() { app.Teardown() })

	It("creates and gets user", func() {
		bob := api.Users.MustCreate("bob@template.com")
		found := api.Users.MustGet(bob.ID)

		Expect(found).To(Equal(bob))
	})

	It("lists users", func() {
		bob := api.Users.MustCreate("bob@template.com")
		dave := api.Users.MustCreate("dave@template.com")

		users := api.Users.MustList()
		Expect(users).To(ConsistOf(bob, dave))
	})

	It("deletes users", func() {
		bob := api.Users.MustCreate("bob@template.com")
		dave := api.Users.MustCreate("dave@template.com")

		api.Users.MustDelete(dave.ID)

		_, err := api.Users.Get(dave.ID)
		Expect(err).To(Equal(driver.RequestFailure{
			Status: http.StatusNotFound,
			Body:   "Not Found",
		}))

		users := api.Users.MustList()
		Expect(users).To(ConsistOf(bob))
	})

	It("switches feature flag", func() {
		bob := api.Users.MustCreate("bob@template.com")
		bobFeatures := api.Users.MustGetFeature(bob.ID)
		Expect(bobFeatures).To(HaveKeyWithValue("cool-feature", "on"))

		frank := api.Users.MustCreate("frank@template.com")
		frankFeatures := api.Users.MustGetFeature(frank.ID)
		Expect(frankFeatures).To(HaveKeyWithValue("cool-feature", "off"))
	})
})
