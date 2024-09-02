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
			userRef := ref.New("111", "user")

			Expect(userRef.ID).To(Equal("111"))
			Expect(userRef.Type).To(Equal("user"))
		})
	})

	Describe("NewFromString", func() {
		It("parses a ref string", func() {
			userRef := ref.NewFromString("user:111")

			Expect(userRef.ID).To(Equal("111"))
			Expect(userRef.Type).To(Equal("user"))
		})
	})
})
