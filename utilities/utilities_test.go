package utilities

import (
	"errors"
	"example-project/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateTimeObject(t *testing.T) {

	mockStart := "2006-Nov-01"
	mockEnd := "2006-Nov-02"

	_, err := CreateTimeObject(mockStart, mockEnd)

	assert.Equal(t, nil, err)
	//	assert.Equal(t, mockStart, *actual.Interval.Start())

}

func TestCraftProposalFromPayload(t *testing.T) {
	okPayload := []model.ProposalPayload{
		model.ProposalPayload{UserId: "1", StartDate: "2006-Nov-01", EndDate: "2006-Nov-10"},
	}
	createObjectErrorMsg := "Error occured in building the time interval for a new proposal"

	var tests = []struct {
		Payload        []model.ProposalPayload
		createObjErr   bool
		expectedErrMsg string
	}{
		{okPayload, false, createObjectErrorMsg},
		{okPayload, false, createObjectErrorMsg},
	}

	for _, tt := range tests {
		if tt.createObjErr {
			_, err := CraftProposalFromPayload(tt.Payload)

			assert.Contains(t, err, tt.expectedErrMsg)
		} else {
			_, err := CraftProposalFromPayload(tt.Payload)
			assert.Equal(t, nil, err)
		}
	}
}

func TestCustomOverlaps(t *testing.T) {
	t1, _ := CreateTimeObject("2006-Nov-01", "2006-Nov-05")
	t2, _ := CreateTimeObject("2006-Nov-02", "2006-Nov-03")
	t2Overlapp, _ := CreateTimeObject("2006-Nov-01", "2006-Nov-10")

	proposal1 := model.Proposal{UserId: "1", StartDate: "2006-Nov-01", EndDate: "2006-Nov-05", TimeObject: t1}
	proposal2 := model.Proposal{UserId: "1", StartDate: "2006-Nov-06", EndDate: "2006-Nov-10", TimeObject: t2}
	proposal2Overlapp := model.Proposal{UserId: "1", StartDate: "2006-Nov-04", EndDate: "2006-Nov-10", TimeObject: t2Overlapp}

	var tests = []struct {
		p1          model.Proposal
		p2          model.Proposal
		hasOverlapp bool
	}{
		{proposal1, proposal2Overlapp, true},
		{proposal2Overlapp, proposal1, true},
		{proposal1, proposal2, false},
	}

	for _, tt := range tests {
		if tt.hasOverlapp {
			actual := CustomOverlaps(tt.p1, tt.p2)
			assert.Equal(t, actual, tt.hasOverlapp)
		} else {
			actual := CustomOverlaps(tt.p1, tt.p2)
			assert.Equal(t, actual, tt.hasOverlapp)
		}
	}
}

func TestStartDateExceedsEndDate(t *testing.T) {

	testArrayOk := []model.Proposal{
		model.Proposal{UserId: "1", StartDate: "2006-Nov-01", EndDate: "2006-Nov-05"},
		model.Proposal{UserId: "1", StartDate: "2006-Nov-02", EndDate: "2006-Nov-05"},
	}

	testArrayErr := []model.Proposal{
		model.Proposal{UserId: "1", StartDate: "2006-Nov-01", EndDate: "2006-Nov-05"},
		model.Proposal{UserId: "1", StartDate: "2006-Nov-12", EndDate: "2006-Nov-05"},
	}

	var tests = []struct {
		Arr    []model.Proposal
		hasErr bool
	}{
		{testArrayOk, false},
		{testArrayErr, true},
	}

	for _, tt := range tests {
		if !tt.hasErr {
			actual := StartDateExceedsEndDate(tt.Arr)
			assert.Equal(t, tt.hasErr, actual)
		} else {
			actual := StartDateExceedsEndDate(tt.Arr)
			assert.Equal(t, tt.hasErr, actual)
		}
	}
}

func TestProposalTimeIntersectsProposals(t *testing.T) {
	testArray := []model.Proposal{
		model.Proposal{UserId: "1", StartDate: "2006-Nov-01", EndDate: "2006-Nov-05"},
		model.Proposal{UserId: "1", StartDate: "2006-Nov-02", EndDate: "2006-Nov-05"},
	}
	testArray1 := []model.Proposal{
		model.Proposal{UserId: "1", StartDate: "2006-Nov-02", EndDate: "2006-Nov-05"},
	}
	tEquals, _ := CreateTimeObject("2006-Nov-01", "2006-Nov-05")
	proposalEquals := model.Proposal{UserId: "1", StartDate: "2006-Nov-01", EndDate: "2006-Nov-05", TimeObject: tEquals}

	tOverlapps, _ := CreateTimeObject("2006-Nov-01", "2006-Nov-04")
	proposalOverlapps := model.Proposal{UserId: "1", StartDate: "2006-Nov-01", EndDate: "2006-Nov-05", TimeObject: tOverlapps}

	tDuring, _ := CreateTimeObject("2006-Nov-02", "2006-Nov-04")
	proposalDuring := model.Proposal{UserId: "1", StartDate: "2006-Nov-02", EndDate: "2006-Nov-04", TimeObject: tDuring}

	tDuring1, _ := CreateTimeObject("2006-Nov-01", "2006-Nov-10")
	proposalDuring1 := model.Proposal{UserId: "1", StartDate: "2006-Nov-01", EndDate: "2006-Nov-10", TimeObject: tDuring1}

	tOk, _ := CreateTimeObject("2006-Nov-10", "2006-Nov-15")
	proposalOk := model.Proposal{UserId: "1", StartDate: "2006-Nov-10", EndDate: "2006-Nov-15", TimeObject: tOk}

	var tests = []struct {
		Arr                    []model.Proposal
		Proposal               model.Proposal
		proposalsEquals        bool
		proposalOverlapps      bool
		proposalDuring         bool
		proposalDuringRedirect bool
		noIntersect            bool
	}{
		{testArray, proposalEquals, true, false, false, false, false},
		{testArray, proposalOverlapps, false, true, false, false, false},
		{testArray, proposalDuring, false, false, true, false, false},
		{testArray1, proposalDuring1, false, false, false, true, false},
		{testArray1, proposalOk, false, false, false, false, true},
	}

	for _, tt := range tests {
		if tt.proposalsEquals {
			actual := ProposalTimeIntersectsProposals(tt.Proposal, tt.Arr)
			assert.Equal(t, actual, tt.proposalsEquals)
		}
		if tt.proposalOverlapps {
			actual := ProposalTimeIntersectsProposals(tt.Proposal, tt.Arr)
			assert.Equal(t, actual, tt.proposalOverlapps)
		}
		if tt.proposalDuring {
			actual := ProposalTimeIntersectsProposals(tt.Proposal, tt.Arr)
			assert.Equal(t, actual, tt.proposalDuring)
		}
		if tt.proposalDuringRedirect {
			actual := ProposalTimeIntersectsProposals(tt.Proposal, tt.Arr)
			assert.Equal(t, actual, tt.proposalDuringRedirect)
		}
		if tt.noIntersect {
			actual := ProposalTimeIntersectsProposals(tt.Proposal, tt.Arr)
			assert.Equal(t, actual, !tt.noIntersect)
		}

	}
}

func TestGenerateToken(t *testing.T) {
	fakeUserId := primitive.NewObjectID()

	resultString := GenerateToken(fakeUserId)
	assert.NotEmpty(t, resultString)
}

func TestValidateToken(t *testing.T) {
	fakeToken := "fakeToken"

	_, _, err := ValidateToken(fakeToken)

	assert.Error(t, err)
}

func TestUserUpdateSetter(t *testing.T) {
	var tests = []struct {
		user         model.UpdateUserPayload
		userGroup    string
		returnSetter bson.D
		isError      error
	}{
		{user: model.UpdateUserPayload{Username: "Thomas", Group: "admin", FirstName: "Jeff", LastName: "TheFirst", Email: "j.thefirst@gmail.com", Password: "mynameisjeff", TotalWorkingDays: 10, VacationDays: 5, Team: "okapi", Projects: []string{"okapi", "tetris"}}, userGroup: "admin", isError: nil, returnSetter: bson.D{primitive.E{Key: "$set", Value: primitive.D{primitive.E{Key: "firstname", Value: "Jeff"}, primitive.E{Key: "lastname", Value: "TheFirst"}, primitive.E{Key: "email", Value: "j.thefirst@gmail.com"}, primitive.E{Key: "team", Value: "okapi"}, primitive.E{Key: "totalWorkingDays", Value: float32(10)}, primitive.E{Key: "vacationDays", Value: 5}, primitive.E{Key: "username", Value: "Thomas"}, primitive.E{Key: "password", Value: [32]uint8{0x5c, 0x7c, 0x56, 0x22, 0xef, 0x18, 0x9d, 0x75, 0x4f, 0xf9, 0xcc, 0x6, 0xe0, 0x3e, 0xa9, 0xf9, 0x1f, 0xb6, 0x98, 0xe2, 0x7, 0xc3, 0x8, 0x67, 0x46, 0xd4, 0x92, 0x5, 0xa2, 0xd1, 0xfc, 0x0}}, primitive.E{Key: "group", Value: "admin"}, primitive.E{Key: "projects", Value: []string{"okapi", "tetris"}}}}}},
		{user: model.UpdateUserPayload{Group: "admin"}, userGroup: "user", isError: errors.New("no data changed on user"), returnSetter: nil},
	}

	for _, tt := range tests {
		result, err := UserUpdateSetter(tt.user, tt.userGroup)
		assert.Equal(t, tt.returnSetter, result)
		assert.Equal(t, err, tt.isError)
	}
}

func TestGetWeekdaysBetween(t *testing.T) {
	var tests = []struct {
		startDate time.Time
		endDate   time.Time
		result    int
	}{
		{
			startDate: time.Date(2022, 10, 3, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2022, 10, 12, 0, 0, 0, 0, time.UTC),
			result:    7,
		},
	}
	for _, tt := range tests {
		result := GetWeekdaysBetween(tt.startDate, tt.endDate)
		assert.Equal(t, result, tt.result)
	}
}

func TestFormGetAllProposalsFilter(t *testing.T) {
	fakeContext := gin.Context{}
	fakeContext.Request = httptest.NewRequest("GET", "http://localhost:9090/proposals?sort=asce&type=vacation&status=approved&username=jack", nil)
	filter, sort := FormGetAllProposalsFilter(primitive.NewObjectID().Hex(), &fakeContext)
	assert.Equal(t, filter, filter)
	assert.Equal(t, sort, sort)
}

func TestFormGetAllProposalsFilterDESC(t *testing.T) {
	fakeContext := gin.Context{}
	fakeContext.Request = httptest.NewRequest("GET", "http://localhost:9090/proposals?sort=desc&type=vacation&status=approved&username=jack", nil)
	filter, sort := FormGetAllProposalsFilter(primitive.NewObjectID().Hex(), &fakeContext)
	assert.Equal(t, filter, filter)
	assert.Equal(t, sort, sort)
}
