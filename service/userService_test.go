package service_test

import (
	"crypto/sha256"
	"errors"
	"example-project/model"
	"example-project/routes"
	"example-project/service"
	"example-project/service/servicefakes"
	"example-project/utils"
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

func TestGetUserID_Success_ReturnsUserId(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	dbReturn := model.User{}
	fakeDB.GetUserByUsernameReturns(dbReturn, nil)
	serviceInstance := service.NewEmployeeService(fakeDB)

	fakeUsername := "Rafael"

	_, err := serviceInstance.GetUserId(fakeUsername)
	assert.Nil(t, err)
}

func TestGetUserID_InvalidUserName_ReturnsError(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	dbReturn := model.User{}
	fakeDB.GetUserByUsernameReturns(dbReturn, errors.New("no user found to that username"))
	serviceInstance := service.NewEmployeeService(fakeDB)

	fakeUsername := "Rafael"

	_, err := serviceInstance.GetUserId(fakeUsername)
	assert.Contains(t, err.Error(), "no user found to that username")
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
	dbInsertedIDS := primitive.NewObjectID()
	var dbReturn interface{} = dbInsertedIDS
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.CreateUserReturns(dbReturn, nil)
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
	fakeUser := model.UserSignupPayload{Username: "pganz", Password: "123", LastName: "Ganz", Email: "p.ganz@gmail.com"}
	var dbReturn interface{}
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	fakeDB.CreateUserReturns(dbReturn, nil)
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
	fakeDB.GetUserByUsernameReturns(model.User{Username: "123"}, nil)
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

func TestEmployeeService_LoginUser(t *testing.T) {
	fakePw := "pa55word"
	fakePwWrong := "password"
	fakeUserPwHash := sha256.Sum256([]byte(fakePw))
	fakeUser := model.User{ID: primitive.NewObjectID(), Username: "Hans", Password: fakeUserPwHash}
	fakeError := errors.New("invalid login")

	var tests = []struct {
		GetUserError   bool
		PasswordsError bool
	}{
		{true, false},
		{false, true},
		{false, false},
	}

	for _, tt := range tests {
		fakeDB := &servicefakes.FakeDatabaseInterface{}
		serviceInstance := service.NewEmployeeService(fakeDB)
		if tt.GetUserError {
			fakeDB.GetUserByUsernameReturns(fakeUser, fakeError)
			_, err := serviceInstance.LoginUser(fakeUser.Username, fakePw)

			assert.Equal(t, fakeError.Error(), err.Error())
		}

		if tt.PasswordsError {
			fakeDB.GetUserByUsernameReturns(fakeUser, nil)
			_, err := serviceInstance.LoginUser(fakeUser.Username, fakePwWrong)

			assert.Equal(t, fakeError.Error(), err.Error())
		}

		fakeDB.GetUserByUsernameReturns(fakeUser, nil)
		_, err := serviceInstance.LoginUser(fakeUser.Username, fakePw)

		assert.Nil(t, err)
	}
}

func TestEmployeeService_LogoutUserF(t *testing.T) {
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	serviceInstance := service.NewEmployeeService(fakeDB)
	sMap := service.SessionMap
	sMap["fake"] = "fakeToken"

	result := serviceInstance.LogoutUser("fake")

	assert.Equal(t, true, result)

}

func TestEmployeeService_LogoutUser(t *testing.T) {
	var tests = []struct {
		userInMap bool
	}{
		{true},
		{false},
	}

	for _, tt := range tests {
		fakeDB := &servicefakes.FakeDatabaseInterface{}
		serviceInstance := service.NewEmployeeService(fakeDB)
		sMap := service.SessionMap
		if tt.userInMap {
			sMap["fake"] = "fakeToken"
			result := serviceInstance.LogoutUser("fake")

			assert.Equal(t, tt.userInMap, result)
		} else {
			result := serviceInstance.LogoutUser("faker")

			assert.Equal(t, tt.userInMap, result)

		}

	}
}

func TestRefreshToken(t *testing.T) {
	fakeToken := utils.GenerateToken(primitive.NewObjectID())

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	serviceInstance := service.NewEmployeeService(fakeDB)
	var tests = []struct {
		token          string
		expectedResult string
		doExpectErr    bool
	}{
		{fakeToken, fakeToken, false},
		{"banana", "", true},
	}

	for _, tt := range tests {
		actualResult, actualErr := serviceInstance.RefreshToken(tt.token)

		assert.Equal(t, tt.doExpectErr, actualErr != nil)
		if !tt.doExpectErr {
			assert.Equal(t, tt.token, actualResult)
		}
	}
}

func TestAuthenticateUser(t *testing.T) {
	fakeUserID := primitive.NewObjectID()
	fakeToken := utils.GenerateToken(fakeUserID)
	fakeUser := model.UserPayload{Group: "user"}
	fakeAdmin := model.UserPayload{Group: "admin"}

	routes.PermissionList.Permissions = append(routes.PermissionList.Permissions, model.Permission{Uri: "/user/", Methods: []string{"GET"}, GetSameUser: true, Group: "user"})
	routes.PermissionList.Permissions = append(routes.PermissionList.Permissions, model.Permission{Uri: "/user/", Methods: []string{"GET", "POST", "PUT", "DELETE"}, GetSameUser: false, Group: "admin"})

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	serviceInstance := service.NewEmployeeService(fakeDB)

	var tests = []struct {
		url            string
		userPayload    model.UserPayload
		method         string
		token          string
		expectedResult bool
		doExpectErr    bool
	}{
		{"/user/" + fakeUserID.Hex(), fakeUser, "GET", fakeToken, true, false},
		{"/user/", fakeAdmin, "GET", fakeToken, true, false},
		{"/user/", fakeUser, "POST", fakeToken, false, true},
		{"/user/", fakeAdmin, "POST", fakeToken, true, false},
	}

	for _, tt := range tests {
		fakeDB.GetUserByIDReturns(tt.userPayload, nil)

		actualResult, actualErr := serviceInstance.AuthenticateUser(tt.url, tt.method, tt.token)

		assert.Equal(t, tt.expectedResult, actualResult)
		assert.Equal(t, tt.doExpectErr, actualErr != nil)
	}
}

func TestAuthenticateUserErrorGetUserByID(t *testing.T) {
	fakeUserID := primitive.NewObjectID()
	fakeToken := utils.GenerateToken(fakeUserID)

	routes.PermissionList.Permissions = append(routes.PermissionList.Permissions, model.Permission{Uri: "/user/", Methods: []string{"GET"}, GetSameUser: true, Group: "user"})
	routes.PermissionList.Permissions = append(routes.PermissionList.Permissions, model.Permission{Uri: "/user/", Methods: []string{"GET", "POST", "PUT", "DELETE"}, GetSameUser: false, Group: "admin"})

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	serviceInstance := service.NewEmployeeService(fakeDB)
	fakeError := errors.New("test error")

	fakeDB.GetUserByIDReturns(model.UserPayload{}, fakeError)

	actualResult, actualErr := serviceInstance.AuthenticateUser("/user/", "GET", fakeToken)

	assert.Equal(t, false, actualResult)
	assert.Equal(t, fakeError, actualErr)
}
