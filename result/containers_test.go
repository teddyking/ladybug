package result_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/ladybug/result"

	"errors"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/garden/gardenfakes"
	"github.com/teddyking/ladybug/sys/sysfakes"
)

var _ = Describe("Containers", func() {
	var (
		fakeContainer    *gardenfakes.FakeContainer
		fakeContainer2   *gardenfakes.FakeContainer
		fakeInfo         garden.ContainerInfo
		fakeInfo2        garden.ContainerInfo
		containers       []garden.Container
		containersResult Containers
	)

	BeforeEach(func() {
		fakeContainer = &gardenfakes.FakeContainer{}
		fakeContainer2 = &gardenfakes.FakeContainer{}
		fakeInfo = garden.ContainerInfo{
			ContainerIP: "192.0.2.10",
			MappedPorts: []garden.PortMapping{{0, 0}},
		}
		fakeInfo2 = garden.ContainerInfo{
			ContainerIP: "192.0.2.11",
			MappedPorts: []garden.PortMapping{{0, 1}},
		}

		fakeContainer.HandleReturns("fake-container")
		fakeContainer2.HandleReturns("fake-container-2")

		fakeContainer.InfoReturns(fakeInfo, nil)
		fakeContainer2.InfoReturns(fakeInfo2, nil)

		containers = []garden.Container{
			fakeContainer,
			fakeContainer2,
		}
		containersResult = make(Containers, len(containers))
	})

	Describe("WithHandles", func() {
		It("returns a func that adds handles to the Containers", func() {
			withHandlesFunc := WithHandles(containers)

			updatedResult, err := withHandlesFunc(containersResult)
			Expect(err).NotTo(HaveOccurred())

			_, ok := updatedResult["fake-container"]
			Expect(ok).To(BeTrue())
			_, ok = updatedResult["fake-container-2"]
			Expect(ok).To(BeTrue())
		})
	})

	Describe("WithIPs", func() {
		It("returns a func that adds IPs to the Containers", func() {
			withIPsFunc := WithIPs(containers)

			updatedResult, err := withIPsFunc(containersResult)
			Expect(err).NotTo(HaveOccurred())

			Expect(updatedResult["fake-container"].Ip).To(Equal("192.0.2.10"))
			Expect(updatedResult["fake-container-2"].Ip).To(Equal("192.0.2.11"))
		})

		Context("when ContainerInfo returns an error", func() {
			BeforeEach(func() {
				fakeContainer.InfoReturns(garden.ContainerInfo{}, errors.New("error-fetching-container-info"))
			})

			It("returns the error", func() {
				withIPsFunc := WithIPs(containers)

				_, err := withIPsFunc(containersResult)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("error-fetching-container-info"))
			})
		})
	})

	Describe("WithProcessNames", func() {
		var fakeHost *sysfakes.FakeHost

		BeforeEach(func() {
			fakeHost = &sysfakes.FakeHost{}

			fakeHost.ContainerPidsReturns([]string{"container-pid"}, nil)
			fakeHost.ContainerProcessNameReturns("ruby", nil)
		})

		It("returns a func that adds process names to the Containers", func() {
			withProcessNamesFunc := WithProcessNames(containers, fakeHost)

			updatedResult, err := withProcessNamesFunc(containersResult)
			Expect(err).NotTo(HaveOccurred())

			Expect(updatedResult["fake-container"].ProcessName).To(Equal("ruby"))
			Expect(updatedResult["fake-container-2"].ProcessName).To(Equal("ruby"))
		})

		Context("when ContainerPids returns an error", func() {
			BeforeEach(func() {
				fakeHost.ContainerPidsReturns(nil, errors.New("error-fetching-pids"))
			})

			It("returns the error", func() {
				withProcessNamesFunc := WithProcessNames(containers, fakeHost)

				_, err := withProcessNamesFunc(containersResult)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("error-fetching-pids"))
			})
		})

		Context("when ContainerProcessName returns an error", func() {
			BeforeEach(func() {
				fakeHost.ContainerProcessNameReturns("", errors.New("error-fetching-process-name"))
			})

			It("returns the error", func() {
				withProcessNamesFunc := WithProcessNames(containers, fakeHost)

				_, err := withProcessNamesFunc(containersResult)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("error-fetching-process-name"))
			})
		})
	})

	Describe("WithCreatedAtTimes", func() {
		var fakeHost *sysfakes.FakeHost

		BeforeEach(func() {
			fakeHost = &sysfakes.FakeHost{}

			fakeHost.ContainerCreationTimeReturns("27-07-1989", nil)
		})

		It("returns a func that adds created at times to the Containers", func() {
			withCreatedAtTimesFunc := WithCreatedAtTimes(containers, fakeHost)

			updatedResult, err := withCreatedAtTimesFunc(containersResult)
			Expect(err).NotTo(HaveOccurred())

			Expect(updatedResult["fake-container"].CreatedAt).To(Equal("27-07-1989"))
			Expect(updatedResult["fake-container-2"].CreatedAt).To(Equal("27-07-1989"))
		})

		Context("when ContainerCreationTime returns an error", func() {
			BeforeEach(func() {
				fakeHost.ContainerCreationTimeReturns("", errors.New("error-fetching-created-at-time"))
			})

			It("returns the error", func() {
				withCreatedAtTimesFunc := WithCreatedAtTimes(containers, fakeHost)

				_, err := withCreatedAtTimesFunc(containersResult)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("error-fetching-created-at-time"))
			})
		})
	})

	Describe("WithPortMappings", func() {
		It("returns a func that adds port mappings to the Containers", func() {
			withPortMappingsFunc := WithPortMappings(containers)

			updatedResult, err := withPortMappingsFunc(containersResult)
			Expect(err).NotTo(HaveOccurred())

			Expect(updatedResult["fake-container"].PortMappings).To(Equal([]garden.PortMapping{{0, 0}}))
			Expect(updatedResult["fake-container-2"].PortMappings).To(Equal([]garden.PortMapping{{0, 1}}))
		})

		Context("when ContainerInfo returns an error", func() {
			BeforeEach(func() {
				fakeContainer.InfoReturns(garden.ContainerInfo{}, errors.New("error-fetching-container-info"))
			})

			It("returns the error", func() {
				withPortMappingsFunc := WithPortMappings(containers)

				_, err := withPortMappingsFunc(containersResult)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("error-fetching-container-info"))
			})
		})
	})

	Describe("Generate", func() {
		It("loops over the provided resultModifier funcs, updating the Containers for each", func() {
			containersResult.Generate(
				WithHandles(containers),
				WithIPs(containers),
			)

			Expect(len(containersResult)).To(Equal(2))
			Expect(containersResult["fake-container"].Ip).To(Equal("192.0.2.10"))
			Expect(containersResult["fake-container-2"].Ip).To(Equal("192.0.2.11"))
		})

		Context("when one of the resultModifier funcs returns an error", func() {
			It("returns the error", func() {
				erroringResultModifierFunc := func(Containers) (Containers, error) {
					return nil, errors.New("resultModifiers-func-error")
				}

				err := containersResult.Generate(erroringResultModifierFunc)
				Expect(err).NotTo(Succeed())
				Expect(err.Error()).To(Equal("resultModifiers-func-error"))
			})
		})
	})
})
