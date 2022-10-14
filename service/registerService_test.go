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
	"time"
)

func TestProposalService_GetProposalsByID(t *testing.T) {
	fakeDb := &servicefakes.FakeDatabaseInterface{}
	fakePayload := []model.Proposal{
		{UserId: "1", Status: "approved"},
		{UserId: "1", Status: "denied"},
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
		fakeDb.GetProposalsByUserIDReturns(tt.payload, tt.err)
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
	mockStartDate := time.Now()
	mockEndDate := time.Now().Add(time.Hour * 72)
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

		UserId: fakeuserid, StartDate: mockStartDate, EndDate: mockEndDate, Status: "denied"}

	mockProposalAdmin := model.Proposal{
		UserId: fakeAdminId, StartDate: mockStartDate, EndDate: mockStartDate, Status: "approved"}

	mockError := errors.New("fake userId")
	mockErrorUser := errors.New("user can not update ")

	result := &mongo.UpdateResult{}

	var tests = []struct {
		Proposal  model.Proposal
		Date      time.Time
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

		actual, err := serviceInstance.UpdateProposalByDate(tt.Proposal, tt.Proposal.StartDate.String(), tt.context)

		assert.Equal(t, actual, tt.Result)
		assert.Equal(t, tt.mockError, err)

	}

}

func TestEmployeeService_CreateProposals(t *testing.T) {
	mockStartDate := time.Now()
	mockEndDate := time.Now().Add(time.Hour * 72)
	StartExceedsEnd := []model.ProposalPayload{
		{UserId: "1", StartDate: "2006-Nov-08", EndDate: "2005-Nov-07"},
	}
	okPayload := []model.ProposalPayload{
		{UserId: "1", StartDate: "2006-Nov-04", EndDate: "2006-Nov-10"},
	}
	GetByIdReturn := []model.Proposal{
		{UserId: "1", StartDate: mockStartDate, EndDate: mockEndDate},
	}
	GetByIdReturnOverlapp := []model.Proposal{
		{UserId: "1", StartDate: mockStartDate, EndDate: mockEndDate},
		{UserId: "1", StartDate: mockStartDate, EndDate: mockEndDate},
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
			fakeDb.GetProposalsByUserIDReturns(tt.GetProposalsReturn, nil)

			actual, err := serviceInstance.CreateProposals(tt.Payload, "2006-Nov-07")
			assert.Equal(t, nil, actual)
			assert.Equal(t, errors.New(tt.expectedError), err)
		}
		if !tt.overlappingErr && !tt.startDateExceedsEndDate {
			fakeDb.GetProposalsByUserIDReturns(tt.GetProposalsReturn, nil)

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

func TestGetAllProposals(t *testing.T) {
	var exampleUserID1 = primitive.NewObjectID()
	var exampleUserID2 = primitive.NewObjectID()
	mockStartDate := time.Now()

	type ProposalCalls struct {
		Proposals []model.Proposal
		Error     error
	}

	var tests = []struct {
		GetAllUserReturn   []model.UserPayload
		GetAllUserError    error
		GetProposalsReturn []ProposalCalls
		Return             []model.ProposalsByUser
		Error              error
	}{
		{GetAllUserReturn: []model.UserPayload{
			{Username: "peter", ID: exampleUserID1},
			{Username: "frank", ID: exampleUserID2},
		}, GetProposalsReturn: []ProposalCalls{
			{Proposals: []model.Proposal{
				{StartDate: mockStartDate, Type: "vacation"},
				{StartDate: mockStartDate, Type: "sickness"},
			}, Error: nil},
			{Proposals: []model.Proposal{
				{UserId: exampleUserID2.Hex()},
			}, Error: nil},
		}, Return: []model.ProposalsByUser{
			{Username: "peter", Userid: exampleUserID1, VacationProposals: []model.Proposal{{StartDate: mockStartDate, Type: "vacation"}}},
			{Username: "frank", Userid: exampleUserID2, VacationProposals: []model.Proposal{{StartDate: mockStartDate, Type: "sickness"}}},
		}, Error: nil},
		{
			GetAllUserReturn:   []model.UserPayload{},
			GetAllUserError:    errors.New("no Users exist"),
			GetProposalsReturn: []ProposalCalls{},
			Return:             nil,
			Error:              errors.New("no Users exist"),
		},
		{
			GetAllUserReturn: []model.UserPayload{
				{Username: "peter1", ID: exampleUserID1},
				{Username: "frank1", ID: exampleUserID2},
			},
			GetAllUserError: nil,
			GetProposalsReturn: []ProposalCalls{
				{
					Proposals: []model.Proposal{},
					Error:     nil,
				},

				{
					Proposals: []model.Proposal{},
					Error:     errors.New("some db error"),
				},
			},
			Return: nil,
			Error:  errors.New("some db error"),
		},
	}
	for _, tt := range tests {
		fakeDB := &servicefakes.FakeDatabaseInterface{}
		serviceInstance := service.NewEmployeeService(fakeDB)
		fakeDB.GetAllUserReturns(tt.GetAllUserReturn, tt.GetAllUserError)
		for index, pCall := range tt.GetProposalsReturn {
			fakeDB.GetProposalsByFilterReturnsOnCall(index, pCall.Proposals, pCall.Error)
		}
		ctx := gin.Context{}
		result, err := serviceInstance.GetAllProposals(&ctx)
		assert.Equal(t, result, result)
		assert.Equal(t, err, err)

	}
}

func TestGetTotalAbsence(t *testing.T) {
	mockStartDate := time.Now()
	mockEndDate := time.Now().Add(time.Hour * 72)
	userID := primitive.NewObjectID().Hex()
	currentTime := time.Now()
	var tests = []struct {
		UserID            string
		GetUserByIDReturn model.UserPayload
		GetUserByIDError  error
		GetProposals      []model.Proposal
		GetProposalError  error
		Return            model.AbsenceObject
		Error             error
	}{
		{
			UserID: userID,
			GetUserByIDReturn: model.UserPayload{
				VacationDays: 30,
				EntryTime:    time.Date(currentTime.Year(), 8, 1, 0, 0, 0, 0, time.UTC),
			},
			GetUserByIDError: nil,
			GetProposals: []model.Proposal{
				{StartDate: mockStartDate, EndDate: mockEndDate, Type: "sickness"},
				{StartDate: mockStartDate, EndDate: mockEndDate, Type: "vacation"},
			},
			GetProposalError: nil,
			Return: model.AbsenceObject{
				TotalVacationDays: 10,
				VacationDays:      3,
				SicknessDays:      3,
			},
			Error: nil,
		},
		{
			UserID:            userID,
			GetUserByIDReturn: model.UserPayload{},
			GetUserByIDError:  errors.New("some user error"),
			GetProposals:      nil,
			GetProposalError:  nil,
			Return:            model.AbsenceObject{},
			Error:             errors.New("some user error"),
		},
		{
			UserID: userID,
			GetUserByIDReturn: model.UserPayload{
				VacationDays: 30,
				EntryTime:    time.Date(currentTime.Year(), 8, 1, 0, 0, 0, 0, time.UTC),
			},
			GetUserByIDError: nil,
			GetProposals:     nil,
			GetProposalError: errors.New("some proposal service error"),
			Return:           model.AbsenceObject{},
			Error:            errors.New("some proposal service error"),
		},
		{
			UserID: userID,
			GetUserByIDReturn: model.UserPayload{
				VacationDays: 30,
				EntryTime:    time.Date(currentTime.Year(), 8, 1, 0, 0, 0, 0, time.UTC),
			},
			GetUserByIDError: nil,
			GetProposals: []model.Proposal{
				{StartDate: mockStartDate, EndDate: mockEndDate, Type: "sickness"},
			},
			GetProposalError: nil,
			Return: model.AbsenceObject{
				TotalVacationDays: 0,
				VacationDays:      0,
				SicknessDays:      0,
			},
			Error: &time.ParseError{Layout: "2006-Jan-02", Value: "2022--10", LayoutElem: "Jan", ValueElem: "-10", Message: ""},
		},
		{
			UserID: userID,
			GetUserByIDReturn: model.UserPayload{
				VacationDays: 30,
				EntryTime:    time.Date(currentTime.Year(), 8, 1, 0, 0, 0, 0, time.UTC),
			},
			GetUserByIDError: nil,
			GetProposals: []model.Proposal{
				{StartDate: mockStartDate, EndDate: mockEndDate, Type: "vacation"},
			},
			GetProposalError: nil,
			Return: model.AbsenceObject{
				TotalVacationDays: 0,
				VacationDays:      0,
				SicknessDays:      0,
			},
			Error: &time.ParseError{Layout: "2006-Jan-02", Value: "2022-Sep15", LayoutElem: "-", ValueElem: "15", Message: ""},
		},
		{
			UserID: "",
			GetUserByIDReturn: model.UserPayload{
				VacationDays: 30,
				EntryTime:    time.Date(currentTime.Year(), 8, 1, 0, 0, 0, 0, time.UTC),
			},
			GetUserByIDError: nil,
			GetProposals:     []model.Proposal{},
			GetProposalError: nil,
			Return:           model.AbsenceObject{},
			Error:            errors.New("the provided hex string is not a valid ObjectID"),
		},
	}

	for _, tt := range tests {
		fakeDB := &servicefakes.FakeDatabaseInterface{}
		serviceInstance := service.NewEmployeeService(fakeDB)
		fakeDB.GetUserByIDReturns(tt.GetUserByIDReturn, tt.GetUserByIDError)
		fakeDB.GetProposalsByUserIDReturns(tt.GetProposals, tt.GetProposalError)
		result, err := serviceInstance.GetTotalAbsence(tt.UserID)
		assert.Equal(t, result, tt.Return)
		assert.Equal(t, err, tt.Error)
	}

}
