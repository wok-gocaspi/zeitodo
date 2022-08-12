package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"example-project/handler"
	"example-project/handler/handlerfakes"
	"example-project/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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

func TestCreateUser_Return_valid_200_multi(t *testing.T) {
	var fakeUserSignupPayload = []model.UserSignupPayload{
		{FirstName: "Peter", LastName: "Test", Email: "peter@test.com", Username: "ptest", Password: "123"},
		{FirstName: "Peter", LastName: "Test", Email: "peter@test.com", Username: "ptest", Password: "123"},
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

func TestCreateUser_Return_invalid_500_multi_insufficent_data(t *testing.T) {
	var fakeUserSignupPayload = []model.UserSignupPayload{
		{FirstName: "Peter", LastName: "Test", Email: "peter@test.com", Password: "123"},
		{FirstName: "Peter", LastName: "Test", Email: "peter@test.com", Username: "ptest", Password: "123"},
	}

	fakeUserSignupPayloadString, _ := json.Marshal(fakeUserSignupPayload)
	fmt.Println(string(fakeUserSignupPayloadString))
	responseRecorder := httptest.NewRecorder()

	fakeContext, _ := gin.CreateTestContext(responseRecorder)
	body := bytes.NewBufferString(string(fakeUserSignupPayloadString))
	fakeContext.Request = httptest.NewRequest("POST", "http://localhost:9090/user/create", body)

	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.CreateUserReturns([]model.UserSignupPayload{fakeUserSignupPayload[0]}, errors.New("insufficent user data"))
	responseRecorder.Body = body

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.CreateUserHandler(fakeContext)

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
}

func TestCreateUser_Return_invalid_500_multi_database_error(t *testing.T) {
	var fakeUserSignupPayload = []model.UserSignupPayload{
		{FirstName: "Peter", LastName: "Test", Email: "peter@test.com", Username: "ptest", Password: "123"},
		{FirstName: "Peter", LastName: "Test", Email: "peter@test.com", Username: "ptest", Password: "123"},
	}

	fakeUserSignupPayloadString, _ := json.Marshal(fakeUserSignupPayload)
	fmt.Println(string(fakeUserSignupPayloadString))
	responseRecorder := httptest.NewRecorder()

	fakeContext, _ := gin.CreateTestContext(responseRecorder)
	body := bytes.NewBufferString(string(fakeUserSignupPayloadString))
	fakeContext.Request = httptest.NewRequest("POST", "http://localhost:9090/user/create", body)

	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.CreateUserReturns(nil, errors.New("some db error"))
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

func TestHandler_GetProposalsById(t *testing.T) {

	filterReturn := []model.Proposal{
		model.Proposal{UserId: "1", Approved: false},
		model.Proposal{UserId: "2", Approved: true},
	}

	filterEmptyReturn := []model.Proposal{}
	fakeError := errors.New("fake error triggered")
	var tests = []struct {
		noQueryParams bool
		serviceErr    bool
		expectedCode  int
		Return        []model.Proposal
		err           error
	}{
		{true, false, 404, filterEmptyReturn, fakeError},
		{false, true, 404, filterReturn, fakeError},
		{false, false, 200, filterReturn, nil},
	}

	for _, tt := range tests {
		fakeRecorder := httptest.NewRecorder()
		fakeContext, _ := gin.CreateTestContext(fakeRecorder)
		fakeContext.Request = httptest.NewRequest("GET", "http://localhost:9090/employee/1/proposals", nil)

		fakeService := &handlerfakes.FakeServiceInterface{}
		fakeService.GetProposalsByIDReturns(tt.Return, tt.err)

		if tt.noQueryParams {
			fakeContext.Request = httptest.NewRequest("GET", "http://localhost:9090/xyz", nil)
			handlerInstance := handler.NewHandler(fakeService)
			handlerInstance.GetProposalsById(fakeContext)
			assert.Equal(t, tt.expectedCode, fakeRecorder.Code)
		}

		if tt.serviceErr {
			handlerInstance := handler.NewHandler(fakeService)
			handlerInstance.GetProposalsById(fakeContext)
			assert.Equal(t, tt.expectedCode, fakeRecorder.Code)

		}
		fakeContext.Params = append(fakeContext.Params, gin.Param{Key: "id", Value: "1"})
		handlerInstance := handler.NewHandler(fakeService)
		handlerInstance.GetProposalsById(fakeContext)
		assert.Equal(t, tt.expectedCode, fakeRecorder.Code)
	}
}

func TestHandler_UpdateProposalsHandler(t *testing.T) {
	jsonPayload := ` {
         "userId": "62eb8ba621c88b9be608b757",
        "startDate": "2013-Nov-20",
        "endDate": "2013-Nov-23",
        "type": "sickness"
        }
    `

	badPayload := ` 
        "startDate": "2013-Nov-20",
        "type": "sickness"
        }
    `
	noQueryError := "No date was given in the query parameter!"
	badPayloadError := "invalid payload"
	returnError := "invalid payload"

	var tests = []struct {
		Payload      string
		noDate       bool
		badPayload   bool
		returnError  bool
		expectedErr  string
		expectedCode int
	}{
		{jsonPayload, true, false, false, noQueryError, 404},
		{badPayload, false, true, false, badPayloadError, 400},
		{jsonPayload, false, false, true, returnError, 400},
		{jsonPayload, false, false, false, "", 200},
	}

	for _, tt := range tests {
		responseRecoder := httptest.NewRecorder()
		var mockUpdate model.Proposal
		json.Unmarshal([]byte(tt.Payload), &mockUpdate)
		body := bytes.NewBufferString(tt.Payload)

		fakeContest, _ := gin.CreateTestContext(responseRecoder)
		fakeService := &handlerfakes.FakeServiceInterface{}
		handlerInstance := handler.NewHandler(fakeService)
		if tt.noDate {
			fakeContest.Request = httptest.NewRequest("POST", "http://localhost:9090/employee/proposals/patch", body)
			handlerInstance.UpdateProposalsHandler(fakeContest)

			assert.Contains(t, responseRecoder.Body.String(), tt.expectedErr)
			assert.Equal(t, responseRecoder.Code, tt.expectedCode)
		}

		if tt.badPayload {
			fakeContest.Request = httptest.NewRequest("POST", "http://localhost:9090/employee/proposals/patch?date=2022-Nov-01", body)
			handlerInstance.UpdateProposalsHandler(fakeContest)

			assert.Contains(t, responseRecoder.Body.String(), tt.expectedErr)
			assert.Equal(t, responseRecoder.Code, tt.expectedCode)
		}

		if tt.returnError {
			fakeService.UpdateProposalByDateReturns(&mongo.UpdateResult{}, errors.New(tt.expectedErr))

			fakeContest.Request = httptest.NewRequest("POST", "http://localhost:9090/employee/proposals/patch?date=2022-Nov-01", body)
			handlerInstance.UpdateProposalsHandler(fakeContest)

			assert.Contains(t, responseRecoder.Body.String(), tt.expectedErr)
			assert.Equal(t, responseRecoder.Code, tt.expectedCode)
		}

		if !(tt.noDate || tt.returnError || tt.badPayload) {
			fakeService.UpdateProposalByDateReturns(&mongo.UpdateResult{}, nil)

			fakeContest.Request = httptest.NewRequest("POST", "http://localhost:9090/employee/proposals/patch?date=2022-Nov-01", body)
			handlerInstance.UpdateProposalsHandler(fakeContest)

			assert.Equal(t, responseRecoder.Code, tt.expectedCode)
		}
	}
}

func TestHandler_CreateProposalsHandler(t *testing.T) {
	jsonPayload := ` {
         "userId": "62eb8ba621c88b9be608b757",
        "startDate": "2013-Nov-20",
        "endDate": "2013-Nov-23",
        "type": "sickness"
        }
    `

	jsonPayloadRight := `    [
        {
         "userId": "62eb8ba621c88b9be608b757",
        "startDate": "2013-Nov-10",
        "endDate": "2013-Nov-11",
        "type": "vacation"
        }
    ]`
	badPayload := ` 
        "startDate": "2013-Nov-20",
        "endDate": "2013-Nov-23",
        "type": "sickness"
        }
    `

	badParamsMsg := "id is not given"
	badPayloadMsg := "invalid payload"

	var tests = []struct {
		Payload           string
		badParams         bool
		badPayload        bool
		createProposalErr bool
		expectedError     string
		expectedCode      int
	}{
		{jsonPayload, true, false, false, badParamsMsg, 404},
		{badPayload, false, true, false, badPayloadMsg, 400},
		{jsonPayloadRight, false, false, true, badPayloadMsg, 400},
		{jsonPayloadRight, false, false, false, badPayloadMsg, 200},
	}

	for _, tt := range tests {
		responseRecoder := httptest.NewRecorder()
		var mockUpdate model.Proposal
		json.Unmarshal([]byte(tt.Payload), &mockUpdate)
		body := bytes.NewBufferString(tt.Payload)

		fakeContest, _ := gin.CreateTestContext(responseRecoder)
		fakeService := &handlerfakes.FakeServiceInterface{}
		handlerInstance := handler.NewHandler(fakeService)
		if tt.badParams {
			fakeContest.Request = httptest.NewRequest("POST", "http://localhost:9090/employee/1/proposals/create", body)
			handlerInstance.CreateProposalsHandler(fakeContest)

			assert.Contains(t, responseRecoder.Body.String(), tt.expectedError)
			assert.Equal(t, responseRecoder.Code, tt.expectedCode)
		}

		if tt.badPayload {
			fakeContest.Params = append(fakeContest.Params, gin.Param{Key: "id", Value: "1"})
			fakeContest.Request = httptest.NewRequest("POST", "http://localhost:9090/employee/1/proposals/create", body)
			handlerInstance.CreateProposalsHandler(fakeContest)

			assert.Contains(t, responseRecoder.Body.String(), tt.expectedError)
			assert.Equal(t, responseRecoder.Code, tt.expectedCode)
		}

		if tt.createProposalErr {

			fakeService.CreateProposalsReturns(mongo.UpdateResult{}, errors.New(tt.expectedError))
			fakeContest.Params = append(fakeContest.Params, gin.Param{Key: "id", Value: "1"})
			fakeContest.Request = httptest.NewRequest("POST", "http://localhost:9090/employee/1/proposals/create", body)
			handlerInstance.CreateProposalsHandler(fakeContest)

			assert.Contains(t, responseRecoder.Body.String(), tt.expectedError)
			assert.Equal(t, responseRecoder.Code, tt.expectedCode)
		}

		if !tt.createProposalErr && !tt.badParams && !tt.badPayload {
			fakeService.CreateProposalsReturns(mongo.UpdateResult{}, nil)
			fakeContest.Params = append(fakeContest.Params, gin.Param{Key: "id", Value: "1"})
			fakeContest.Request = httptest.NewRequest("POST", "http://localhost:9090/employee/1/proposals/create", body)
			handlerInstance.CreateProposalsHandler(fakeContest)

			assert.Equal(t, responseRecoder.Code, tt.expectedCode)
		}

	}
}

func TestHandler_DeleteProposalHandler(t *testing.T) {
	fakeServiceErr := errors.New("error in Service")

	var tests = []struct {
		hasID         bool
		hasDate       bool
		hasServiceErr bool
		date          string
		expectedCode  int
	}{
		{false, false, false, "", 400},
		{true, false, false, "", 404},
		{true, true, true, "banana", 404},
		{true, true, false, "banana", 200},
	}

	for _, tt := range tests {
		fakeRecorder := httptest.NewRecorder()
		fakeContext, _ := gin.CreateTestContext(fakeRecorder)

		url := "http://localhost:9090/employee/1/proposals/delete"

		fakeContext.Request = httptest.NewRequest("DELETE", url, nil)
		if tt.hasID {
			fakeContext.Params = append(fakeContext.Params, gin.Param{Key: "id", Value: "1"})
		}

		if tt.hasDate {
			fakeContext.Request.URL.RawQuery = "date=" + tt.date
		}

		fakeService := &handlerfakes.FakeServiceInterface{}
		if tt.hasServiceErr {
			fakeService.DeleteProposalsByIDReturns(fakeServiceErr)
		} else {
			fakeService.DeleteProposalsByIDReturns(nil)
		}

		handlerInstance := handler.NewHandler(fakeService)
		handlerInstance.DeleteProposalHandler(fakeContext)

		if !tt.hasID {
			assert.Equal(t, tt.expectedCode, fakeRecorder.Code)
			assert.Contains(t, fakeRecorder.Body.String(), "id is not given")
		}
		if tt.hasID && !tt.hasDate {
			assert.Equal(t, tt.expectedCode, fakeRecorder.Code)
			assert.Contains(t, fakeRecorder.Body.String(), "No date was given in the query parameter!")
		}

		if tt.hasID && tt.hasDate && tt.hasServiceErr {
			assert.Equal(t, tt.expectedCode, fakeRecorder.Code)
			assert.Contains(t, fakeRecorder.Body.String(), "error in Service")
		}

		if tt.hasID && tt.hasDate && !tt.hasServiceErr {
			assert.Equal(t, tt.expectedCode, fakeRecorder.Code)
			assert.Equal(t, fakeRecorder.Body.String(), "\"\"")
		}
	}

}
