package matchers

import (
	. "github.com/onsi/gomega"
)

func Must(notError error) {
	Expect(notError).To(BeNil())
}

func Must2[T any](val T, notError error) T {
	Expect(notError).To(BeNil())

	return val
}
