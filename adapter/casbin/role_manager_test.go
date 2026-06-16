package casbin_test

import (
	"app/adapter/casbin"
	"app/internal/ref"
	"app/test/harness"
	. "app/test/matchers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("casbin role manager", Ordered, func() {
	user := ref.New("user", "111")
	admin := ref.New("role", "admin")
	customer := ref.New("role", "customer")

	var (
		app *harness.Harness
		sut *casbin.RoleManager
	)

	BeforeAll(func() { app = harness.Setup(harness.Populate(&sut)) })
	BeforeEach(func() { app.BeforeEach() })
	AfterEach(func() { app.AfterEach() })
	AfterAll(func() { app.Teardown() })

	It("grants roles", func() {
		Must(sut.Assign(user, admin))
	})

	It("retrieves roles", func() {
		Must(sut.Assign(user, admin))
		Must(sut.Assign(user, customer))

		Expect(sut.GetAll(user)).To(ConsistOf(admin, customer))
	})

	It("revokes roles", func() {
		Must(sut.Assign(user, admin))

		Must(sut.Revoke(user, admin))

		Expect(sut.GetAll(user)).To(BeEmpty())
	})

	It("removes all roles for a user", func() {
		Must(sut.Assign(user, admin))
		Must(sut.Assign(user, customer))

		Expect(sut.GetAll(user)).To(ConsistOf(admin, customer))

		Must(sut.RevokeAll(user))

		Expect(sut.GetAll(user)).To(BeEmpty())
	})
})
