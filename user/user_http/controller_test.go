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
	BeforeEach(func() { app.BeforeEach() })
	AfterEach(func() { app.AfterEach() })
	AfterAll(func() { app.Teardown() })

	It("creates and gets user", func() {
		bob := app.Users.MustCreate("bob@template.com")
		found := app.Users.MustGet(bob.ID)

		Expect(found).To(Equal(bob))
	})

	It("lists users", func() {
		bob := app.Users.MustCreate("bob@template.com")
		dave := app.Users.MustCreate("dave@template.com")

		users := app.Users.MustList()
		Expect(users).To(ConsistOf(bob, dave))
	})

	It("deletes users", func() {
		bob := app.Users.MustCreate("bob@template.com")
		dave := app.Users.MustCreate("dave@template.com")

		app.Users.MustDelete(dave.ID)

		_, err := app.Users.Get(dave.ID)
		Expect(err).To(Equal(driver.RequestFailure{
			Status: http.StatusNotFound,
			Body:   "Not Found",
		}))

		users := app.Users.MustList()
		Expect(users).To(ConsistOf(bob))
	})

	It("switches feature flag", func() {
		bob := app.Users.MustCreate("bob@template.com")
		bobFeatures := app.Users.MustGetFeature(bob.ID)
		Expect(bobFeatures).To(HaveKeyWithValue("cool-feature", "on"))

		frank := app.Users.MustCreate("frank@template.com")
		frankFeatures := app.Users.MustGetFeature(frank.ID)
		Expect(frankFeatures).To(HaveKeyWithValue("cool-feature", "off"))
	})
})
