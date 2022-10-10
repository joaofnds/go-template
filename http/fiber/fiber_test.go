package fiber_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFiber(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "fiber suite")
}
