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
