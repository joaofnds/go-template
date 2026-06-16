package user_test

import (
	"testing"

	"app/test/harness"
	. "app/test/matchers"
	"app/user"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUserService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "user service suite")
}

var _ = Describe("user service", Ordered, func() {
	var (
		app         *harness.Harness
		userService *user.Service
	)

	BeforeAll(func() { app = harness.Setup(harness.Populate(&userService)) })
	BeforeEach(func() { app.BeforeEach() })
	AfterEach(func() { app.AfterEach() })
	AfterAll(func() { app.Teardown() })

	Describe("DeleteAll", func() {
		It("removes all users", func(ctx SpecContext) {
			Must2(userService.CreateUser(ctx, "joao@template.com"))
			Expect(userService.List(ctx)).NotTo(BeEmpty())

			Must(userService.DeleteAll(ctx))
			Expect(userService.List(ctx)).To(BeEmpty())
		})
	})

	It("created users can be found by id", func(ctx SpecContext) {
		user := Must2(userService.CreateUser(ctx, "joao@template.com"))

		found := Must2(userService.FindByID(ctx, user.ID))
		Expect(found).To(Equal(user))
	})

	It("created users can be found by email", func(ctx SpecContext) {
		user := Must2(userService.CreateUser(ctx, "joao@template.com"))

		found := Must2(userService.FindByEmail(ctx, user.Email))
		Expect(found).To(Equal(user))
	})

	It("lists created users", func(ctx SpecContext) {
		user := Must2(userService.CreateUser(ctx, "joao@template.com"))
		Expect(userService.List(ctx)).To(ContainElement(user))
	})

	It("removed users are not listed", func(ctx SpecContext) {
		user := Must2(userService.CreateUser(ctx, "joao@template.com"))
		Must(userService.Remove(ctx, user))

		users := Must2(userService.List(ctx))
		Expect(users).NotTo(ContainElement(user))
	})
})
