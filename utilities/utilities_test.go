package utilities

import (
	"example-project/model"
	"github.com/stretchr/testify/assert"
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
