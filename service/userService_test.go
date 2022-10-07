package service_test

import (
	"crypto/sha256"
	"errors"
	"example-project/model"
	"example-project/routes"
	"example-project/service"
	"example-project/service/servicefakes"
	"example-project/utilities"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http/httptest"
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

func TestUpdateUser_Service(t *testing.T) {
	type UserByIDCalls struct {
		callIteration int
		userReturn    model.UserPayload
		err           error
	}
	fakeUserID := primitive.NewObjectID()
	fakeNexUserID := primitive.NewObjectID()

	var tests = []struct {
		users             []model.UpdateUserPayload
		group             string
		userid            string
		firstUserSuccess  bool
		secondUserSuccess bool
		isStatusError     bool
		isDBError         bool
		DBError           error
		DBReturn          *mongo.UpdateResult
		UserByIDCalls     []UserByIDCalls
	}{
		{[]model.UpdateUserPayload{{ID: fakeUserID, Username: "jochen1"}, {ID: fakeNexUserID, Username: "peter2"}}, "admin", fakeUserID.Hex(), true, false, true, false, nil, &mongo.UpdateResult{},
			[]UserByIDCalls{{callIteration: 0, userReturn: model.UserPayload{Username: "jochen1"}, err: nil},
				{callIteration: 1, userReturn: model.UserPayload{}, err: errors.New("some error")},
			},
		},
		{[]model.UpdateUserPayload{{ID: fakeUserID, Username: "jochen1", FirstName: "jochen", LastName: "schweizer", Email: "j.schweizer@gmail.com", Team: "unkreativername", Projects: []string{"team1", "team2"}, TotalWorkingDays: 10, VacationDays: 20},
			{ID: fakeNexUserID, Username: "peter1", FirstName: "Peter", LastName: "Vogel", Email: "p.vogel@outlook.com", Team: "unkreativername", Projects: []string{"team1, team2"}, Group: "admin", Password: "1234567"}},
			"admin", fakeUserID.Hex(), true, true, true, false, nil, &mongo.UpdateResult{}, []UserByIDCalls{
				{callIteration: 0, userReturn: model.UserPayload{Username: "jochen1"}, err: nil},
				{callIteration: 1, userReturn: model.UserPayload{Username: "peter1"}, err: nil},
			},
		},
		{[]model.UpdateUserPayload{{ID: fakeUserID},
			{ID: fakeNexUserID, Username: "peter1", FirstName: "Peter", LastName: "Vogel"}},
			"admin", fakeUserID.Hex(), false, true, true, false, nil, &mongo.UpdateResult{}, []UserByIDCalls{
				{callIteration: 0, userReturn: model.UserPayload{Username: "jochen1"}, err: nil},
				{callIteration: 1, userReturn: model.UserPayload{Username: "peter1"}, err: nil},
			},
		},
		{[]model.UpdateUserPayload{{ID: fakeUserID, Username: "jochen1", FirstName: "jochen", LastName: "schweizer", Email: "j.schweizer@gmail.com", Team: "unkreativername", Projects: []string{"team1", "team2"}, TotalWorkingDays: 10, VacationDays: 20},
			{ID: fakeNexUserID, Username: "peter1", FirstName: "Peter", LastName: "Vogel", Email: "p.vogel@outlook.com", Team: "unkreativername", Projects: []string{"team1, team2"}, Password: "1234567"}},
			"user", fakeUserID.Hex(), true, false, true, false, nil, &mongo.UpdateResult{}, []UserByIDCalls{
				{callIteration: 0, userReturn: model.UserPayload{Username: "jochen1"}, err: nil},
				{callIteration: 1, userReturn: model.UserPayload{Username: "peter1"}, err: nil},
			},
		},
		{[]model.UpdateUserPayload{{ID: fakeUserID, Username: "jochen1", FirstName: "jochen", LastName: "schweizer", Email: "j.schweizer@gmail.com", Team: "unkreativername", Projects: []string{"team1", "team2"}, TotalWorkingDays: 10, VacationDays: 20},
			{ID: fakeNexUserID, Username: "peter1", FirstName: "Peter", LastName: "Vogel", Email: "p.vogel@outlook.com", Team: "unkreativername", Projects: []string{"team1, team2"}, Group: "admin", Password: "1234567"}},
			"admin", fakeUserID.Hex(), true, true, false, true, errors.New("some db error"), &mongo.UpdateResult{}, []UserByIDCalls{
				{callIteration: 0, userReturn: model.UserPayload{Username: "jochen1"}, err: nil},
				{callIteration: 1, userReturn: model.UserPayload{Username: "peter1"}, err: nil},
			},
		},
	}
	for _, tt := range tests {
		fakeDB := &servicefakes.FakeDatabaseInterface{}

		fakeDB.UpdateUserByIDReturns(tt.DBReturn, tt.DBError)
		serviceInstance := service.NewEmployeeService(fakeDB)
		for _, call := range tt.UserByIDCalls {
			fakeDB.GetUserByIDReturnsOnCall(call.callIteration, call.userReturn, call.err)
		}
		actual, err := serviceInstance.UpdateUsers(tt.users, tt.userid, tt.group)
		actualObj := actual.([]model.UserUpdateResult)
		if tt.isStatusError {
			assert.Equal(t, actualObj[0].Success, tt.firstUserSuccess)
			assert.Equal(t, actualObj[1].Success, tt.secondUserSuccess)
		}
		if tt.isDBError {
			assert.NotNil(t, err)
		}
		fakeDB = nil

	}
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
	fakeToken := utilities.GenerateToken(primitive.NewObjectID())

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
	fakeToken := utilities.GenerateToken(fakeUserID)
	fakeUser := model.UserPayload{Group: "user"}
	fakeAdmin := model.UserPayload{Group: "admin"}

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	serviceInstance := service.NewEmployeeService(fakeDB)

	var tests = []struct {
		userPayload model.UserPayload
		token       string
		doExpectErr bool
		userid      string
		usergroup   string
		dbErr       error
		isdberr     bool
	}{
		{fakeUser, fakeToken, false, fakeUserID.Hex(), "user", nil, false},
		{fakeAdmin, fakeToken, true, fakeUserID.Hex(), "user", errors.New("test"), false},
	}

	for _, tt := range tests {

		fakeDB.GetUserByIDReturns(tt.userPayload, tt.dbErr)

		_, _, actualErr := serviceInstance.AuthenticateUser(tt.token)

		assert.Equal(t, tt.doExpectErr, actualErr != nil)
	}
}

func TestAuthenticateUserErrorGetUserByID(t *testing.T) {
	fakeUserID := primitive.NewObjectID()
	fakeToken := utilities.GenerateToken(fakeUserID)

	routes.PermissionList.Permissions = append(routes.PermissionList.Permissions, model.Permission{Uri: "/user/", Methods: []string{"GET"}, GetSameUser: true, Group: "user"})
	routes.PermissionList.Permissions = append(routes.PermissionList.Permissions, model.Permission{Uri: "/user/", Methods: []string{"GET", "POST", "PUT", "DELETE"}, GetSameUser: false, Group: "admin"})

	fakeDB := &servicefakes.FakeDatabaseInterface{}
	serviceInstance := service.NewEmployeeService(fakeDB)
	fakeError := errors.New("test error")

	fakeDB.GetUserByIDReturns(model.UserPayload{}, fakeError)

	_, _, actualErr := serviceInstance.AuthenticateUser(fakeToken)

	assert.Equal(t, fakeError, actualErr)
}
func TestCheckUserPolicy(t *testing.T) {
	var permissionList model.PermissionList

	permissionList.Permissions = append(permissionList.Permissions, model.Permission{Uri: "/user/", Methods: []string{"GET"}, GetSameUser: true, Group: "user"})
	baseurl := "/user"
	userid := primitive.NewObjectID().Hex()
	adminid := primitive.NewObjectID().Hex()

	var tests = []struct {
		url    string
		urlid  string
		userid string
		group  string
		isErr  bool
		err    error
	}{
		{baseurl, userid, userid, "user", false, nil},
		{baseurl, adminid, userid, "user", false, errors.New("requesting user data of other users is not allowed")},
	}

	for _, tt := range tests {
		fakeDB := &servicefakes.FakeDatabaseInterface{}
		serviceInstance := service.NewEmployeeService(fakeDB)
		fakecontext := gin.Context{}

		fakecontext.Set("userid", tt.userid)
		fakecontext.Set("group", tt.group)
		fakecontext.Request = httptest.NewRequest("GET", tt.url+"/"+tt.urlid, nil)
		fakecontext.Params = append(fakecontext.Params, gin.Param{Key: "id", Value: tt.urlid})

		err := serviceInstance.CheckUserPolicy(&fakecontext, permissionList)

		assert.Equal(t, err, tt.err)

	}
}

func TestIsSameUser(t *testing.T) {
	var permissionList model.PermissionList
	permissionList.Permissions = append(permissionList.Permissions, model.Permission{Uri: "/user/", Methods: []string{"GET"}, GetSameUser: true, Group: "user"})
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	serviceInstance := service.NewEmployeeService(fakeDB)
	fakecontext := gin.Context{}
	userid := primitive.NewObjectID().Hex()
	fakecontext.Set("userid", userid)

	result := serviceInstance.CheckIsSameUser(&fakecontext, permissionList, userid)

	assert.Equal(t, nil, result)

}
func TestIsSameUserr(t *testing.T) {
	var permissionList model.PermissionList
	permissionList.Permissions = append(permissionList.Permissions, model.Permission{Uri: "/user/", Methods: []string{"GET"}, GetSameUser: true, Group: "user"})
	fakeDB := &servicefakes.FakeDatabaseInterface{}
	serviceInstance := service.NewEmployeeService(fakeDB)
	fakecontext := gin.Context{}
	userid := primitive.NewObjectID().Hex()

	result := serviceInstance.CheckIsSameUser(&fakecontext, permissionList, userid)

	assert.Equal(t, errors.New("requesting user data of other users is not allowed"), result)

}
