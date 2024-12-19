package user_http_test

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
		bob := app.User.MustCreate("bob@template.com")
		found := app.User.MustGet(bob.ID)

		Expect(found).To(Equal(bob))
	})

	It("lists users", func() {
		bob := app.User.MustCreate("bob@template.com")
		dave := app.User.MustCreate("dave@template.com")

		users := app.User.MustList()
		Expect(users).To(ConsistOf(bob, dave))
	})

	It("deletes users", func() {
		bob := app.User.MustCreate("bob@template.com")
		dave := app.User.MustCreate("dave@template.com")

		app.User.MustDelete(dave.ID)

		_, err := app.User.Get(dave.ID)
		Expect(err).To(Equal(driver.RequestFailure{
			Status: http.StatusNotFound,
			Body:   "Not Found",
		}))

		users := app.User.MustList()
		Expect(users).To(ConsistOf(bob))
	})

	It("switches feature flag", func() {
		bob := app.User.MustCreate("bob@template.com")
		bobFeatures := app.User.MustGetFeature(bob.ID)
		Expect(bobFeatures).To(HaveKeyWithValue("cool-feature", "on"))

		frank := app.User.MustCreate("frank@template.com")
		frankFeatures := app.User.MustGetFeature(frank.ID)
		Expect(frankFeatures).To(HaveKeyWithValue("cool-feature", "off"))
	})
})
