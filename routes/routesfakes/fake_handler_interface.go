// Code generated by counterfeiter. DO NOT EDIT.
package routesfakes

import (
	"example-project/routes"
	"sync"

	"github.com/gin-gonic/gin"
)

type FakeHandlerInterface struct {
	CreatTimeEntryStub        func(*gin.Context)
	creatTimeEntryMutex       sync.RWMutex
	creatTimeEntryArgsForCall []struct {
		arg1 *gin.Context
	}
	CreateEmployeeHandlerStub        func(*gin.Context)
	createEmployeeHandlerMutex       sync.RWMutex
	createEmployeeHandlerArgsForCall []struct {
		arg1 *gin.Context
	}
	DeleteTimeEntryStub        func(*gin.Context)
	deleteTimeEntryMutex       sync.RWMutex
	deleteTimeEntryArgsForCall []struct {
		arg1 *gin.Context
	}
	GetAllTimeEntryStub        func(*gin.Context)
	getAllTimeEntryMutex       sync.RWMutex
	getAllTimeEntryArgsForCall []struct {
		arg1 *gin.Context
	}
	GetEmployeeHandlerStub        func(*gin.Context)
	getEmployeeHandlerMutex       sync.RWMutex
	getEmployeeHandlerArgsForCall []struct {
		arg1 *gin.Context
	}
	GetTimeEntryStub        func(*gin.Context)
	getTimeEntryMutex       sync.RWMutex
	getTimeEntryArgsForCall []struct {
		arg1 *gin.Context
	}
	UpdateTimeEntryStub        func(*gin.Context)
	updateTimeEntryMutex       sync.RWMutex
	updateTimeEntryArgsForCall []struct {
		arg1 *gin.Context
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeHandlerInterface) CreatTimeEntry(arg1 *gin.Context) {
	fake.creatTimeEntryMutex.Lock()
	fake.creatTimeEntryArgsForCall = append(fake.creatTimeEntryArgsForCall, struct {
		arg1 *gin.Context
	}{arg1})
	stub := fake.CreatTimeEntryStub
	fake.recordInvocation("CreatTimeEntry", []interface{}{arg1})
	fake.creatTimeEntryMutex.Unlock()
	if stub != nil {
		fake.CreatTimeEntryStub(arg1)
	}
}

func (fake *FakeHandlerInterface) CreatTimeEntryCallCount() int {
	fake.creatTimeEntryMutex.RLock()
	defer fake.creatTimeEntryMutex.RUnlock()
	return len(fake.creatTimeEntryArgsForCall)
}

func (fake *FakeHandlerInterface) CreatTimeEntryCalls(stub func(*gin.Context)) {
	fake.creatTimeEntryMutex.Lock()
	defer fake.creatTimeEntryMutex.Unlock()
	fake.CreatTimeEntryStub = stub
}

func (fake *FakeHandlerInterface) CreatTimeEntryArgsForCall(i int) *gin.Context {
	fake.creatTimeEntryMutex.RLock()
	defer fake.creatTimeEntryMutex.RUnlock()
	argsForCall := fake.creatTimeEntryArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeHandlerInterface) CreateEmployeeHandler(arg1 *gin.Context) {
	fake.createEmployeeHandlerMutex.Lock()
	fake.createEmployeeHandlerArgsForCall = append(fake.createEmployeeHandlerArgsForCall, struct {
		arg1 *gin.Context
	}{arg1})
	stub := fake.CreateEmployeeHandlerStub
	fake.recordInvocation("CreateEmployeeHandler", []interface{}{arg1})
	fake.createEmployeeHandlerMutex.Unlock()
	if stub != nil {
		fake.CreateEmployeeHandlerStub(arg1)
	}
}

func (fake *FakeHandlerInterface) CreateEmployeeHandlerCallCount() int {
	fake.createEmployeeHandlerMutex.RLock()
	defer fake.createEmployeeHandlerMutex.RUnlock()
	return len(fake.createEmployeeHandlerArgsForCall)
}

func (fake *FakeHandlerInterface) CreateEmployeeHandlerCalls(stub func(*gin.Context)) {
	fake.createEmployeeHandlerMutex.Lock()
	defer fake.createEmployeeHandlerMutex.Unlock()
	fake.CreateEmployeeHandlerStub = stub
}

func (fake *FakeHandlerInterface) CreateEmployeeHandlerArgsForCall(i int) *gin.Context {
	fake.createEmployeeHandlerMutex.RLock()
	defer fake.createEmployeeHandlerMutex.RUnlock()
	argsForCall := fake.createEmployeeHandlerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeHandlerInterface) DeleteTimeEntry(arg1 *gin.Context) {
	fake.deleteTimeEntryMutex.Lock()
	fake.deleteTimeEntryArgsForCall = append(fake.deleteTimeEntryArgsForCall, struct {
		arg1 *gin.Context
	}{arg1})
	stub := fake.DeleteTimeEntryStub
	fake.recordInvocation("DeleteTimeEntry", []interface{}{arg1})
	fake.deleteTimeEntryMutex.Unlock()
	if stub != nil {
		fake.DeleteTimeEntryStub(arg1)
	}
}

func (fake *FakeHandlerInterface) DeleteTimeEntryCallCount() int {
	fake.deleteTimeEntryMutex.RLock()
	defer fake.deleteTimeEntryMutex.RUnlock()
	return len(fake.deleteTimeEntryArgsForCall)
}

func (fake *FakeHandlerInterface) DeleteTimeEntryCalls(stub func(*gin.Context)) {
	fake.deleteTimeEntryMutex.Lock()
	defer fake.deleteTimeEntryMutex.Unlock()
	fake.DeleteTimeEntryStub = stub
}

func (fake *FakeHandlerInterface) DeleteTimeEntryArgsForCall(i int) *gin.Context {
	fake.deleteTimeEntryMutex.RLock()
	defer fake.deleteTimeEntryMutex.RUnlock()
	argsForCall := fake.deleteTimeEntryArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeHandlerInterface) GetAllTimeEntry(arg1 *gin.Context) {
	fake.getAllTimeEntryMutex.Lock()
	fake.getAllTimeEntryArgsForCall = append(fake.getAllTimeEntryArgsForCall, struct {
		arg1 *gin.Context
	}{arg1})
	stub := fake.GetAllTimeEntryStub
	fake.recordInvocation("GetAllTimeEntry", []interface{}{arg1})
	fake.getAllTimeEntryMutex.Unlock()
	if stub != nil {
		fake.GetAllTimeEntryStub(arg1)
	}
}

func (fake *FakeHandlerInterface) GetAllTimeEntryCallCount() int {
	fake.getAllTimeEntryMutex.RLock()
	defer fake.getAllTimeEntryMutex.RUnlock()
	return len(fake.getAllTimeEntryArgsForCall)
}

func (fake *FakeHandlerInterface) GetAllTimeEntryCalls(stub func(*gin.Context)) {
	fake.getAllTimeEntryMutex.Lock()
	defer fake.getAllTimeEntryMutex.Unlock()
	fake.GetAllTimeEntryStub = stub
}

func (fake *FakeHandlerInterface) GetAllTimeEntryArgsForCall(i int) *gin.Context {
	fake.getAllTimeEntryMutex.RLock()
	defer fake.getAllTimeEntryMutex.RUnlock()
	argsForCall := fake.getAllTimeEntryArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeHandlerInterface) GetEmployeeHandler(arg1 *gin.Context) {
	fake.getEmployeeHandlerMutex.Lock()
	fake.getEmployeeHandlerArgsForCall = append(fake.getEmployeeHandlerArgsForCall, struct {
		arg1 *gin.Context
	}{arg1})
	stub := fake.GetEmployeeHandlerStub
	fake.recordInvocation("GetEmployeeHandler", []interface{}{arg1})
	fake.getEmployeeHandlerMutex.Unlock()
	if stub != nil {
		fake.GetEmployeeHandlerStub(arg1)
	}
}

func (fake *FakeHandlerInterface) GetEmployeeHandlerCallCount() int {
	fake.getEmployeeHandlerMutex.RLock()
	defer fake.getEmployeeHandlerMutex.RUnlock()
	return len(fake.getEmployeeHandlerArgsForCall)
}

func (fake *FakeHandlerInterface) GetEmployeeHandlerCalls(stub func(*gin.Context)) {
	fake.getEmployeeHandlerMutex.Lock()
	defer fake.getEmployeeHandlerMutex.Unlock()
	fake.GetEmployeeHandlerStub = stub
}

func (fake *FakeHandlerInterface) GetEmployeeHandlerArgsForCall(i int) *gin.Context {
	fake.getEmployeeHandlerMutex.RLock()
	defer fake.getEmployeeHandlerMutex.RUnlock()
	argsForCall := fake.getEmployeeHandlerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeHandlerInterface) GetTimeEntry(arg1 *gin.Context) {
	fake.getTimeEntryMutex.Lock()
	fake.getTimeEntryArgsForCall = append(fake.getTimeEntryArgsForCall, struct {
		arg1 *gin.Context
	}{arg1})
	stub := fake.GetTimeEntryStub
	fake.recordInvocation("GetTimeEntry", []interface{}{arg1})
	fake.getTimeEntryMutex.Unlock()
	if stub != nil {
		fake.GetTimeEntryStub(arg1)
	}
}

func (fake *FakeHandlerInterface) GetTimeEntryCallCount() int {
	fake.getTimeEntryMutex.RLock()
	defer fake.getTimeEntryMutex.RUnlock()
	return len(fake.getTimeEntryArgsForCall)
}

func (fake *FakeHandlerInterface) GetTimeEntryCalls(stub func(*gin.Context)) {
	fake.getTimeEntryMutex.Lock()
	defer fake.getTimeEntryMutex.Unlock()
	fake.GetTimeEntryStub = stub
}

func (fake *FakeHandlerInterface) GetTimeEntryArgsForCall(i int) *gin.Context {
	fake.getTimeEntryMutex.RLock()
	defer fake.getTimeEntryMutex.RUnlock()
	argsForCall := fake.getTimeEntryArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeHandlerInterface) UpdateTimeEntry(arg1 *gin.Context) {
	fake.updateTimeEntryMutex.Lock()
	fake.updateTimeEntryArgsForCall = append(fake.updateTimeEntryArgsForCall, struct {
		arg1 *gin.Context
	}{arg1})
	stub := fake.UpdateTimeEntryStub
	fake.recordInvocation("UpdateTimeEntry", []interface{}{arg1})
	fake.updateTimeEntryMutex.Unlock()
	if stub != nil {
		fake.UpdateTimeEntryStub(arg1)
	}
}

func (fake *FakeHandlerInterface) UpdateTimeEntryCallCount() int {
	fake.updateTimeEntryMutex.RLock()
	defer fake.updateTimeEntryMutex.RUnlock()
	return len(fake.updateTimeEntryArgsForCall)
}

func (fake *FakeHandlerInterface) UpdateTimeEntryCalls(stub func(*gin.Context)) {
	fake.updateTimeEntryMutex.Lock()
	defer fake.updateTimeEntryMutex.Unlock()
	fake.UpdateTimeEntryStub = stub
}

func (fake *FakeHandlerInterface) UpdateTimeEntryArgsForCall(i int) *gin.Context {
	fake.updateTimeEntryMutex.RLock()
	defer fake.updateTimeEntryMutex.RUnlock()
	argsForCall := fake.updateTimeEntryArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeHandlerInterface) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.creatTimeEntryMutex.RLock()
	defer fake.creatTimeEntryMutex.RUnlock()
	fake.createEmployeeHandlerMutex.RLock()
	defer fake.createEmployeeHandlerMutex.RUnlock()
	fake.deleteTimeEntryMutex.RLock()
	defer fake.deleteTimeEntryMutex.RUnlock()
	fake.getAllTimeEntryMutex.RLock()
	defer fake.getAllTimeEntryMutex.RUnlock()
	fake.getEmployeeHandlerMutex.RLock()
	defer fake.getEmployeeHandlerMutex.RUnlock()
	fake.getTimeEntryMutex.RLock()
	defer fake.getTimeEntryMutex.RUnlock()
	fake.updateTimeEntryMutex.RLock()
	defer fake.updateTimeEntryMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeHandlerInterface) recordInvocation(key string, args []interface{}) {
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

var _ routes.HandlerInterface = new(FakeHandlerInterface)
