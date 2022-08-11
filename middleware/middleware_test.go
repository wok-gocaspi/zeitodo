package middleware

import (
	"example-project/service/servicefakes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetupEngine_RequestToUnknownPathReturnsNotFound(t *testing.T) {
	// arrange
	fakeDb := servicefakes.FakeDatabaseInterface{}
	fakeHandlerFunc1 := SetupService(&fakeDb)
	fakeHandlerFunc2 := SetupService(&fakeDb)
	var fakeHandlerArr []gin.HandlerFunc
	fakeHandlerArr = []gin.HandlerFunc{fakeHandlerFunc1, fakeHandlerFunc2}

	engine := SetupEngine(fakeHandlerArr)
	request := httptest.NewRequest("GET", "/blop", nil)
	responseRecorder := httptest.NewRecorder()

	// act
	engine.ServeHTTP(responseRecorder, request)

	// assert
	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
}
func TestSetupService(t *testing.T) {
	fakeDb := servicefakes.FakeDatabaseInterface{}
	actualHandler := SetupService(&fakeDb)
	assert.NotNil(t, actualHandler)
}

func TestSetupEngine(t *testing.T) {
	fakeDb := servicefakes.FakeDatabaseInterface{}
	fakeHandlerFunc1 := SetupService(&fakeDb)
	fakeHandlerFunc2 := SetupService(&fakeDb)
	var fakeHandlerArr []gin.HandlerFunc
	fakeHandlerArr = []gin.HandlerFunc{fakeHandlerFunc1, fakeHandlerFunc2}

	actualHandler := SetupEngine(fakeHandlerArr)
	assert.NotNil(t, actualHandler)

}
