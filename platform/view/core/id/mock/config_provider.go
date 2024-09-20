// Code generated by counterfeiter. DO NOT EDIT.
package mock

import (
	"sync"

	"github.com/hyperledger-labs/fabric-smart-client/platform/view/core/id"
)

type ConfigProvider struct {
	GetPathStub        func(string) string
	getPathMutex       sync.RWMutex
	getPathArgsForCall []struct {
		arg1 string
	}
	getPathReturns struct {
		result1 string
	}
	getPathReturnsOnCall map[int]struct {
		result1 string
	}
	GetStringSliceStub        func(string) []string
	getStringSliceMutex       sync.RWMutex
	getStringSliceArgsForCall []struct {
		arg1 string
	}
	getStringSliceReturns struct {
		result1 []string
	}
	getStringSliceReturnsOnCall map[int]struct {
		result1 []string
	}
	TranslatePathStub        func(string) string
	translatePathMutex       sync.RWMutex
	translatePathArgsForCall []struct {
		arg1 string
	}
	translatePathReturns struct {
		result1 string
	}
	translatePathReturnsOnCall map[int]struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ConfigProvider) GetPath(arg1 string) string {
	fake.getPathMutex.Lock()
	ret, specificReturn := fake.getPathReturnsOnCall[len(fake.getPathArgsForCall)]
	fake.getPathArgsForCall = append(fake.getPathArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetPathStub
	fakeReturns := fake.getPathReturns
	fake.recordInvocation("GetPath", []interface{}{arg1})
	fake.getPathMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *ConfigProvider) GetPathCallCount() int {
	fake.getPathMutex.RLock()
	defer fake.getPathMutex.RUnlock()
	return len(fake.getPathArgsForCall)
}

func (fake *ConfigProvider) GetPathCalls(stub func(string) string) {
	fake.getPathMutex.Lock()
	defer fake.getPathMutex.Unlock()
	fake.GetPathStub = stub
}

func (fake *ConfigProvider) GetPathArgsForCall(i int) string {
	fake.getPathMutex.RLock()
	defer fake.getPathMutex.RUnlock()
	argsForCall := fake.getPathArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ConfigProvider) GetPathReturns(result1 string) {
	fake.getPathMutex.Lock()
	defer fake.getPathMutex.Unlock()
	fake.GetPathStub = nil
	fake.getPathReturns = struct {
		result1 string
	}{result1}
}

func (fake *ConfigProvider) GetPathReturnsOnCall(i int, result1 string) {
	fake.getPathMutex.Lock()
	defer fake.getPathMutex.Unlock()
	fake.GetPathStub = nil
	if fake.getPathReturnsOnCall == nil {
		fake.getPathReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.getPathReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *ConfigProvider) GetStringSlice(arg1 string) []string {
	fake.getStringSliceMutex.Lock()
	ret, specificReturn := fake.getStringSliceReturnsOnCall[len(fake.getStringSliceArgsForCall)]
	fake.getStringSliceArgsForCall = append(fake.getStringSliceArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetStringSliceStub
	fakeReturns := fake.getStringSliceReturns
	fake.recordInvocation("GetStringSlice", []interface{}{arg1})
	fake.getStringSliceMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *ConfigProvider) GetStringSliceCallCount() int {
	fake.getStringSliceMutex.RLock()
	defer fake.getStringSliceMutex.RUnlock()
	return len(fake.getStringSliceArgsForCall)
}

func (fake *ConfigProvider) GetStringSliceCalls(stub func(string) []string) {
	fake.getStringSliceMutex.Lock()
	defer fake.getStringSliceMutex.Unlock()
	fake.GetStringSliceStub = stub
}

func (fake *ConfigProvider) GetStringSliceArgsForCall(i int) string {
	fake.getStringSliceMutex.RLock()
	defer fake.getStringSliceMutex.RUnlock()
	argsForCall := fake.getStringSliceArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ConfigProvider) GetStringSliceReturns(result1 []string) {
	fake.getStringSliceMutex.Lock()
	defer fake.getStringSliceMutex.Unlock()
	fake.GetStringSliceStub = nil
	fake.getStringSliceReturns = struct {
		result1 []string
	}{result1}
}

func (fake *ConfigProvider) GetStringSliceReturnsOnCall(i int, result1 []string) {
	fake.getStringSliceMutex.Lock()
	defer fake.getStringSliceMutex.Unlock()
	fake.GetStringSliceStub = nil
	if fake.getStringSliceReturnsOnCall == nil {
		fake.getStringSliceReturnsOnCall = make(map[int]struct {
			result1 []string
		})
	}
	fake.getStringSliceReturnsOnCall[i] = struct {
		result1 []string
	}{result1}
}

func (fake *ConfigProvider) TranslatePath(arg1 string) string {
	fake.translatePathMutex.Lock()
	ret, specificReturn := fake.translatePathReturnsOnCall[len(fake.translatePathArgsForCall)]
	fake.translatePathArgsForCall = append(fake.translatePathArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.TranslatePathStub
	fakeReturns := fake.translatePathReturns
	fake.recordInvocation("TranslatePath", []interface{}{arg1})
	fake.translatePathMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *ConfigProvider) TranslatePathCallCount() int {
	fake.translatePathMutex.RLock()
	defer fake.translatePathMutex.RUnlock()
	return len(fake.translatePathArgsForCall)
}

func (fake *ConfigProvider) TranslatePathCalls(stub func(string) string) {
	fake.translatePathMutex.Lock()
	defer fake.translatePathMutex.Unlock()
	fake.TranslatePathStub = stub
}

func (fake *ConfigProvider) TranslatePathArgsForCall(i int) string {
	fake.translatePathMutex.RLock()
	defer fake.translatePathMutex.RUnlock()
	argsForCall := fake.translatePathArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ConfigProvider) TranslatePathReturns(result1 string) {
	fake.translatePathMutex.Lock()
	defer fake.translatePathMutex.Unlock()
	fake.TranslatePathStub = nil
	fake.translatePathReturns = struct {
		result1 string
	}{result1}
}

func (fake *ConfigProvider) TranslatePathReturnsOnCall(i int, result1 string) {
	fake.translatePathMutex.Lock()
	defer fake.translatePathMutex.Unlock()
	fake.TranslatePathStub = nil
	if fake.translatePathReturnsOnCall == nil {
		fake.translatePathReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.translatePathReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *ConfigProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getPathMutex.RLock()
	defer fake.getPathMutex.RUnlock()
	fake.getStringSliceMutex.RLock()
	defer fake.getStringSliceMutex.RUnlock()
	fake.translatePathMutex.RLock()
	defer fake.translatePathMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ConfigProvider) recordInvocation(key string, args []interface{}) {
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

var _ id.ConfigProvider = new(ConfigProvider)
