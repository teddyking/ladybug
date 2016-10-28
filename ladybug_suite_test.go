package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/onsi/gomega/gexec"
)

var pathToLadybug string

func TestLadybug(t *testing.T) {
	BeforeSuite(func() {
		var err error
		pathToLadybug, err = gexec.Build("github.com/teddyking/ladybug")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Ladybug Suite")
}
