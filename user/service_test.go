package user_test

import (
	"web/user"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("user service", func() {
	var userService user.UserService

	BeforeEach(func() {
		userService.DeleteAll()
	})

	Describe("DeleteAll", func() {
		It("removes all users", func() {
			userService.CreateUser("joao")
			Expect(userService.List()).NotTo(BeEmpty())

			userService.DeleteAll()
			Expect(userService.List()).To(BeEmpty())
		})
	})

	It("created users can be found by name", func() {
		user := userService.CreateUser("joao")

		found, ok := userService.FindByName(user.Name)

		Expect(ok).To(BeTrue())
		Expect(found).To(Equal(user))
	})

	It("created users appear on users listing", func() {
		user := userService.CreateUser("joao")
		Expect(userService.List()).To(ContainElement(user))
	})

	It("removed users do not appear on users listing", func() {
		user := userService.CreateUser("joao")
		userService.Remove(user)
		Expect(userService.List()).NotTo(ContainElement(user))
	})
})
