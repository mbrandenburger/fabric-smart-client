// Code generated by counterfeiter. DO NOT EDIT.
package mock

import (
	"sync"

	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/kvs"
)

type ConfigProvider struct {
	GetIntStub        func(string) int
	getIntMutex       sync.RWMutex
	getIntArgsForCall []struct {
		arg1 string
	}
	getIntReturns struct {
		result1 int
	}
	getIntReturnsOnCall map[int]struct {
		result1 int
	}
	IsSetStub        func(string) bool
	isSetMutex       sync.RWMutex
	isSetArgsForCall []struct {
		arg1 string
	}
	isSetReturns struct {
		result1 bool
	}
	isSetReturnsOnCall map[int]struct {
		result1 bool
	}
	UnmarshalKeyStub        func(string, interface{}) error
	unmarshalKeyMutex       sync.RWMutex
	unmarshalKeyArgsForCall []struct {
		arg1 string
		arg2 interface{}
	}
	unmarshalKeyReturns struct {
		result1 error
	}
	unmarshalKeyReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ConfigProvider) GetInt(arg1 string) int {
	fake.getIntMutex.Lock()
	ret, specificReturn := fake.getIntReturnsOnCall[len(fake.getIntArgsForCall)]
	fake.getIntArgsForCall = append(fake.getIntArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetIntStub
	fakeReturns := fake.getIntReturns
	fake.recordInvocation("GetInt", []interface{}{arg1})
	fake.getIntMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *ConfigProvider) GetIntCallCount() int {
	fake.getIntMutex.RLock()
	defer fake.getIntMutex.RUnlock()
	return len(fake.getIntArgsForCall)
}

func (fake *ConfigProvider) GetIntCalls(stub func(string) int) {
	fake.getIntMutex.Lock()
	defer fake.getIntMutex.Unlock()
	fake.GetIntStub = stub
}

func (fake *ConfigProvider) GetIntArgsForCall(i int) string {
	fake.getIntMutex.RLock()
	defer fake.getIntMutex.RUnlock()
	argsForCall := fake.getIntArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ConfigProvider) GetIntReturns(result1 int) {
	fake.getIntMutex.Lock()
	defer fake.getIntMutex.Unlock()
	fake.GetIntStub = nil
	fake.getIntReturns = struct {
		result1 int
	}{result1}
}

func (fake *ConfigProvider) GetIntReturnsOnCall(i int, result1 int) {
	fake.getIntMutex.Lock()
	defer fake.getIntMutex.Unlock()
	fake.GetIntStub = nil
	if fake.getIntReturnsOnCall == nil {
		fake.getIntReturnsOnCall = make(map[int]struct {
			result1 int
		})
	}
	fake.getIntReturnsOnCall[i] = struct {
		result1 int
	}{result1}
}

func (fake *ConfigProvider) IsSet(arg1 string) bool {
	fake.isSetMutex.Lock()
	ret, specificReturn := fake.isSetReturnsOnCall[len(fake.isSetArgsForCall)]
	fake.isSetArgsForCall = append(fake.isSetArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.IsSetStub
	fakeReturns := fake.isSetReturns
	fake.recordInvocation("IsSet", []interface{}{arg1})
	fake.isSetMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *ConfigProvider) IsSetCallCount() int {
	fake.isSetMutex.RLock()
	defer fake.isSetMutex.RUnlock()
	return len(fake.isSetArgsForCall)
}

func (fake *ConfigProvider) IsSetCalls(stub func(string) bool) {
	fake.isSetMutex.Lock()
	defer fake.isSetMutex.Unlock()
	fake.IsSetStub = stub
}

func (fake *ConfigProvider) IsSetArgsForCall(i int) string {
	fake.isSetMutex.RLock()
	defer fake.isSetMutex.RUnlock()
	argsForCall := fake.isSetArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ConfigProvider) IsSetReturns(result1 bool) {
	fake.isSetMutex.Lock()
	defer fake.isSetMutex.Unlock()
	fake.IsSetStub = nil
	fake.isSetReturns = struct {
		result1 bool
	}{result1}
}

func (fake *ConfigProvider) IsSetReturnsOnCall(i int, result1 bool) {
	fake.isSetMutex.Lock()
	defer fake.isSetMutex.Unlock()
	fake.IsSetStub = nil
	if fake.isSetReturnsOnCall == nil {
		fake.isSetReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isSetReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *ConfigProvider) UnmarshalKey(arg1 string, arg2 interface{}) error {
	fake.unmarshalKeyMutex.Lock()
	ret, specificReturn := fake.unmarshalKeyReturnsOnCall[len(fake.unmarshalKeyArgsForCall)]
	fake.unmarshalKeyArgsForCall = append(fake.unmarshalKeyArgsForCall, struct {
		arg1 string
		arg2 interface{}
	}{arg1, arg2})
	stub := fake.UnmarshalKeyStub
	fakeReturns := fake.unmarshalKeyReturns
	fake.recordInvocation("UnmarshalKey", []interface{}{arg1, arg2})
	fake.unmarshalKeyMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *ConfigProvider) UnmarshalKeyCallCount() int {
	fake.unmarshalKeyMutex.RLock()
	defer fake.unmarshalKeyMutex.RUnlock()
	return len(fake.unmarshalKeyArgsForCall)
}

func (fake *ConfigProvider) UnmarshalKeyCalls(stub func(string, interface{}) error) {
	fake.unmarshalKeyMutex.Lock()
	defer fake.unmarshalKeyMutex.Unlock()
	fake.UnmarshalKeyStub = stub
}

func (fake *ConfigProvider) UnmarshalKeyArgsForCall(i int) (string, interface{}) {
	fake.unmarshalKeyMutex.RLock()
	defer fake.unmarshalKeyMutex.RUnlock()
	argsForCall := fake.unmarshalKeyArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *ConfigProvider) UnmarshalKeyReturns(result1 error) {
	fake.unmarshalKeyMutex.Lock()
	defer fake.unmarshalKeyMutex.Unlock()
	fake.UnmarshalKeyStub = nil
	fake.unmarshalKeyReturns = struct {
		result1 error
	}{result1}
}

func (fake *ConfigProvider) UnmarshalKeyReturnsOnCall(i int, result1 error) {
	fake.unmarshalKeyMutex.Lock()
	defer fake.unmarshalKeyMutex.Unlock()
	fake.UnmarshalKeyStub = nil
	if fake.unmarshalKeyReturnsOnCall == nil {
		fake.unmarshalKeyReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.unmarshalKeyReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ConfigProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getIntMutex.RLock()
	defer fake.getIntMutex.RUnlock()
	fake.isSetMutex.RLock()
	defer fake.isSetMutex.RUnlock()
	fake.unmarshalKeyMutex.RLock()
	defer fake.unmarshalKeyMutex.RUnlock()
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

var _ kvs.ConfigProvider = new(ConfigProvider)
