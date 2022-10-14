package service_test

import (
	"bytes"
	"errors"
	"example-project/model"
	"example-project/service"
	"example-project/service/servicefakes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http/httptest"
	"testing"
	"time"
	//"example-project/service"
)

func TestDeleteTimeEntry(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.DeleteTimeEntryByIdReturns(&mongo.DeleteResult{DeletedCount: 1}, nil)
	serviceInstance := service.NewEmployeeService(fakeDB)
	result, err := serviceInstance.DeleteTimeEntries(primitive.NewObjectID().Hex(), time.Now())
	assert.Nil(t, err)
	assert.Equal(t, &mongo.DeleteResult{DeletedCount: 1}, result)
}
func TestGetTimeEntry_userId(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.GetTimeEntryByIDReturns([]model.TimeEntry{{UserId: "147", Start: time.Time{}, End: time.Time{}, BreakStart: time.Time{}, BreakEnd: time.Time{}, Project: "145"}})

	servicefakes := service.NewEmployeeService(fakeDB)
	result := servicefakes.GetTimeEntries("1")

	assert.Equal(t, []model.TimeEntry{{UserId: "147", Start: time.Time{}, End: time.Time{}, BreakStart: time.Time{}, BreakEnd: time.Time{}, Project: "145"}}, result)

}
func TestGetAll_TimeEntries(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.GetAllTimeEntryReturns([]model.TimeEntry{}, nil)

	servicefakes := service.NewEmployeeService(fakeDB)
	result, err := servicefakes.GetAllTimeEntries()
	assert.Nil(t, err)
	assert.Equal(t, []model.TimeEntry{}, result)

}

func TestColl_Timeentries(t *testing.T) {

	first, _ := time.Parse(time.RFC3339, "2021-08-02T08:00:00.801Z")
	end, _ := time.Parse(time.RFC3339, "2021-08-06T08:00:00.801Z")
	check := model.TimeEntry{Start: first, End: end}

	checkcomplit1, _ := time.Parse(time.RFC3339, "2021-08-03T08:00:00.801Z")
	checkcomplit2, _ := time.Parse(time.RFC3339, "2021-08-04T08:00:00.801Z")
	checkcomplit := model.TimeEntry{Start: checkcomplit1, End: checkcomplit2}

	overlabsleft1, _ := time.Parse(time.RFC3339, "2021-08-04T08:00:00.801Z")
	overlabsleft2, _ := time.Parse(time.RFC3339, "2021-08-07T08:00:00.801Z")
	overlabsleft := model.TimeEntry{Start: overlabsleft1, End: overlabsleft2}

	overlabright1, _ := time.Parse(time.RFC3339, "2021-08-01T08:00:00.801Z")
	overlabright2, _ := time.Parse(time.RFC3339, "2021-08-04T08:00:00.801Z")
	overlabright := model.TimeEntry{Start: overlabright1, End: overlabright2}

	contains1, _ := time.Parse(time.RFC3339, "2021-08-01T08:00:00.801Z")
	contains2, _ := time.Parse(time.RFC3339, "2021-08-07T08:00:00.801Z")
	contains := model.TimeEntry{Start: contains1, End: contains2}

	notcall1, _ := time.Parse(time.RFC3339, "2021-07-01T08:00:00.801Z")
	notcall2, _ := time.Parse(time.RFC3339, "2021-07-08T08:00:00.801Z")
	notcall := model.TimeEntry{Start: notcall1, End: notcall2}

	equalstart1, _ := time.Parse(time.RFC3339, "2021-08-02T08:00:00.801Z")
	equalstart2, _ := time.Parse(time.RFC3339, "2021-08-03T08:00:00.801Z")
	equalstart := model.TimeEntry{Start: equalstart1, End: equalstart2}

	equalend1, _ := time.Parse(time.RFC3339, "2021-08-01T08:00:00.801Z")
	equalend2, _ := time.Parse(time.RFC3339, "2021-08-07T08:00:00.801Z")
	equalend := model.TimeEntry{Start: equalend1, End: equalend2}

	var tests = []struct {
		checker  model.TimeEntry
		tester   model.TimeEntry
		expected bool
	}{
		{check, check, true},
		{check, checkcomplit, true},
		{check, overlabsleft, true},
		{check, overlabright, true},
		{check, contains, true},
		{check, notcall, false},
		{check, equalstart, true},
		{check, equalend, true},
	}
	for _, tt := range tests {
		fakeDB := &servicefakes.FakeDatabaseInterface{}
		servicefakes := service.NewEmployeeService(fakeDB)
		result := servicefakes.CollideTimeEntry(tt.checker, tt.tester)
		assert.Equal(t, result, tt.expected)
		assert.Nil(t, nil)
	}
}
func TestCreate_timeEntries_UserId(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.CreatTimeEntryByIdReturns(model.TimeEntry{}, nil)

	serviceInstance := service.NewEmployeeService(fakeDB)
	result, err := serviceInstance.CreatTimeEntries(model.TimeEntry{})

	assert.Equal(t, nil, err)
	assert.Equal(t, result, model.TimeEntry{})
}

func TestCreate_TimeEntries_if(t *testing.T) {

	var fakeJSONString = `
		[
			{
				"userId": "123456789",

			{
		]
	`

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	servicefakes := service.NewEmployeeService(fakeDB)
	body := bytes.NewBufferString(fakeJSONString)
	responseRecorder := httptest.NewRecorder()

	fakeerr := errors.New("fakeDB err")
	fakeDB.GetAllTimeEntryReturns(nil, fakeerr)
	result, err := servicefakes.CreatTimeEntries(model.TimeEntry{})
	responseRecorder.Body = body
	assert.Equal(t, fakeerr, err)
	assert.Nil(t, result)
}

func TestCreate_timeEntries_(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	servicefakes := service.NewEmployeeService(fakeDB)

	faketimeentries := []model.TimeEntry{
		{
			UserId: "1",
		},
	}
	fakeDB.GetAllTimeEntryReturns(faketimeentries, nil)
	result, err := servicefakes.CreatTimeEntries(model.TimeEntry{UserId: "2"})

	assert.Equal(t, nil, err)
	assert.Nil(t, result)
}
func TestCreate_timeEntries_coll(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	servicefakes := service.NewEmployeeService(fakeDB)

	faketimeentries := []model.TimeEntry{
		{
			UserId: "1",
		},
	}
	fakeDB.GetAllTimeEntryReturns(faketimeentries, nil)
	result, err := servicefakes.CreatTimeEntries(model.TimeEntry{UserId: "1"})

	assert.NotNil(t, err)
	assert.Nil(t, result)
}
func TestUpdateTimeEntry(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.UpdateTimeEntryByIdReturns(&mongo.UpdateResult{UpsertedCount: 1}, nil)

	servicefakes := service.NewEmployeeService(fakeDB)
	result, err := servicefakes.UpdateTimeEntries(model.TimeEntry{})

	assert.Equal(t, &mongo.UpdateResult{UpsertedCount: 1}, result)
	assert.Nil(t, err)

}
func TestUpdate_TimeEntry_Coll(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.UpdateTimeEntryByIdReturns(&mongo.UpdateResult{}, errors.New(""))

	servicefakes := service.NewEmployeeService(fakeDB)
	_, err := servicefakes.UpdateTimeEntries(model.TimeEntry{})

	assert.NotNil(t, err)

}
func TestUpdate_timeEntries_(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	servicefakes := service.NewEmployeeService(fakeDB)
	fakeresult := mongo.UpdateResult{}

	faketimeentries := []model.TimeEntry{
		{
			UserId: "1",
		},
	}
	fakeDB.GetAllTimeEntryReturns(faketimeentries, nil)
	fakeDB.UpdateTimeEntryByIdReturns(&fakeresult, nil)
	result, err := servicefakes.UpdateTimeEntries(model.TimeEntry{UserId: "2"})

	assert.Equal(t, nil, err)
	assert.Equal(t, &fakeresult, result)
}
func TestUpdate_timeEntries_coll(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	servicefakes := service.NewEmployeeService(fakeDB)
	fakeresult := mongo.UpdateResult{}
	starttime1, _ := time.Parse(time.RFC3339, "2021-08-01T08:00:00.801Z")
	starttime2, _ := time.Parse(time.RFC3339, "2021-08-01T09:00:00.801Z")
	endtime, _ := time.Parse(time.RFC3339, "2021-08-01T10:00:00.801Z")
	faketimeentries := []model.TimeEntry{
		{
			UserId: "1",
			Start:  starttime1, End: endtime,
		},
	}
	fakeDB.GetAllTimeEntryReturns(faketimeentries, nil)
	fakeDB.UpdateTimeEntryByIdReturns(&fakeresult, nil)

	_, err := servicefakes.UpdateTimeEntries(model.TimeEntry{UserId: "1", Start: starttime2, End: endtime})

	assert.NotNil(t, err)
}

/*func TestUpdate_TimeEntries_err(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.UpdateTimeEntryByIdReturns(&mongo.UpdateResult{UpsertedCount: 1}, nil)
	servicefakes := service.NewEmployeeService(fakeDB)
	result, _ := servicefakes.UpdateTimeEntries(model.TimeEntry{})
	assert.Equal(t, &mongo.UpdateResult{UpsertedCount: 1}, nil)
	assert.Nil(t, result)
}*/

func TestCalcultimeEntry_err(t *testing.T) {
	mockStartDate := time.Now()
	mockEndDate := time.Now().Add(time.Hour * 72)
	url := fmt.Sprintf("http://localhost:9090/test?start=%v&end=%v", mockStartDate.String(), mockEndDate.String())
	responseRecorder := httptest.NewRecorder()
	fakeContest, _ := gin.CreateTestContext(responseRecorder)
	fakeContest.Request = httptest.NewRequest("GET", url, nil)

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	servicefakes := service.NewEmployeeService(fakeDB)

	fakeerr := errors.New("fakeDB err")
	fakeDB.GetAllTimeEntryReturns(nil, fakeerr)
	result, err := servicefakes.CalculateTimeEntries(fakeContest)

	assert.Equal(t, fakeerr, err)
	assert.Nil(t, result)
}

func TestCalcultimeEntry(t *testing.T) {

	mockStartDate := time.Now()
	mockEndDate := time.Now().Add(time.Hour * 72)
	url := fmt.Sprintf("http://localhost:9090/test?start=%v&end=%v", mockStartDate.String(), mockEndDate.String())
	responseRecorder := httptest.NewRecorder()
	fakeContest, _ := gin.CreateTestContext(responseRecorder)
	fakeContest.Request = httptest.NewRequest("GET", url, nil)

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	servicefakes := service.NewEmployeeService(fakeDB)

	//fakeerr := errors.New("fakeDB err")
	faketimeentries := []model.TimeEntry{
		model.TimeEntry{
			UserId: "123", Start: time.Time{}, End: time.Time{}, BreakStart: time.Time{}, BreakEnd: time.Time{}, Project: "135"},
	}
	fakeDB.GetAllTimeEntryReturns(faketimeentries, nil)
	result, err := servicefakes.CalculateTimeEntries(fakeContest)

	assert.Equal(t, nil, err)
	assert.NotNil(t, result)
}
func TestCalcultimeEntryend(t *testing.T) {

	mockStartDate := time.Now()
	mockEndDate := time.Now().Add(time.Hour * 72)
	url := fmt.Sprintf("http://localhost:9090/test?start=%v&end=%v", mockStartDate.String(), mockEndDate.String())
	responseRecorder := httptest.NewRecorder()
	fakeContest, _ := gin.CreateTestContext(responseRecorder)
	fakeContest.Request = httptest.NewRequest("GET", url, nil)

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	servicefakes := service.NewEmployeeService(fakeDB)

	//fakeerr := errors.New("fakeDB err")
	faketimeentries := []model.TimeEntry{
		model.TimeEntry{
			UserId: "1", Start: time.Time{}, End: time.Time{}, BreakStart: time.Time{}, BreakEnd: time.Time{}, Project: "135"},
	}
	fakeDB.GetAllTimeEntryReturns(faketimeentries, nil)
	result, err := servicefakes.CalculateTimeEntries(fakeContest)

	assert.Equal(t, nil, err)
	assert.NotNil(t, result)
}
func TestCalcul_timeEntry_end(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	servicefakes := service.NewEmployeeService(fakeDB)

	mockStartDate := time.Now()
	mockEndDate := time.Now().Add(time.Hour * 72)
	url := fmt.Sprintf("http://localhost:9090/test?start=%v&end=%v", mockStartDate.String(), mockEndDate.String())
	responseRecorder := httptest.NewRecorder()
	fakeContest, _ := gin.CreateTestContext(responseRecorder)
	fakeContest.Request = httptest.NewRequest("GET", url, nil)

	//fakeerr := errors.New("fakeDB err")
	faketimeentries := []model.TimeEntry{
		model.TimeEntry{
			UserId: "1", Start: time.Time{}, End: time.Time{}, BreakStart: time.Time{}, BreakEnd: time.Time{}, Project: "135"},
	}
	fakeDB.GetAllTimeEntryReturns(faketimeentries, nil)
	result, err := servicefakes.CalculateTimeEntries(fakeContest)

	assert.Equal(t, nil, err)
	assert.NotNil(t, result)
}

func TestEmployeeService_CalculateTimeEntries(t *testing.T) {

	fakeHexId, _ := primitive.ObjectIDFromHex("6346c9d1b8489ecce7c010f8")

	timeEntryReturn := []model.TimeEntry{
		model.TimeEntry{UserId: "1"},
	}
	proposalReturn := []model.Proposal{
		model.Proposal{UserId: "1"},
	}
	getUserReturn := model.UserPayload{ID: fakeHexId}

	var tests = []struct {
		hasUserIdQuery          bool
		hasUserIdFormatationErr bool
		proposalReturn          []model.Proposal
		hasProposalReturnErr    bool
		proposalReturnErr       error
		timeEntryReturn         []model.TimeEntry
		hasTimeEntryReturnErr   bool
		timeEntryReturnErr      error
		getUserByIdReturn       model.UserPayload
		hasUserReturnErr        bool
		userByIdReturnErr       error
	}{
		{true, false, proposalReturn, false, nil, timeEntryReturn, false, nil, getUserReturn, false, nil},
		{false, false, proposalReturn, false, nil, timeEntryReturn, false, nil, getUserReturn, false, nil},
		{true, true, proposalReturn, false, nil, timeEntryReturn, false, nil, getUserReturn, false, nil},
		{true, false, nil, true, errors.New("fakeProposalReturnErr"), timeEntryReturn, false, nil, getUserReturn, false, nil},
		{true, false, proposalReturn, false, nil, timeEntryReturn, true, errors.New("fakeTimeEntryErr"), getUserReturn, false, nil},
	}
	for _, tt := range tests {
		fakeDB := &servicefakes.FakeDatabaseInterface{}
		serviceInstance := service.NewEmployeeService(fakeDB)
		fakeDB.GetUserByIDReturns(tt.getUserByIdReturn, nil)
		fakeDB.GetProposalsByFilterReturns(tt.proposalReturn, nil)
		fakeDB.GetTimeEntriesByFilterReturns(tt.timeEntryReturn, nil)
		if tt.hasUserIdQuery && !tt.hasUserIdFormatationErr {
			responseRecoder := httptest.NewRecorder()
			fakeContest, _ := gin.CreateTestContext(responseRecoder)
			fakeContest.Request = httptest.NewRequest("POST", "http://localhost:9090/user?userid=6346c9d1b8489ecce7c010f8", nil)

			_, err := serviceInstance.CalculateTimeEntries(fakeContest)
			assert.Equal(t, err, nil)
		}

		if !tt.hasUserIdQuery {
			responseRecoder := httptest.NewRecorder()
			fakeContest, _ := gin.CreateTestContext(responseRecoder)
			fakeContest.Request = httptest.NewRequest("POST", "http://localhost:9090/user", nil)
			expectedErr := "no user id supplied"
			_, err := serviceInstance.CalculateTimeEntries(fakeContest)
			assert.Contains(t, err.Error(), expectedErr)
		}

		if tt.hasUserIdFormatationErr {
			responseRecoder := httptest.NewRecorder()
			fakeContest, _ := gin.CreateTestContext(responseRecoder)
			fakeContest.Request = httptest.NewRequest("POST", "http://localhost:9090/user?userid=1", nil)
			expectedErr := "odd length hex string"
			_, err := serviceInstance.CalculateTimeEntries(fakeContest)
			assert.Contains(t, err.Error(), expectedErr)
		}

		if tt.hasProposalReturnErr {
			fakeDB.GetProposalsByFilterReturns(tt.proposalReturn, tt.proposalReturnErr)
			responseRecoder := httptest.NewRecorder()
			fakeContest, _ := gin.CreateTestContext(responseRecoder)
			fakeContest.Request = httptest.NewRequest("POST", "http://localhost:9090/user?userid=6346c9d1b8489ecce7c010f8", nil)
			expectedErr := "fakeProposalReturnErr"
			_, err := serviceInstance.CalculateTimeEntries(fakeContest)
			assert.Contains(t, err.Error(), expectedErr)
		}

		if tt.hasTimeEntryReturnErr {
			fakeDB.GetTimeEntriesByFilterReturns(tt.timeEntryReturn, tt.timeEntryReturnErr)
			responseRecoder := httptest.NewRecorder()
			fakeContest, _ := gin.CreateTestContext(responseRecoder)
			fakeContest.Request = httptest.NewRequest("POST", "http://localhost:9090/user?userid=6346c9d1b8489ecce7c010f8", nil)
			expectedErr := "fakeTimeEntryErr"
			_, err := serviceInstance.CalculateTimeEntries(fakeContest)
			assert.Contains(t, err.Error(), expectedErr)
		}

	}

}
