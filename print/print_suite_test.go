package print_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPrint(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Print Suite")
}
