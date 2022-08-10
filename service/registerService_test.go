package service_test

import (
	"errors"
	"example-project/model"
	"example-project/service"
	"example-project/service/servicefakes"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestGetEmployeeById(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}

	data := model.Employee{
		ID:        "1",
		FirstName: "jon",
		LastName:  "kock",
		Email:     "jon@gmail.com",
	}

	fakeDB.GetByIDReturns(data)

	serviceInstance := service.NewEmployeeService(fakeDB)
	actual := serviceInstance.GetEmployeeById("1")
	assert.Equal(t, data, actual)

}

func TestCreateEmployees(t *testing.T) {
	//here comes your first unit test which should cover the function CreateEmployees
}

func TestProposalService_GetProposalsByID(t *testing.T) {
	fakeDb := &servicefakes.FakeDatabaseInterface{}
	fakePayload := []model.Proposal{
		model.Proposal{UserId: "1", Approved: false},
		model.Proposal{UserId: "1", Approved: true},
	}
	fakeNilPayload := []model.Proposal{}
	fakeDecodeErr := errors.New("Decode went wrong")
	fakeNoResultErr := errors.New("No results could be found to your query")
	var tests = []struct {
		hasDecodeErr bool
		hasNoPayload bool
		payload      []model.Proposal
		err          error
	}{
		{false, false, fakePayload, nil},
		{true, false, fakeNilPayload, fakeDecodeErr},
		{false, true, fakeNilPayload, nil},
	}

	for _, tt := range tests {
		fakeDb.GetProposalsReturns(tt.payload, tt.err)
		serviceInstance := service.NewEmployeeService(fakeDb)

		if !tt.hasNoPayload && !tt.hasDecodeErr && tt.err == nil {
			actualResult, actualErr := serviceInstance.GetProposalsByID(tt.payload[0].UserId)
			assert.Equal(t, fakePayload, actualResult)
			assert.Equal(t, tt.err, actualErr)
		}
		if tt.hasDecodeErr {
			actualResult, actualErr := serviceInstance.GetProposalsByID("fakeDepartment")
			assert.Equal(t, tt.payload, actualResult)
			assert.Equal(t, tt.err, actualErr)
		}

		if tt.hasNoPayload {
			actualResult, actualErr := serviceInstance.GetProposalsByID("fakeDepartment")
			assert.Equal(t, tt.payload, actualResult)
			assert.Equal(t, fakeNoResultErr, actualErr)
		}
	}
}

func TestProposalService_DeleteProposalsByID(t *testing.T) {
	fakeDb := &servicefakes.FakeDatabaseInterface{}

	fakeDeletedCount0 := mongo.DeleteResult{
		DeletedCount: 0,
	}
	fakeDeletedCount1 := mongo.DeleteResult{
		DeletedCount: 1,
	}
	fakeDBErr := errors.New("Error in Database")
	fakeDeleterror := errors.New("the Employee id is not existing")

	var tests = []struct {
		hasValidID   bool
		hasValidDate bool
		hasDBErr     bool
		id           string
		date         string
		deletedCount *mongo.DeleteResult
		err          error
	}{
		{true, true, false, "1", "TODO", &fakeDeletedCount1, nil},
		{false, false, true, "1", "TODO", &fakeDeletedCount0, fakeDBErr},
		{true, true, true, "1", "TODO", &fakeDeletedCount1, fakeDBErr},
		{false, true, false, "1", "TODO", &fakeDeletedCount0, nil},
	}

	for _, tt := range tests {
		fakeDb.DeleteProposalByIdAndDateReturns(tt.deletedCount, tt.err)
		serviceInstance := service.NewEmployeeService(fakeDb)
		actualErr := serviceInstance.DeleteProposalsByID(tt.id, tt.date)

		if !tt.hasDBErr && tt.hasValidID && tt.hasValidDate {
			assert.Equal(t, nil, actualErr)
		}
		if tt.hasDBErr {
			assert.Equal(t, fakeDBErr, actualErr)
		}
		if !tt.hasDBErr && !tt.hasValidID {
			assert.Equal(t, fakeDeleterror, actualErr)
		}

	}
}
