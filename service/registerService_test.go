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

func TestEmployeeService_UpdateEmployee(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}

	mockProposal := model.Proposal{
		UserId: "1", StartDate: "2006-Nov-06", EndDate: "2006-Nov-02", Approved: false}
	mockError := errors.New("fake error")
	result := &mongo.UpdateResult{}

	var tests = []struct {
		Proposal  model.Proposal
		Date      string
		hasError  bool
		mockError error
		Result    *mongo.UpdateResult
	}{
		{mockProposal, mockProposal.StartDate, false, nil, result},
		{mockProposal, mockProposal.StartDate, true, mockError, result},
	}

	for _, tt := range tests {
		fakeDB.UpdateProposalReturns(tt.Result, tt.mockError)
		serviceInstance := service.NewEmployeeService(fakeDB)

		actual, err := serviceInstance.UpdateProposalByDate(mockProposal, mockProposal.StartDate)

		assert.Equal(t, actual, tt.Result)
		assert.Equal(t, tt.mockError, err)

	}

}

func TestEmployeeService_CreateProposals(t *testing.T) {
	StartExceedsEnd := []model.ProposalPayload{
		model.ProposalPayload{UserId: "1", StartDate: "2006-Nov-08", EndDate: "2005-Nov-07"},
	}
	okPayload := []model.ProposalPayload{
		model.ProposalPayload{UserId: "1", StartDate: "2006-Nov-04", EndDate: "2006-Nov-10"},
	}
	GetByIdReturn := []model.Proposal{
		model.Proposal{UserId: "1", StartDate: "2006-Nov-06", EndDate: "2006-Nov-07"},
	}
	GetByIdReturnOverlapp := []model.Proposal{
		model.Proposal{UserId: "1", StartDate: "2006-Nov-07", EndDate: "2006-Nov-08"},
		model.Proposal{UserId: "1", StartDate: "2006-Nov-06", EndDate: "2006-Nov-09"},
	}

	StartExceedsEndMsg := "The startdate must be before the enddate"
	overlappingErrorMsg := "There cant be overlapping proposals."
	var tests = []struct {
		Payload                 []model.ProposalPayload
		GetProposalsReturn      []model.Proposal
		startDateExceedsEndDate bool
		overlappingErr          bool
		expectedError           string
	}{
		{StartExceedsEnd, GetByIdReturn, true, false, StartExceedsEndMsg},
		{okPayload, GetByIdReturnOverlapp, false, true, overlappingErrorMsg},
		{okPayload, GetByIdReturn, false, false, overlappingErrorMsg},
	}

	for _, tt := range tests {
		fakeDb := &servicefakes.FakeDatabaseInterface{}
		serviceInstance := service.NewEmployeeService(fakeDb)

		if tt.startDateExceedsEndDate {
			actual, err := serviceInstance.CreateProposals(tt.Payload, "2006-Nov-07")
			assert.Equal(t, nil, actual)
			assert.Equal(t, errors.New(tt.expectedError), err)
		}

		if tt.overlappingErr {
			fakeDb.GetProposalsReturns(tt.GetProposalsReturn, nil)

			actual, err := serviceInstance.CreateProposals(tt.Payload, "2006-Nov-07")
			assert.Equal(t, nil, actual)
			assert.Equal(t, errors.New(tt.expectedError), err)
		}
		if !tt.overlappingErr && !tt.startDateExceedsEndDate {
			fakeDb.GetProposalsReturns(tt.GetProposalsReturn, nil)

			actual, err := serviceInstance.CreateProposals(tt.Payload, "2010-Jan-07")
			assert.Equal(t, nil, actual)
			assert.Equal(t, errors.New(tt.expectedError), err)
		}

	}
}
