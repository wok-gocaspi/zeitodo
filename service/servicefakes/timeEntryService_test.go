package servicefakes_test

import (
	"errors"
	"example-project/handler/handlerfakes"
	"example-project/model"
	"example-project/service"
	"example-project/service/servicefakes"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	assert.Equal(t, model.TimeEntry{UserId: "147", Start: time.Time{}, End: time.Time{}, BreakStart: time.Time{}, BreakEnd: time.Time{}, Project: "145"}, result)

}
func TestCreate_timeEntries_UserId(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.CreatTimeEntryByIdReturns(model.TimeEntry{}, nil)

	serviceInstance := service.NewEmployeeService(fakeDB)
	result, err := serviceInstance.CreatTimeEntries(model.TimeEntry{UserId: "123", Start: time.Time{}, End: time.Time{}, BreakStart: time.Time{}, BreakEnd: time.Time{}, Project: "135"})
	assert.Nil(t, err)
	assert.Equal(t, model.TimeEntry{UserId: "147", Start: time.Time{}, End: time.Time{}, BreakStart: time.Time{}, BreakEnd: time.Time{}, Project: "145"}, result)

}
func TestCreate_timeEntries_(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.CreatTimeEntryByIdReturns(model.TimeEntry{}, nil)
	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.CreatTimeEntriesReturns(nil, errors.New("insufficent user data"))
	serviceInstance := service.NewEmployeeService(fakeDB)
	result, err := serviceInstance.CreatTimeEntries(model.TimeEntry{UserId: "123", Start: time.Time{}, End: time.Time{}, BreakStart: time.Time{}, BreakEnd: time.Time{}, Project: "135"})
	assert.Nil(t, err)
	assert.Equal(t, model.TimeEntry{UserId: "147", Start: time.Time{}, End: time.Time{}, BreakStart: time.Time{}, BreakEnd: time.Time{}, Project: "145"}, result)

}
func TestUpdateTimeEntry(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.UpdateTimeEntryByIdReturns(&mongo.UpdateResult{UpsertedCount: 1}, nil)

	servicefakes := service.NewEmployeeService(fakeDB)
	result, err := servicefakes.UpdateTimeEntries(model.TimeEntry{UserId: "123", Start: time.Time{}, End: time.Time{}, BreakStart: time.Time{}, BreakEnd: time.Time{}, Project: "135"})

	assert.Equal(t, &mongo.UpdateResult{UpsertedCount: 1}, result)
	assert.Nil(t, err)

}
func TestUpdate_TimeEntry_Coll(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.UpdateTimeEntryByIdReturns(&mongo.UpdateResult{}, errors.New(""))

	servicefakes := service.NewEmployeeService(fakeDB)
	result := servicefakes.CollideTimeEntry(model.TimeEntry{UserId: "123", Start: time.Time{}, End: time.Time{}, BreakStart: time.Time{}, BreakEnd: time.Time{}, Project: "135"}, model.TimeEntry{UserId: "123", Start: time.Time{}, End: time.Time{}, BreakStart: time.Time{}, BreakEnd: time.Time{}, Project: "135"})

	assert.Equal(t, &mongo.UpdateResult{UpsertedCount: 1}, result)
	assert.Nil(t, result)

}
func TestUpdate_TimeEntries_err(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.UpdateTimeEntryByIdReturns(&mongo.UpdateResult{UpsertedCount: 1}, nil)

}
func TestCollideTimeEntry(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}

}

func TestCalcultimeEntry(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}

}
