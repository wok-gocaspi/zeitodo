package service_test

import (
	"errors"
	"example-project/model"
	"example-project/service"
	"example-project/service/servicefakes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestProposalService_GetProposalsByID(t *testing.T) {
	fakeDb := &servicefakes.FakeDatabaseInterface{}
	fakePayload := []model.Proposal{
		{UserId: "1", Approved: false},
		{UserId: "1", Approved: true},
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
	fakeAdminId := primitive.NewObjectID().String()
	fakecontextAdmin := gin.Context{}
	fakecontextAdmin.Set("userid", primitive.NewObjectID().String())
	fakecontextAdmin.Set("group", "admin")
	fakecontextUser := gin.Context{}
	fakeuserid := primitive.NewObjectID().String()
	fakecontextUser.Set("userid", fakeuserid)
	fakecontextUser.Set("group", "user")

	mockProposaluser := model.Proposal{

		UserId: fakeuserid, StartDate: "2006-Nov-06", EndDate: "2006-Nov-02", Approved: false}

	mockProposalAdmin := model.Proposal{
		UserId: fakeAdminId, StartDate: "2006-Nov-06", EndDate: "2006-Nov-02", Approved: true}

	mockError := errors.New("fake userId")
	mockErrorUser := errors.New("user can not update ")

	result := &mongo.UpdateResult{}

	var tests = []struct {
		Proposal  model.Proposal
		Date      string
		hasError  bool
		mockError error
		Result    *mongo.UpdateResult
		context   *gin.Context
	}{
		{mockProposaluser, mockProposaluser.StartDate, false, nil, result, &fakecontextUser},
		{mockProposalAdmin, mockProposaluser.StartDate, true, mockErrorUser, nil, &fakecontextUser},
		{mockProposalAdmin, mockProposalAdmin.StartDate, false, nil, result, &fakecontextAdmin},
		{mockProposalAdmin, mockProposalAdmin.StartDate, true, mockError, result, &fakecontextAdmin},
	}

	for _, tt := range tests {
		fakeDB.UpdateProposalReturns(tt.Result, tt.mockError)
		serviceInstance := service.NewEmployeeService(fakeDB)

		actual, err := serviceInstance.UpdateProposalByDate(tt.Proposal, tt.Proposal.StartDate, tt.context)

		assert.Equal(t, actual, tt.Result)
		assert.Equal(t, tt.mockError, err)

	}

}

func TestEmployeeService_CreateProposals(t *testing.T) {
	StartExceedsEnd := []model.ProposalPayload{
		{UserId: "1", StartDate: "2006-Nov-08", EndDate: "2005-Nov-07"},
	}
	okPayload := []model.ProposalPayload{
		{UserId: "1", StartDate: "2006-Nov-04", EndDate: "2006-Nov-10"},
	}
	GetByIdReturn := []model.Proposal{
		{UserId: "1", StartDate: "2006-Nov-06", EndDate: "2006-Nov-07"},
	}
	GetByIdReturnOverlapp := []model.Proposal{
		{UserId: "1", StartDate: "2006-Nov-07", EndDate: "2006-Nov-08"},
		{UserId: "1", StartDate: "2006-Nov-06", EndDate: "2006-Nov-09"},
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
