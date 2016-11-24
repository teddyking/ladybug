package result_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/ladybug/result"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/garden/gardenfakes"
)

var _ = Describe("Info", func() {
	Describe("Generate", func() {
		var (
			infoResult    Info
			containers    []garden.Container
			fakeContainer *gardenfakes.FakeContainer
		)

		BeforeEach(func() {
			fakeContainer = &gardenfakes.FakeContainer{}
			containers = []garden.Container{fakeContainer}
		})

		It("generates an Info", func() {
			infoResult.Generate(containers)

			Expect(infoResult.ContainersCount).To(Equal(1))
		})
	})
})
