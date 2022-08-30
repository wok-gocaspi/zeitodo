package utilities

import (
	"example-project/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
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
