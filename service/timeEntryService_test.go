package service_test

import (
	"example-project/model"
	"example-project/service"
	"example-project/service/servicefakes"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

func TestDeleteTimeEntry(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.DeleteTimeEntryByIdReturns(&mongo.DeleteResult{DeletedCount: 1}, nil)

	serviceInstance := service.NewEmployeeService(fakeDB)
	result, err := serviceInstance.DeleteTimeEntries("1")

	assert.Equal(t, &mongo.DeleteResult{DeletedCount: 1}, result)
	assert.Nil(t, err)
}
func TestCreattimeEntry(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.CreatTimeEntryByIdReturns(&mongo.InsertManyResult{}, nil)

	serviceInstance := service.NewEmployeeService(fakeDB)
	result, err := serviceInstance.CreatTimeEntries("1")

	assert.Equal(t, &mongo.InsertManyResult{}, result)
	assert.Nil(t, err)
}
func TestUpdateTimeEntry(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.UpdateTimeEntryByIdReturns(&mongo.UpdateResult{UpsertedCount: 1}, nil)

	servicefakes := service.NewEmployeeService(fakeDB)
	result, err := servicefakes.UpdateTimeEntries(model.TimeEntry{UserId: "123", Start: time.Time{}, End: time.Time{}, BreakStart: time.Time{}, BreakEnd: time.Time{}, Project: "135"})

	assert.Equal(t, &mongo.UpdateResult{UpsertedCount: 1}, result)
	assert.Nil(t, err)

}

func TestGetTimeEntry(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.GetTimeEntryByIdReturns(model.TimeEntry{UserId: "147", Start: time.Time{}, End: time.Time{}, BreakStart: time.Time{}, BreakEnd: time.Time{}, Project: "145"})

	servicefakes := service.NewEmployeeService(fakeDB)
	result := servicefakes.GetTimeEntries("1")

	assert.Equal(t, model.TimeEntry{UserId: "147", Start: time.Time{}, End: time.Time{}, BreakStart: time.Time{}, BreakEnd: time.Time{}, Project: "145"}, result)

}

func TestGetAllTimeEntriesTimeEntry(t *testing.T) {

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.GetAllTimeEntriesByIdReturns(model.TimeEntry{UserId: "147", Start: time.Time{}, End: time.Time{}, BreakStart: time.Time{}, BreakEnd: time.Time{}, Project: "145"})

	servicefakes := service.NewEmployeeService(fakeDB)
	result := servicefakes.GetAllTimeEntries("1")
	assert.Equal(t, model.TimeEntry{UserId: "147", Start: time.Time{}, End: time.Time{}, BreakStart: time.Time{}, BreakEnd: time.Time{}, Project: "145"}, result)

}
