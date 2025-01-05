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

	BeforeAll(func() { app = harness.Setup() })
	BeforeEach(func() { app.BeforeEach() })
	AfterEach(func() { app.AfterEach() })
	AfterAll(func() { app.Teardown() })

	It("creates and gets user", func() {
		bob := app.NewUser("bob@template.com", "p455w0rd")
		found := bob.App.Users.MustGet(bob.Entity.ID)

		Expect(found).To(Equal(bob.Entity))
	})

	It("lists users", func() {
		bob := app.NewUser("bob@template.com", "p455w0rd")
		dave := app.NewUser("dave@template.com", "p455w0rd")

		users := app.NewDriver().Users.MustList()
		Expect(users).To(ConsistOf(bob.Entity, dave.Entity))
	})

	It("deletes users", func() {
		bob := app.NewUser("bob@template.com", "p455w0rd")
		dave := app.NewUser("dave@template.com", "p455w0rd")

		dave.App.Users.MustDelete(dave.Entity.ID)

		_, err := bob.App.Users.Get(dave.Entity.ID)
		Expect(err).To(Equal(driver.RequestFailure{
			Status: http.StatusNotFound,
			Body:   "Not Found",
		}))

		users := bob.App.Users.MustList()
		Expect(users).To(ConsistOf(bob.Entity))
	})

	It("switches feature flag", func() {
		bob := app.NewUser("bob@template.com", "p455w0rd")
		bobFeatures := bob.App.Users.MustGetFeature(bob.Entity.ID)
		Expect(bobFeatures).To(HaveKeyWithValue("cool-feature", "on"))

		frank := app.NewUser("frank@template.com", "p455w0rd")
		frankFeatures := frank.App.Users.MustGetFeature(frank.Entity.ID)
		Expect(frankFeatures).To(HaveKeyWithValue("cool-feature", "off"))
	})
})
