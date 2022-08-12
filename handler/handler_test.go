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
	"go.mongodb.org/mongo-driver/mongo"
	"net/http/httptest"
	"testing"
)

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
