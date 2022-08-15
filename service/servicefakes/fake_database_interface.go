// Code generated by counterfeiter. DO NOT EDIT.
package servicefakes

import (
	"example-project/model"
	"example-project/service"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

type FakeDatabaseInterface struct {
	DeleteProposalByIdAndDateStub        func(string, string) (*mongo.DeleteResult, error)
	deleteProposalByIdAndDateMutex       sync.RWMutex
	deleteProposalByIdAndDateArgsForCall []struct {
		arg1 string
		arg2 string
	}
	deleteProposalByIdAndDateReturns struct {
		result1 *mongo.DeleteResult
		result2 error
	}
	deleteProposalByIdAndDateReturnsOnCall map[int]struct {
		result1 *mongo.DeleteResult
		result2 error
	}
	GetByIDStub        func(string) model.Employee
	getByIDMutex       sync.RWMutex
	getByIDArgsForCall []struct {
		arg1 string
	}
	getByIDReturns struct {
		result1 model.Employee
	}
	getByIDReturnsOnCall map[int]struct {
		result1 model.Employee
	}
	GetProposalsStub        func(string) ([]model.Proposal, error)
	getProposalsMutex       sync.RWMutex
	getProposalsArgsForCall []struct {
		arg1 string
	}
	getProposalsReturns struct {
		result1 []model.Proposal
		result2 error
	}
	getProposalsReturnsOnCall map[int]struct {
		result1 []model.Proposal
		result2 error
	}
	SaveProposalsStub        func([]interface{}) (interface{}, error)
	saveProposalsMutex       sync.RWMutex
	saveProposalsArgsForCall []struct {
		arg1 []interface{}
	}
	saveProposalsReturns struct {
		result1 interface{}
		result2 error
	}
	saveProposalsReturnsOnCall map[int]struct {
		result1 interface{}
		result2 error
	}
	UpdateProposalStub        func(model.Proposal, string) (*mongo.UpdateResult, error)
	updateProposalMutex       sync.RWMutex
	updateProposalArgsForCall []struct {
		arg1 model.Proposal
		arg2 string
	}
	updateProposalReturns struct {
		result1 *mongo.UpdateResult
		result2 error
	}
	updateProposalReturnsOnCall map[int]struct {
		result1 *mongo.UpdateResult
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeDatabaseInterface) DeleteProposalByIdAndDate(arg1 string, arg2 string) (*mongo.DeleteResult, error) {
	fake.deleteProposalByIdAndDateMutex.Lock()
	ret, specificReturn := fake.deleteProposalByIdAndDateReturnsOnCall[len(fake.deleteProposalByIdAndDateArgsForCall)]
	fake.deleteProposalByIdAndDateArgsForCall = append(fake.deleteProposalByIdAndDateArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.DeleteProposalByIdAndDateStub
	fakeReturns := fake.deleteProposalByIdAndDateReturns
	fake.recordInvocation("DeleteProposalByIdAndDate", []interface{}{arg1, arg2})
	fake.deleteProposalByIdAndDateMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDatabaseInterface) DeleteProposalByIdAndDateCallCount() int {
	fake.deleteProposalByIdAndDateMutex.RLock()
	defer fake.deleteProposalByIdAndDateMutex.RUnlock()
	return len(fake.deleteProposalByIdAndDateArgsForCall)
}

func (fake *FakeDatabaseInterface) DeleteProposalByIdAndDateCalls(stub func(string, string) (*mongo.DeleteResult, error)) {
	fake.deleteProposalByIdAndDateMutex.Lock()
	defer fake.deleteProposalByIdAndDateMutex.Unlock()
	fake.DeleteProposalByIdAndDateStub = stub
}

func (fake *FakeDatabaseInterface) DeleteProposalByIdAndDateArgsForCall(i int) (string, string) {
	fake.deleteProposalByIdAndDateMutex.RLock()
	defer fake.deleteProposalByIdAndDateMutex.RUnlock()
	argsForCall := fake.deleteProposalByIdAndDateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeDatabaseInterface) DeleteProposalByIdAndDateReturns(result1 *mongo.DeleteResult, result2 error) {
	fake.deleteProposalByIdAndDateMutex.Lock()
	defer fake.deleteProposalByIdAndDateMutex.Unlock()
	fake.DeleteProposalByIdAndDateStub = nil
	fake.deleteProposalByIdAndDateReturns = struct {
		result1 *mongo.DeleteResult
		result2 error
	}{result1, result2}
}

func (fake *FakeDatabaseInterface) DeleteProposalByIdAndDateReturnsOnCall(i int, result1 *mongo.DeleteResult, result2 error) {
	fake.deleteProposalByIdAndDateMutex.Lock()
	defer fake.deleteProposalByIdAndDateMutex.Unlock()
	fake.DeleteProposalByIdAndDateStub = nil
	if fake.deleteProposalByIdAndDateReturnsOnCall == nil {
		fake.deleteProposalByIdAndDateReturnsOnCall = make(map[int]struct {
			result1 *mongo.DeleteResult
			result2 error
		})
	}
	fake.deleteProposalByIdAndDateReturnsOnCall[i] = struct {
		result1 *mongo.DeleteResult
		result2 error
	}{result1, result2}
}

func (fake *FakeDatabaseInterface) GetByID(arg1 string) model.Employee {
	fake.getByIDMutex.Lock()
	ret, specificReturn := fake.getByIDReturnsOnCall[len(fake.getByIDArgsForCall)]
	fake.getByIDArgsForCall = append(fake.getByIDArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetByIDStub
	fakeReturns := fake.getByIDReturns
	fake.recordInvocation("GetByID", []interface{}{arg1})
	fake.getByIDMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeDatabaseInterface) GetByIDCallCount() int {
	fake.getByIDMutex.RLock()
	defer fake.getByIDMutex.RUnlock()
	return len(fake.getByIDArgsForCall)
}

func (fake *FakeDatabaseInterface) GetByIDCalls(stub func(string) model.Employee) {
	fake.getByIDMutex.Lock()
	defer fake.getByIDMutex.Unlock()
	fake.GetByIDStub = stub
}

func (fake *FakeDatabaseInterface) GetByIDArgsForCall(i int) string {
	fake.getByIDMutex.RLock()
	defer fake.getByIDMutex.RUnlock()
	argsForCall := fake.getByIDArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeDatabaseInterface) GetByIDReturns(result1 model.Employee) {
	fake.getByIDMutex.Lock()
	defer fake.getByIDMutex.Unlock()
	fake.GetByIDStub = nil
	fake.getByIDReturns = struct {
		result1 model.Employee
	}{result1}
}

func (fake *FakeDatabaseInterface) GetByIDReturnsOnCall(i int, result1 model.Employee) {
	fake.getByIDMutex.Lock()
	defer fake.getByIDMutex.Unlock()
	fake.GetByIDStub = nil
	if fake.getByIDReturnsOnCall == nil {
		fake.getByIDReturnsOnCall = make(map[int]struct {
			result1 model.Employee
		})
	}
	fake.getByIDReturnsOnCall[i] = struct {
		result1 model.Employee
	}{result1}
}

func (fake *FakeDatabaseInterface) GetProposals(arg1 string) ([]model.Proposal, error) {
	fake.getProposalsMutex.Lock()
	ret, specificReturn := fake.getProposalsReturnsOnCall[len(fake.getProposalsArgsForCall)]
	fake.getProposalsArgsForCall = append(fake.getProposalsArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetProposalsStub
	fakeReturns := fake.getProposalsReturns
	fake.recordInvocation("GetProposals", []interface{}{arg1})
	fake.getProposalsMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDatabaseInterface) GetProposalsCallCount() int {
	fake.getProposalsMutex.RLock()
	defer fake.getProposalsMutex.RUnlock()
	return len(fake.getProposalsArgsForCall)
}

func (fake *FakeDatabaseInterface) GetProposalsCalls(stub func(string) ([]model.Proposal, error)) {
	fake.getProposalsMutex.Lock()
	defer fake.getProposalsMutex.Unlock()
	fake.GetProposalsStub = stub
}

func (fake *FakeDatabaseInterface) GetProposalsArgsForCall(i int) string {
	fake.getProposalsMutex.RLock()
	defer fake.getProposalsMutex.RUnlock()
	argsForCall := fake.getProposalsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeDatabaseInterface) GetProposalsReturns(result1 []model.Proposal, result2 error) {
	fake.getProposalsMutex.Lock()
	defer fake.getProposalsMutex.Unlock()
	fake.GetProposalsStub = nil
	fake.getProposalsReturns = struct {
		result1 []model.Proposal
		result2 error
	}{result1, result2}
}

func (fake *FakeDatabaseInterface) GetProposalsReturnsOnCall(i int, result1 []model.Proposal, result2 error) {
	fake.getProposalsMutex.Lock()
	defer fake.getProposalsMutex.Unlock()
	fake.GetProposalsStub = nil
	if fake.getProposalsReturnsOnCall == nil {
		fake.getProposalsReturnsOnCall = make(map[int]struct {
			result1 []model.Proposal
			result2 error
		})
	}
	fake.getProposalsReturnsOnCall[i] = struct {
		result1 []model.Proposal
		result2 error
	}{result1, result2}
}

func (fake *FakeDatabaseInterface) SaveProposals(arg1 []interface{}) (interface{}, error) {
	var arg1Copy []interface{}
	if arg1 != nil {
		arg1Copy = make([]interface{}, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.saveProposalsMutex.Lock()
	ret, specificReturn := fake.saveProposalsReturnsOnCall[len(fake.saveProposalsArgsForCall)]
	fake.saveProposalsArgsForCall = append(fake.saveProposalsArgsForCall, struct {
		arg1 []interface{}
	}{arg1Copy})
	stub := fake.SaveProposalsStub
	fakeReturns := fake.saveProposalsReturns
	fake.recordInvocation("SaveProposals", []interface{}{arg1Copy})
	fake.saveProposalsMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDatabaseInterface) SaveProposalsCallCount() int {
	fake.saveProposalsMutex.RLock()
	defer fake.saveProposalsMutex.RUnlock()
	return len(fake.saveProposalsArgsForCall)
}

func (fake *FakeDatabaseInterface) SaveProposalsCalls(stub func([]interface{}) (interface{}, error)) {
	fake.saveProposalsMutex.Lock()
	defer fake.saveProposalsMutex.Unlock()
	fake.SaveProposalsStub = stub
}

func (fake *FakeDatabaseInterface) SaveProposalsArgsForCall(i int) []interface{} {
	fake.saveProposalsMutex.RLock()
	defer fake.saveProposalsMutex.RUnlock()
	argsForCall := fake.saveProposalsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeDatabaseInterface) SaveProposalsReturns(result1 interface{}, result2 error) {
	fake.saveProposalsMutex.Lock()
	defer fake.saveProposalsMutex.Unlock()
	fake.SaveProposalsStub = nil
	fake.saveProposalsReturns = struct {
		result1 interface{}
		result2 error
	}{result1, result2}
}

func (fake *FakeDatabaseInterface) SaveProposalsReturnsOnCall(i int, result1 interface{}, result2 error) {
	fake.saveProposalsMutex.Lock()
	defer fake.saveProposalsMutex.Unlock()
	fake.SaveProposalsStub = nil
	if fake.saveProposalsReturnsOnCall == nil {
		fake.saveProposalsReturnsOnCall = make(map[int]struct {
			result1 interface{}
			result2 error
		})
	}
	fake.saveProposalsReturnsOnCall[i] = struct {
		result1 interface{}
		result2 error
	}{result1, result2}
}

func (fake *FakeDatabaseInterface) UpdateProposal(arg1 model.Proposal, arg2 string) (*mongo.UpdateResult, error) {
	fake.updateProposalMutex.Lock()
	ret, specificReturn := fake.updateProposalReturnsOnCall[len(fake.updateProposalArgsForCall)]
	fake.updateProposalArgsForCall = append(fake.updateProposalArgsForCall, struct {
		arg1 model.Proposal
		arg2 string
	}{arg1, arg2})
	stub := fake.UpdateProposalStub
	fakeReturns := fake.updateProposalReturns
	fake.recordInvocation("UpdateProposal", []interface{}{arg1, arg2})
	fake.updateProposalMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDatabaseInterface) UpdateProposalCallCount() int {
	fake.updateProposalMutex.RLock()
	defer fake.updateProposalMutex.RUnlock()
	return len(fake.updateProposalArgsForCall)
}

func (fake *FakeDatabaseInterface) UpdateProposalCalls(stub func(model.Proposal, string) (*mongo.UpdateResult, error)) {
	fake.updateProposalMutex.Lock()
	defer fake.updateProposalMutex.Unlock()
	fake.UpdateProposalStub = stub
}

func (fake *FakeDatabaseInterface) UpdateProposalArgsForCall(i int) (model.Proposal, string) {
	fake.updateProposalMutex.RLock()
	defer fake.updateProposalMutex.RUnlock()
	argsForCall := fake.updateProposalArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeDatabaseInterface) UpdateProposalReturns(result1 *mongo.UpdateResult, result2 error) {
	fake.updateProposalMutex.Lock()
	defer fake.updateProposalMutex.Unlock()
	fake.UpdateProposalStub = nil
	fake.updateProposalReturns = struct {
		result1 *mongo.UpdateResult
		result2 error
	}{result1, result2}
}

func (fake *FakeDatabaseInterface) UpdateProposalReturnsOnCall(i int, result1 *mongo.UpdateResult, result2 error) {
	fake.updateProposalMutex.Lock()
	defer fake.updateProposalMutex.Unlock()
	fake.UpdateProposalStub = nil
	if fake.updateProposalReturnsOnCall == nil {
		fake.updateProposalReturnsOnCall = make(map[int]struct {
			result1 *mongo.UpdateResult
			result2 error
		})
	}
	fake.updateProposalReturnsOnCall[i] = struct {
		result1 *mongo.UpdateResult
		result2 error
	}{result1, result2}
}

func (fake *FakeDatabaseInterface) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.deleteProposalByIdAndDateMutex.RLock()
	defer fake.deleteProposalByIdAndDateMutex.RUnlock()
	fake.getByIDMutex.RLock()
	defer fake.getByIDMutex.RUnlock()
	fake.getProposalsMutex.RLock()
	defer fake.getProposalsMutex.RUnlock()
	fake.saveProposalsMutex.RLock()
	defer fake.saveProposalsMutex.RUnlock()
	fake.updateProposalMutex.RLock()
	defer fake.updateProposalMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeDatabaseInterface) recordInvocation(key string, args []interface{}) {
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

var _ service.DatabaseInterface = new(FakeDatabaseInterface)
