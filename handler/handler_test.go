package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"example-project/handler"
	"example-project/handler/handlerfakes"
	"example-project/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetUserHandler_Return_valid_200(t *testing.T) {
	responseRecoder := httptest.NewRecorder()

	fakeContest, _ := gin.CreateTestContext(responseRecoder)
	fakeContest.Params = append(fakeContest.Params, gin.Param{Key: "id", Value: "1"})

	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.GetUserByIDReturns(model.UserPayload{
		FirstName: "Joe",
		LastName:  "Hitch",
	}, nil)

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.GetUserHandler(fakeContest)

	assert.Equal(t, http.StatusOK, responseRecoder.Code)

}

func TestGetUserHandler_Return_invalid_400(t *testing.T) {
	responseRecoder := httptest.NewRecorder()

	fakeContest, _ := gin.CreateTestContext(responseRecoder)

	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.GetUserByIDReturns(model.UserPayload{
		FirstName: "Joe",
		LastName:  "Hitch",
	}, nil)

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.GetUserHandler(fakeContest)

	assert.Equal(t, http.StatusBadRequest, responseRecoder.Code)
}

func TestGetUserHandler_Return_invalid_404(t *testing.T) {
	responseRecoder := httptest.NewRecorder()

	fakeContest, _ := gin.CreateTestContext(responseRecoder)
	fakeContest.Params = append(fakeContest.Params, gin.Param{Key: "id", Value: "1"})

	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.GetUserByIDReturns(model.UserPayload{}, errors.New("some db error"))

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.GetUserHandler(fakeContest)

	assert.Equal(t, http.StatusNotFound, responseRecoder.Code)

}

func TestGetAllUserHandler_Return_valid_200(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	fakeContext, _ := gin.CreateTestContext(responseRecorder)
	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.GetAllUserReturns([]model.UserPayload{}, nil)

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.GetAllUserHandler(fakeContext)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestGetAllUserHandler_Return_invalid_404(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	fakeContext, _ := gin.CreateTestContext(responseRecorder)
	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.GetAllUserReturns(nil, errors.New("some error"))

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.GetAllUserHandler(fakeContext)

	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
}

func TestGetTeamMemberHandler_Return_valid_200_id_given(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	fakeContext, _ := gin.CreateTestContext(responseRecorder)
	fakeContext.Request = httptest.NewRequest("GET", "http://localhost:9090/user/team/get?id=1", nil)
	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.GetTeamMembersByUserIDReturns([]model.TeamMember{}, nil)

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.GetTeamMemberHandler(fakeContext)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestGetTeamMemberHandler_Return_valid_200_team_given(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	fakeContext, _ := gin.CreateTestContext(responseRecorder)
	fakeContext.Request = httptest.NewRequest("GET", "http://localhost:9090/user/team/get?name=powerrangers", nil)
	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.GetTeamMembersByNameReturns([]model.TeamMember{}, nil)

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.GetTeamMemberHandler(fakeContext)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestGetTeamMemberHandler_Return_valid_400_insufficent_queries(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	fakeContext, _ := gin.CreateTestContext(responseRecorder)
	fakeContext.Request = httptest.NewRequest("GET", "http://localhost:9090/user/team/get", nil)
	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.GetTeamMembersByUserIDReturns(nil, errors.New("no queries errors"))

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.GetTeamMemberHandler(fakeContext)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func TestGetTeamMemberHandler_Return_invalid_500_id_error(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	fakeContext, _ := gin.CreateTestContext(responseRecorder)
	fakeContext.Request = httptest.NewRequest("GET", "http://localhost:9090/user/team/get?id=1", nil)
	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.GetTeamMembersByUserIDReturns(nil, errors.New("some error"))

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.GetTeamMemberHandler(fakeContext)

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
}

func TestGetTeamMemberHandler_Return_invalid_500_team_error(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	fakeContext, _ := gin.CreateTestContext(responseRecorder)
	fakeContext.Request = httptest.NewRequest("GET", "http://localhost:9090/user/team/get?name=1", nil)
	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.GetTeamMembersByNameReturns(nil, errors.New("some error"))

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.GetTeamMemberHandler(fakeContext)

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
}

func TestCreateUser_Return_valid_200_single(t *testing.T) {
	var fakeUserSignupPayload = model.UserSignupPayload{
		FirstName: "Peter", LastName: "Test", Email: "peter@test.com", Username: "ptest", Password: "123",
	}

	fakeUserSignupPayloadString, _ := json.Marshal(fakeUserSignupPayload)
	fmt.Println(string(fakeUserSignupPayloadString))
	responseRecorder := httptest.NewRecorder()

	fakeContext, _ := gin.CreateTestContext(responseRecorder)
	body := bytes.NewBufferString(string(fakeUserSignupPayloadString))
	fakeContext.Request = httptest.NewRequest("POST", "http://localhost:9090/user/create", body)

	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.CreateUserReturns(mongo.InsertManyResult{}.InsertedIDs, nil)

	responseRecorder.Body = body

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.CreateUserHandler(fakeContext)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestCreateUser_Return_invalid_500_single_insufficent_data(t *testing.T) {
	var fakeUserSignupPayload = model.UserSignupPayload{
		FirstName: "Peter", LastName: "Test", Email: "peter@test.com", Password: "123",
	}

	fakeUserSignupPayloadString, _ := json.Marshal(fakeUserSignupPayload)
	fmt.Println(string(fakeUserSignupPayloadString))
	responseRecorder := httptest.NewRecorder()

	fakeContext, _ := gin.CreateTestContext(responseRecorder)
	body := bytes.NewBufferString(string(fakeUserSignupPayloadString))
	fakeContext.Request = httptest.NewRequest("POST", "http://localhost:9090/user/create", body)

	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.CreateUserReturns(nil, errors.New("insufficent user data"))
	responseRecorder.Body = body

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.CreateUserHandler(fakeContext)

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
}

func TestCreateUser_Return_invalid_500_invalid_json(t *testing.T) {

	var fakeJSONString = `
		[
			{
				"username": "Peter12355",
				"password": "234",
				"lastname": "Müller",
				"email": "p.mueller@gmx.com"
			,
			{
				"username": "Peter12355",
				"password": "234",
				"lastname": "Müller",
				"email": "p.mueller@gmx.com"
			}
		]
	`

	responseRecorder := httptest.NewRecorder()

	fakeContext, _ := gin.CreateTestContext(responseRecorder)
	body := bytes.NewBufferString(fakeJSONString)
	fakeContext.Request = httptest.NewRequest("POST", "http://localhost:9090/user/create", body)

	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.CreateUserReturns(nil, errors.New("insufficent user data"))
	responseRecorder.Body = body

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.CreateUserHandler(fakeContext)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func TestUpdateUserHandler_Return_valid_200_single(t *testing.T) {
	var fakeUserSignupPayload = model.User{
		FirstName: "Peter", ID: primitive.ObjectID{},
	}

	fakeUserSignupPayloadString, _ := json.Marshal(fakeUserSignupPayload)
	responseRecorder := httptest.NewRecorder()

	fakeContext, _ := gin.CreateTestContext(responseRecorder)
	body := bytes.NewBufferString(string(fakeUserSignupPayloadString))
	fakeContext.Request = httptest.NewRequest("POST", "http://localhost:9090/user/update", body)

	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.UpdateUsersReturns([]model.UserUpdateResult{}, nil)
	responseRecorder.Body = body

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.UpdateUserHandler(fakeContext)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestUpdateUserHandler_Return_invalid_json(t *testing.T) {
	var fakeJSONString = `
		[
			{
				"id": "123",
				"email": "p.mueller@gmx.com"
			,
			{
				"id": "456",
				"email": "p.mueller2@gmx.com"
			}
		]
	`
	responseRecorder := httptest.NewRecorder()

	fakeContext, _ := gin.CreateTestContext(responseRecorder)
	body := bytes.NewBufferString(fakeJSONString)
	fakeContext.Request = httptest.NewRequest("POST", "http://localhost:9090/user/update", body)

	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.UpdateUsersReturns([]model.UserUpdateResult{}, nil)
	responseRecorder.Body = body

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.UpdateUserHandler(fakeContext)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func TestUpdateUserHandler_Return_invalid_500_multi_update_unsuccessful(t *testing.T) {
	var fakeUserSignupPayload = []model.User{
		{FirstName: "Peter", ID: primitive.NewObjectID()},
		{FirstName: "Peter", ID: primitive.NewObjectID()},
	}

	fakeUserSignupPayloadString, _ := json.Marshal(fakeUserSignupPayload)
	responseRecorder := httptest.NewRecorder()

	fakeContext, _ := gin.CreateTestContext(responseRecorder)
	body := bytes.NewBufferString(string(fakeUserSignupPayloadString))
	fakeContext.Request = httptest.NewRequest("POST", "http://localhost:9090/user/update", body)

	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.UpdateUsersReturns([]model.UserUpdateResult{}, errors.New("a few users couldn't be updated"))
	responseRecorder.Body = body

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.UpdateUserHandler(fakeContext)

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
}

func TestUpdateUserHandler_Return_valid_200_multi(t *testing.T) {
	var fakeUserSignupPayload = []model.User{
		{FirstName: "Peter", ID: primitive.ObjectID{}},
		{FirstName: "Peter", ID: primitive.ObjectID{}},
	}

	fakeUserSignupPayloadString, _ := json.Marshal(fakeUserSignupPayload)
	responseRecorder := httptest.NewRecorder()

	fakeContext, _ := gin.CreateTestContext(responseRecorder)
	body := bytes.NewBufferString(string(fakeUserSignupPayloadString))
	fakeContext.Request = httptest.NewRequest("POST", "http://localhost:9090/user/update", body)

	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.UpdateUsersReturns([]model.UserUpdateResult{}, nil)
	responseRecorder.Body = body

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.UpdateUserHandler(fakeContext)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestUpdateUserHandler_Return_invalid_500_single_update_unsuccessful(t *testing.T) {
	var fakeUserSignupPayload = model.User{
		FirstName: "Peter", ID: primitive.NewObjectID(),
	}

	fakeUserSignupPayloadString, _ := json.Marshal(fakeUserSignupPayload)
	responseRecorder := httptest.NewRecorder()

	fakeContext, _ := gin.CreateTestContext(responseRecorder)
	body := bytes.NewBufferString(string(fakeUserSignupPayloadString))
	fakeContext.Request = httptest.NewRequest("POST", "http://localhost:9090/user/update", body)

	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.UpdateUsersReturns([]model.UserUpdateResult{}, errors.New("a few users couldn't be updated"))
	responseRecorder.Body = body

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.UpdateUserHandler(fakeContext)

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
}

func TestDeleteUserHandler_Return_valid_200_delete(t *testing.T) {

	responseRecorder := httptest.NewRecorder()

	fakeContext, _ := gin.CreateTestContext(responseRecorder)
	fakeContext.Params = append(fakeContext.Params, gin.Param{Key: "id", Value: "1"})

	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.DeleteUsersReturns(mongo.DeleteResult{DeletedCount: 1}, nil)

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.DeleteUserHandler(fakeContext)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestDeleteUserHandler_Return_invalid_400_id(t *testing.T) {

	responseRecorder := httptest.NewRecorder()

	fakeContext, _ := gin.CreateTestContext(responseRecorder)
	fakeContext.Params = append(fakeContext.Params, gin.Param{Key: "i", Value: "1"})

	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.DeleteUsersReturns(mongo.DeleteResult{DeletedCount: 1}, nil)

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.DeleteUserHandler(fakeContext)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func TestDeleteUserHandler_Return_invalid_500_no_deletion(t *testing.T) {

	responseRecorder := httptest.NewRecorder()

	fakeContext, _ := gin.CreateTestContext(responseRecorder)
	fakeContext.Params = append(fakeContext.Params, gin.Param{Key: "id", Value: "1"})

	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.DeleteUsersReturns(mongo.DeleteResult{DeletedCount: 0}, errors.New("no user have been deleted, please check the id"))

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.DeleteUserHandler(fakeContext)

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
}

func TestLoginUserHandler(t *testing.T) {
	fakePayload := model.UserAuthPayload{
		Username: "testusername",
		Password: "testuserpwd",
	}
	fakePayloadString, _ := json.Marshal(fakePayload)

	expDate := time.Now().Add(time.Minute * 5)
	fakeCookie := http.Cookie{
		Name:     "token",
		Value:    "this is  sample token",
		Expires:  expDate,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	}

	fakeCookieHeader := "token=this+is++sample+token; Path=/; Domain=localhost; Max-Age=3600; HttpOnly"
	fakeServiceErr := errors.New("user is unauthorized")

	var tests = []struct {
		body                *bytes.Buffer
		serviceResponse     http.Cookie
		serviceErr          error
		expectedStatus      int
		expectedCookieCount int
		expectedCookie      string
	}{
		{bytes.NewBufferString(""), fakeCookie, nil, http.StatusBadRequest, 0, ""},
		{bytes.NewBufferString(string(fakePayloadString)), fakeCookie, nil, http.StatusOK, 1, fakeCookieHeader},
		{bytes.NewBufferString(string(fakePayloadString)), fakeCookie, fakeServiceErr, http.StatusUnauthorized, 0, ""},
	}

	for _, tt := range tests {
		responseRecorder := httptest.NewRecorder()

		fakeContext, _ := gin.CreateTestContext(responseRecorder)
		fakeContext.Request = httptest.NewRequest("POST", "http://localhost:9090/user/login", tt.body)

		fakeService := &handlerfakes.FakeServiceInterface{}
		fakeService.LoginUserReturns(tt.serviceResponse, tt.serviceErr)

		handlerInstance := handler.NewHandler(fakeService)
		handlerInstance.LoginUserHandler(fakeContext)

		assert.Equal(t, tt.expectedStatus, responseRecorder.Code)
		assert.Equal(t, tt.expectedCookieCount, len(responseRecorder.Header()["Set-Cookie"]))
		if tt.expectedCookieCount > 0 {
			assert.Equal(t, tt.expectedCookie, responseRecorder.Header()["Set-Cookie"][0])
		}
	}
}

func TestLogoutUserHandler(t *testing.T) {
	expDate := time.Now().Add(time.Minute * 5)
	fakeCookie := http.Cookie{
		Name:     "token",
		Value:    "this is  sample token",
		Expires:  expDate,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	}

	fakeCookieHeader := "token=; Path=/; HttpOnly; Secure"

	var tests = []struct {
		hasValidToken       bool
		reqCookie           http.Cookie
		expectedStatus      int
		expectedCookieCount int
		expectedCookie      string
	}{
		{true, fakeCookie, http.StatusOK, 1, fakeCookieHeader},
		{false, fakeCookie, http.StatusUnauthorized, 0, ""},
	}

	for _, tt := range tests {
		responseRecorder := httptest.NewRecorder()

		fakeContext, _ := gin.CreateTestContext(responseRecorder)
		fakeContext.Request = httptest.NewRequest("POST", "http://localhost:9090/user/logout", nil)
		if tt.hasValidToken {
			fakeContext.Request.AddCookie(&tt.reqCookie)
		}

		fakeService := &handlerfakes.FakeServiceInterface{}
		handlerInstance := handler.NewHandler(fakeService)
		handlerInstance.LogoutUserHandler(fakeContext)

		assert.Equal(t, tt.expectedStatus, responseRecorder.Code)
		assert.Equal(t, tt.expectedCookieCount, len(responseRecorder.Header()["Set-Cookie"]))
		if tt.expectedCookieCount > 0 {
			assert.Equal(t, tt.expectedCookie, responseRecorder.Header()["Set-Cookie"][0])
		}
	}
}

func TestRefreshTokenHandler(t *testing.T) {
	expDate := time.Now().Add(time.Minute * 5)
	fakeCookie := http.Cookie{
		Name:     "token",
		Value:    "this is  sample token",
		Expires:  expDate,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	}
	fakeServiceToken := "serviceToken"

	fakeCookieHeader := "token=" + fakeServiceToken + "; Path=/; Domain=localhost; Max-Age=3600; HttpOnly"
	fakeServiceErr := errors.New("fake unauthorized user")

	var tests = []struct {
		hasValidToken       bool
		reqCookie           http.Cookie
		serviceToken        string
		serviceErr          error
		expectedStatus      int
		expectedCookieCount int
		expectedCookie      string
	}{
		{false, fakeCookie, "", nil, http.StatusUnauthorized, 0, ""},
		{true, fakeCookie, "", fakeServiceErr, http.StatusUnauthorized, 0, ""},
		{true, fakeCookie, fakeServiceToken, nil, http.StatusOK, 1, fakeCookieHeader},
	}

	for _, tt := range tests {
		responseRecorder := httptest.NewRecorder()

		fakeContext, _ := gin.CreateTestContext(responseRecorder)
		fakeContext.Request = httptest.NewRequest("POST", "http://localhost:9090/user/refresh", nil)
		if tt.hasValidToken {
			fakeContext.Request.AddCookie(&tt.reqCookie)
		}

		fakeService := &handlerfakes.FakeServiceInterface{}
		fakeService.RefreshTokenReturns(tt.serviceToken, tt.serviceErr)
		handlerInstance := handler.NewHandler(fakeService)
		handlerInstance.RefreshTokenHandler(fakeContext)

		assert.Equal(t, tt.expectedStatus, responseRecorder.Code)
		assert.Equal(t, tt.expectedCookieCount, len(responseRecorder.Header()["Set-Cookie"]))
		if tt.expectedCookieCount > 0 {
			assert.Equal(t, tt.expectedCookie, responseRecorder.Header()["Set-Cookie"][0])
		}
	}
}

func TestPermissionMiddleware(t *testing.T) {
	expDate := time.Now().Add(time.Minute * 5)
	fakeCookie := http.Cookie{
		Name:     "token",
		Value:    "this is  sample token",
		Expires:  expDate,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	}

	fakeServiceErr := errors.New("fake unauthorized user")

	var tests = []struct {
		hasValidToken  bool
		reqCookie      http.Cookie
		serviceOk      bool
		serviceErr     error
		expectedStatus int
	}{
		{false, fakeCookie, false, nil, http.StatusUnauthorized},
		{true, fakeCookie, false, fakeServiceErr, http.StatusUnauthorized},
	}

	for _, tt := range tests {

		responseRecorder := httptest.NewRecorder()

		fakeContext, _ := gin.CreateTestContext(responseRecorder)
		fakeContext.Request = httptest.NewRequest("POST", "http://localhost:9090/user/refresh", nil)
		fakeService := &handlerfakes.FakeServiceInterface{}

		if tt.hasValidToken {
			fakeContext.Request.AddCookie(&tt.reqCookie)
			fakeService.AuthenticateUserReturns(tt.serviceOk, tt.serviceErr)
		}

		handlerInstance := handler.NewHandler(fakeService)
		handlerInstance.PermissionMiddleware(fakeContext)

		assert.Equal(t, tt.expectedStatus, responseRecorder.Code)
	}
}
