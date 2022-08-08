package service_test

import (
	"example-project/model"
	"example-project/service/servicefakes"
	"testing"
)

func TestGetUserByID_Return_valid_200(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	dbReturn := model.UserPayload{FirstName: "Tom", LastName: "Cruz", Email: "bielefeld@gibtesnicht.de", Username: "tcruz", Password: "123"}
	fakeDB.GetUserByIDReturns(dbReturn, nil)
}
