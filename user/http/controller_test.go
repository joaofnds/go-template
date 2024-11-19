package http_test

import (
	"net/http"
	"testing"

	"app/test/driver"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUserHTTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "user http suite")
}

var _ = Describe("/users", Ordered, func() {
	var app *driver.Driver

	BeforeAll(func() { app = driver.Setup() })
	BeforeEach(func() { app.BeginTx() })
	AfterEach(func() { app.RollbackTx() })
	AfterAll(func() { app.Teardown() })

	It("creates and gets user", func() {
		bob := app.User.MustCreate("bob")
		found := app.User.MustGet(bob.Name)

		Expect(found).To(Equal(bob))
	})

	It("lists users", func() {
		bob := app.User.MustCreate("bob")
		dave := app.User.MustCreate("dave")

		users := app.User.MustList()
		Expect(users).To(ConsistOf(bob, dave))
	})

	It("deletes users", func() {
		bob := app.User.MustCreate("bob")
		dave := app.User.MustCreate("dave")

		app.User.MustDelete(dave.Name)

		_, err := app.User.Get(dave.Name)
		Expect(err).To(Equal(driver.RequestFailure{
			Status: http.StatusNotFound,
			Body:   "Not Found",
		}))

		users := app.User.MustList()
		Expect(users).To(ConsistOf(bob))
	})

	It("switches feature flag", func() {
		bob := app.User.MustCreate("bob")
		bobFeatures := app.User.MustGetFeature(bob.Name)
		Expect(bobFeatures).To(HaveKeyWithValue("cool-feature", "on"))

		frank := app.User.MustCreate("frank")
		frankFeatures := app.User.MustGetFeature(frank.Name)
		Expect(frankFeatures).To(HaveKeyWithValue("cool-feature", "off"))
	})
})
