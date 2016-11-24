package output_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
	"testing"
)

type ErroringWriter struct{}

func (e ErroringWriter) Write(p []byte) (int, error) {
	return 0, errors.New("error-writing-to-writer")
}

func TestPrint(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Print Suite")
}
