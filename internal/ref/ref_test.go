package ref_test

import (
	"app/internal/ref"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRef(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ref suite")
}

var _ = Describe("Ref", func() {
	Describe("New", func() {
		It("creates a new ref", func() {
			userRef := ref.New("user", "111")

			Expect(userRef.Type).To(Equal("user"))
			Expect(userRef.ID).To(Equal("111"))
		})
	})

	Describe("NewFromString", func() {
		It("parses a ref string", func() {
			userRef := ref.NewFromString("user:111")

			Expect(userRef.Type).To(Equal("user"))
			Expect(userRef.ID).To(Equal("111"))
		})

		When("':' is missing", func() {
			It("panics", func() {
				Expect(func() {
					ref.NewFromString("user111")
				}).To(PanicWith("ref string must have len 2 after split"))
			})
		})

		When("type is missing", func() {
			It("panics", func() {
				Expect(func() {
					ref.NewFromString(":111")
				}).To(PanicWith("ref string must have non-empty parts"))
			})
		})

		When("id is missing", func() {
			It("panics", func() {
				Expect(func() {
					ref.NewFromString("user:")
				}).To(PanicWith("ref string must have non-empty parts"))
			})
		})
	})
})
