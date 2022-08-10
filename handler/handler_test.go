package handler_test

import (
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

func TestGetEmployeeHandler_Return_valid_status_code(t *testing.T) {
	responseRecoder := httptest.NewRecorder()

	fakeContest, _ := gin.CreateTestContext(responseRecoder)
	fakeContest.Params = append(fakeContest.Params, gin.Param{Key: "id", Value: "1"})

	fakeService := &handlerfakes.FakeServiceInterface{}
	fakeService.GetEmployeeByIdReturns(model.Employee{
		ID:        "1",
		FirstName: "Joe",
	})

	handlerInstance := handler.NewHandler(fakeService)
	handlerInstance.GetEmployeeHandler(fakeContest)

	assert.Equal(t, http.StatusOK, responseRecoder.Code)

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
