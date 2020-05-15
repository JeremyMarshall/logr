package logrusr_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogrusr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Logrusr Suite")
}
