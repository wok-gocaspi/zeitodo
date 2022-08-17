package service_test

import (
	"errors"
	"example-project/model"
	"example-project/service"
	"example-project/service/servicefakes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestGetUserByID_Success_ReturnsUser(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	dbReturn := model.UserPayload{}
	fakeDB.GetUserByIDReturns(dbReturn, nil)
	serviceInstance := service.NewEmployeeService(fakeDB)
	fakeID := primitive.NewObjectID()
	fakeIDString := fakeID.Hex()
	_, err := serviceInstance.GetUserByID(fakeIDString)
	assert.Nil(t, err)
}

func TestGetUserByID_InvalidID_ReturnsError(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	dbReturn := model.UserPayload{}
	fakeDB.GetUserByIDReturns(dbReturn, nil)
	serviceInstance := service.NewEmployeeService(fakeDB)
	_, err := serviceInstance.GetUserByID("1")
	assert.EqualError(t, err, "encoding/hex: odd length hex string")
}

func TestGetUserByID_Return_invalid_database_error(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	dbReturn := model.UserPayload{}
	fakeDB.GetUserByIDReturns(dbReturn, errors.New("mongo: no documents in result"))
	serviceInstance := service.NewEmployeeService(fakeDB)
	fakeID := primitive.NewObjectID()
	fakeIDString := fakeID.Hex()
	fmt.Println(fakeIDString)
	_, err := serviceInstance.GetUserByID(fakeIDString)
	assert.EqualError(t, err, "mongo: no documents in result")
}

func TestGetAllUser_Return_valid(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	var dbReturn []model.UserPayload
	fakeDB.GetAllUserReturns(dbReturn, nil)
	serviceInstance := service.NewEmployeeService(fakeDB)
	_, err := serviceInstance.GetAllUser()
	assert.Nil(t, err)
}

func TestGetAllUser_Return_invalid_database_error(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	var dbReturn []model.UserPayload
	fakeDB.GetAllUserReturns(dbReturn, errors.New("mongo: no documents in result"))
	serviceInstance := service.NewEmployeeService(fakeDB)
	_, err := serviceInstance.GetAllUser()
	assert.Error(t, err, "mongo: no documents in result")
}

func TestGetTeamMembersByUserID_Return_valid(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	var dbReturn []model.TeamMember
	fakeDB.GetUserTeamMembersByIDReturns(dbReturn, nil)
	serviceInstance := service.NewEmployeeService(fakeDB)
	fakeIDString := primitive.NewObjectID().Hex()
	_, err := serviceInstance.GetTeamMembersByUserID(fakeIDString)
	assert.Nil(t, err)
}

func TestGetTeamMembersByUserID_Return_invalid_userid(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	var dbReturn []model.TeamMember
	fakeDB.GetUserTeamMembersByIDReturns(dbReturn, nil)
	serviceInstance := service.NewEmployeeService(fakeDB)
	_, err := serviceInstance.GetTeamMembersByUserID("1")
	assert.EqualError(t, err, "encoding/hex: odd length hex string")
}

func TestGetTeamMembersByUserID_Return_invalid_database_error(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	var dbReturn []model.TeamMember
	fakeDB.GetUserTeamMembersByIDReturns(dbReturn, errors.New("assert.Error(t, err, \"mongo: no documents in result\")"))
	serviceInstance := service.NewEmployeeService(fakeDB)
	fakeIDString := primitive.NewObjectID().Hex()
	_, err := serviceInstance.GetTeamMembersByUserID(fakeIDString)
	assert.Error(t, err, "mongo: no documents in result")
}

func TestGetTeamMembersByName_Return_valid(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	var dbReturn []model.TeamMember
	fakeDB.GetUserTeamMembersByNameReturns(dbReturn, nil)
	serviceInstance := service.NewEmployeeService(fakeDB)
	_, err := serviceInstance.GetTeamMembersByName("test")
	assert.Nil(t, err)
}

func TestGetTeamMembersByName_Return_invalid_database_error(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	var dbReturn []model.TeamMember
	fakeDB.GetUserTeamMembersByNameReturns(dbReturn, errors.New("some db error"))
	serviceInstance := service.NewEmployeeService(fakeDB)
	_, err := serviceInstance.GetTeamMembersByName("test")
	assert.Error(t, err, "some db error")
}

func TestCreateUser_Return_valid(t *testing.T) {
	fakeUser := model.UserSignupPayload{Username: "pganz", Password: "123", FirstName: "Peter", LastName: "Ganz", Email: "p.ganz@gmail.com"}
	dbReturn := primitive.NewObjectID()
	var dbInterface interface{} = dbReturn
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.CreateUserReturns(dbInterface, nil)
	serviceInstance := service.NewEmployeeService(fakeDB)
	_, err := serviceInstance.CreateUser(fakeUser)
	assert.Nil(t, err)
}

func TestCreateUser_Return_invalid_database_error(t *testing.T) {
	fakeUser := model.UserSignupPayload{Username: "pganz", Password: "123", FirstName: "Peter", LastName: "Ganz", Email: "p.ganz@gmail.com"}

	var dbReturn interface{}
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.CreateUserReturns(dbReturn, errors.New("some db error"))
	serviceInstance := service.NewEmployeeService(fakeDB)
	_, err := serviceInstance.CreateUser(fakeUser)
	assert.Error(t, err, "some db error")
}

func TestCreateUser_Return_invalid_payload(t *testing.T) {
	fakeUser := model.UserSignupPayload{Username: "pganz", FirstName: "Peter", LastName: "Ganz", Email: "p.ganz@gmail.com"}
	dbReturn := primitive.NewObjectID()
	var dbInterface interface{} = dbReturn
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.CreateUserReturns(dbInterface, nil)
	serviceInstance := service.NewEmployeeService(fakeDB)
	_, err := serviceInstance.CreateUser(fakeUser)
	assert.Error(t, err, "insufficent user data")
}

func TestCreateUser_Return_existing_user(t *testing.T) {
	fakeUser := model.UserSignupPayload{Username: "pganz", FirstName: "Peter", LastName: "Ganz", Email: "p.ganz@gmail.com"}
	dbReturn := primitive.NewObjectID()
	var dbInterface interface{} = dbReturn
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.CreateUserReturns(dbInterface, nil)
	fakeDB.GetUserByUsernameReturns(model.UserPayload{Username: "123"}, nil)
	serviceInstance := service.NewEmployeeService(fakeDB)
	_, err := serviceInstance.CreateUser(fakeUser)
	assert.Error(t, err, "user already exists, please choose another username")
}

func TestUpdateUser_Return_valid(t *testing.T) {
	fakeUserArray := []model.User{
		{ID: primitive.NewObjectID(), LastName: "Hans"},
		{ID: primitive.NewObjectID(), Email: "test@gmail.com"},
	}
	fakeUserUpdateResults := []model.UserUpdateResult{
		{Success: true, User: fakeUserArray[0]},
		{Success: true, User: fakeUserArray[1]},
	}
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.UpdateManyUserByIDReturns(fakeUserUpdateResults)
	serviceInstance := service.NewEmployeeService(fakeDB)
	result, err := serviceInstance.UpdateUsers(fakeUserArray)
	assert.Nil(t, err)
	assert.Equal(t, fakeUserUpdateResults, result)
}

func TestUpdateUser_Return_Unsuccessful(t *testing.T) {
	fakeUserArray := []model.User{
		{ID: primitive.NewObjectID(), LastName: "Hans"},
		{ID: primitive.NewObjectID(), Email: "test@gmail.com"},
	}
	fakeUserUpdateResults := []model.UserUpdateResult{
		{Success: false, User: fakeUserArray[0]},
		{Success: true, User: fakeUserArray[1]},
	}
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.UpdateManyUserByIDReturns(fakeUserUpdateResults)
	serviceInstance := service.NewEmployeeService(fakeDB)
	result, err := serviceInstance.UpdateUsers(fakeUserArray)
	assert.Error(t, err, "a few users couldn't be updated")
	assert.Equal(t, []model.UserUpdateResult{{User: fakeUserUpdateResults[0].User, Success: false}}, result)
}

func TestDeleteUser_Return_Success(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.DeleteUserReturns(&mongo.DeleteResult{DeletedCount: 1}, nil)
	serviceInstance := service.NewEmployeeService(fakeDB)
	result, err := serviceInstance.DeleteUsers(primitive.NewObjectID().Hex())
	assert.Nil(t, err)
	assert.Equal(t, &mongo.DeleteResult{DeletedCount: 1}, result)
}

func TestDeleteUser_Invalid_ID(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.DeleteUserReturns(&mongo.DeleteResult{DeletedCount: 0}, nil)
	serviceInstance := service.NewEmployeeService(fakeDB)
	result, err := serviceInstance.DeleteUsers("1")
	assert.Nil(t, result)
	assert.Error(t, err, "encoding/hex: odd length hex string")
}
