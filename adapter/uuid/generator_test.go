package uuid_test

import (
	"app/adapter/uuid"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUUID(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "uuid suite")
}

var _ = Describe("UUIDGenerator", func() {
	var generator *uuid.Generator

	BeforeEach(func() { generator = uuid.NewGenerator() })

	Describe("NewID", func() {
		It("returns a new UUID", func() {

			Expect(generator.NewID()).To(HaveLen(36))
		})
	})
})
