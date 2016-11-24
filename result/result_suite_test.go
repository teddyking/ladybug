package result_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestResult(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Result Suite")
}
