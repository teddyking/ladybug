package exterminator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
	"testing"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/garden/client"
	"code.cloudfoundry.org/garden/client/connection"
	"github.com/onsi/gomega/gexec"
)

var (
	pathToLadybug string
	gardenClient  garden.Client
)

func TestExterminator(t *testing.T) {
	BeforeSuite(func() {
		var err error

		pathToLadybug, err = gexec.Build("github.com/teddyking/ladybug")
		Expect(err).NotTo(HaveOccurred())

		gardenAddress := os.Getenv("GARDEN_ADDRESS")
		if gardenAddress == "" {
			gardenAddress = "127.0.0.1:7777"
		}

		gardenClient = client.New(connection.New("tcp", gardenAddress))
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Exterminator Suite")
}
