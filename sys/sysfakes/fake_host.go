// This file was generated by counterfeiter
package sysfakes

import (
	"sync"

	"github.com/teddyking/ladybug/sys"
)

type FakeHost struct {
	ContainerPidsStub        func(handle string) ([]string, error)
	containerPidsMutex       sync.RWMutex
	containerPidsArgsForCall []struct {
		handle string
	}
	containerPidsReturns struct {
		result1 []string
		result2 error
	}
	ContainerProcessNameStub        func(pid string) (string, error)
	containerProcessNameMutex       sync.RWMutex
	containerProcessNameArgsForCall []struct {
		pid string
	}
	containerProcessNameReturns struct {
		result1 string
		result2 error
	}
	ContainerCreationTimeStub        func(handle string) (string, error)
	containerCreationTimeMutex       sync.RWMutex
	containerCreationTimeArgsForCall []struct {
		handle string
	}
	containerCreationTimeReturns struct {
		result1 string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeHost) ContainerPids(handle string) ([]string, error) {
	fake.containerPidsMutex.Lock()
	fake.containerPidsArgsForCall = append(fake.containerPidsArgsForCall, struct {
		handle string
	}{handle})
	fake.recordInvocation("ContainerPids", []interface{}{handle})
	fake.containerPidsMutex.Unlock()
	if fake.ContainerPidsStub != nil {
		return fake.ContainerPidsStub(handle)
	} else {
		return fake.containerPidsReturns.result1, fake.containerPidsReturns.result2
	}
}

func (fake *FakeHost) ContainerPidsCallCount() int {
	fake.containerPidsMutex.RLock()
	defer fake.containerPidsMutex.RUnlock()
	return len(fake.containerPidsArgsForCall)
}

func (fake *FakeHost) ContainerPidsArgsForCall(i int) string {
	fake.containerPidsMutex.RLock()
	defer fake.containerPidsMutex.RUnlock()
	return fake.containerPidsArgsForCall[i].handle
}

func (fake *FakeHost) ContainerPidsReturns(result1 []string, result2 error) {
	fake.ContainerPidsStub = nil
	fake.containerPidsReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeHost) ContainerProcessName(pid string) (string, error) {
	fake.containerProcessNameMutex.Lock()
	fake.containerProcessNameArgsForCall = append(fake.containerProcessNameArgsForCall, struct {
		pid string
	}{pid})
	fake.recordInvocation("ContainerProcessName", []interface{}{pid})
	fake.containerProcessNameMutex.Unlock()
	if fake.ContainerProcessNameStub != nil {
		return fake.ContainerProcessNameStub(pid)
	} else {
		return fake.containerProcessNameReturns.result1, fake.containerProcessNameReturns.result2
	}
}

func (fake *FakeHost) ContainerProcessNameCallCount() int {
	fake.containerProcessNameMutex.RLock()
	defer fake.containerProcessNameMutex.RUnlock()
	return len(fake.containerProcessNameArgsForCall)
}

func (fake *FakeHost) ContainerProcessNameArgsForCall(i int) string {
	fake.containerProcessNameMutex.RLock()
	defer fake.containerProcessNameMutex.RUnlock()
	return fake.containerProcessNameArgsForCall[i].pid
}

func (fake *FakeHost) ContainerProcessNameReturns(result1 string, result2 error) {
	fake.ContainerProcessNameStub = nil
	fake.containerProcessNameReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeHost) ContainerCreationTime(handle string) (string, error) {
	fake.containerCreationTimeMutex.Lock()
	fake.containerCreationTimeArgsForCall = append(fake.containerCreationTimeArgsForCall, struct {
		handle string
	}{handle})
	fake.recordInvocation("ContainerCreationTime", []interface{}{handle})
	fake.containerCreationTimeMutex.Unlock()
	if fake.ContainerCreationTimeStub != nil {
		return fake.ContainerCreationTimeStub(handle)
	} else {
		return fake.containerCreationTimeReturns.result1, fake.containerCreationTimeReturns.result2
	}
}

func (fake *FakeHost) ContainerCreationTimeCallCount() int {
	fake.containerCreationTimeMutex.RLock()
	defer fake.containerCreationTimeMutex.RUnlock()
	return len(fake.containerCreationTimeArgsForCall)
}

func (fake *FakeHost) ContainerCreationTimeArgsForCall(i int) string {
	fake.containerCreationTimeMutex.RLock()
	defer fake.containerCreationTimeMutex.RUnlock()
	return fake.containerCreationTimeArgsForCall[i].handle
}

func (fake *FakeHost) ContainerCreationTimeReturns(result1 string, result2 error) {
	fake.ContainerCreationTimeStub = nil
	fake.containerCreationTimeReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeHost) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.containerPidsMutex.RLock()
	defer fake.containerPidsMutex.RUnlock()
	fake.containerProcessNameMutex.RLock()
	defer fake.containerProcessNameMutex.RUnlock()
	fake.containerCreationTimeMutex.RLock()
	defer fake.containerCreationTimeMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeHost) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ sys.Host = new(FakeHost)