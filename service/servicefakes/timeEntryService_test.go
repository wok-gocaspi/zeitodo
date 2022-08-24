package servicefakes_test

import (
	"bytes"
	"errors"
	"example-project/handler/handlerfakes"
	"example-project/model"
	"example-project/service"
	"example-project/service/servicefakes"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDeleteTimeEntry(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.DeleteTimeEntryByIdReturns(&mongo.DeleteResult{DeletedCount: 1}, nil)
	serviceInstance := service.NewEmployeeService(fakeDB)
	result, err := serviceInstance.DeleteTimeEntries(primitive.NewObjectID().Hex())
	assert.Nil(t, err)
	assert.Equal(t, &mongo.DeleteResult{DeletedCount: 1}, result)
}
func TestGetTimeEntry_userId(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.GetTimeEntryByIDReturns([]model.TimeEntry{})

	servicefakes := service.NewEmployeeService(fakeDB)
	result := servicefakes.GetTimeEntries("1")

	assert.Equal(t, model.TimeEntry{UserId: "147", Start: time.Time{}, End: time.Time{}, BreakStart: time.Time{}, BreakEnd: time.Time{}, Project: "145"}, result)

}
func TestGetAll_TimeEntries(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.GetAllTimeEntryReturns([]model.TimeEntry{}, nil)

	servicefakes := service.NewEmployeeService(fakeDB)
	result, err := servicefakes.GetAllTimeEntries()
	assert.Nil(t, err)
	assert.Equal(t, model.TimeEntry{}, result)

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

	assert.Equal(t, result, err)
	assert.Nil(t, nil)
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
	fakeDB.CreatTimeEntryByIdReturns(model.TimeEntry{}, nil)
	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.CreatTimeEntriesReturns(nil, errors.New("insufficent user data"))
	serviceInstance := service.NewEmployeeService(fakeDB)
	result, err := serviceInstance.CreatTimeEntries(model.TimeEntry{})
	assert.Nil(t, err)
	assert.Equal(t, model.TimeEntry{}, result)

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
	result, err := servicefakes.UpdateTimeEntries(model.TimeEntry{})

	assert.Equal(t, &mongo.UpdateResult{UpsertedCount: 1}, result)
	assert.Nil(t, err)

}
func TestUpdate_TimeEntries_err(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.UpdateTimeEntryByIdReturns(&mongo.UpdateResult{UpsertedCount: 1}, nil)
	servicefakes := service.NewEmployeeService(fakeDB)
	result, _ := servicefakes.UpdateTimeEntries(model.TimeEntry{})
	assert.Equal(t, &mongo.UpdateResult{UpsertedCount: 1}, nil)
	assert.Nil(t, result)
}

func TestCalcultimeEntry_err(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	servicefakes := service.NewEmployeeService(fakeDB)

	fakeerr := errors.New("fakeDB err")
	fakeDB.GetAllTimeEntryReturns(nil, fakeerr)
	result, err := servicefakes.CalcultimeEntry("1")

	assert.Equal(t, fakeerr, err)
	assert.Nil(t, result)
}

func TestCalcultimeEntry(t *testing.T) {

	fakestart1, _ := time.Parse(time.RFC3339, "2021-08-01T08:00:00.801Z")
	fakeend1, _ := time.Parse(time.RFC3339, "2021-08-01T17:00:00.801Z")

	faketimeentries := []model.TimeEntry{
		{
			UserId: "123456789",
			Start:  fakestart1,
			End:    fakeend1,
		},
	}
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	servicefakes := service.NewEmployeeService(fakeDB)

	fakeerr := errors.New("fakeDB err")
	fakeDB.GetAllTimeEntryReturns(nil, fakeerr)
	result, err := servicefakes.CalcultimeEntry("1")

	assert.Equal(t, faketimeentries, err)
	assert.Nil(t, result)
}
