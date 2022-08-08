package service

import (
	"errors"
	"example-project/model"
	"fmt"
	"github.com/retailify/go-interval"
	"strings"
	"time"
)

func (s EmployeeService) GetProposalsByID(id string) ([]model.Proposal, error) {
	result, err := s.DbService.GetProposals(id)
	if err != nil {
		return []model.Proposal{}, err
	}
	if len(result) == 0 && err == nil {
		noResultsErr := errors.New("No results could be found to your query")
		return result, noResultsErr
	}

	return result, err
}

/*
func (s EmployeeService) CreateProposals(proposalPayloadArr []model.ProposalPayload, id string) (interface{}, error) {
	proposalArr, err := CraftProposalFromPayload(proposalPayloadArr)
	if err != nil {
		return nil, err
	}

	if StartDateExceedsEndDate(proposalArr) {
		startExceedsEndErrMsg := errors.New("The startdate must be before the enddate")
		return nil, startExceedsEndErrMsg
	}

	overlappingErrMsg := errors.New("There cant be overlapping proposals.")
	actualProposals, err := s.GetProposalsByID(id)

	var proposals []interface{}
	for _, p := range proposalArr {
		if !ProposalTimeIntersectsProposals(p, actualProposals) {
			proposals = append(proposals, p)
		} else {
			return nil, overlappingErrMsg
		}
	}

	result, err := s.DbService.SaveProposals(proposals)
	return result, err
}

*/

func (s EmployeeService) CreateProposalsString(proposalPayloadArr []model.ProposalStringPayload, id string) (interface{}, error) {
	const shortForm = "2006-Jan-02"
	proposalArr, err := CraftProposalFromPayloadString(proposalPayloadArr)
	if err != nil {
		return nil, err
	}

	if StartDateExceedsEndDate(proposalArr) {
		startExceedsEndErrMsg := errors.New("The startdate must be before the enddate")
		return nil, startExceedsEndErrMsg
	}

	overlappingErrMsg := errors.New("There cant be overlapping proposals.")
	actualProposals, err := s.GetProposalsByID(id)
	var actualProposalsString []model.Proposal
	for _, ps := range actualProposals {
		newP := ps

		Start := strings.Split(ps.StartDate, " ")
		End := strings.Split(ps.EndDate, " ")
		newTIme, _ := CreateTimeStringObject(Start[0], End[0])
		newP.TimeObject = newTIme
		actualProposalsString = append(actualProposalsString, newP)
	}

	var proposals []interface{}
	for _, p := range proposalArr {
		if !ProposalTimeIntersectsProposals(p, actualProposalsString) {
			proposals = append(proposals, p)
		} else {
			return nil, overlappingErrMsg
		}
	}

	result, err := s.DbService.SaveProposals(proposals)
	return result, err
}

func CreateTimeObject(start, end time.Time) (model.ProposalTimeObject, error) {

	Interval, err := interval.MakeTimeInterval(&start, &end)
	obj := model.ProposalTimeObject{
		Duration: Interval.Duration(),
		Interval: Interval,
		//		Err:      err,
	}
	return obj, err
}

func CreateTimeStringObject(start, end string) (model.ProposalTimeObject, error) {
	const shortForm = "2006-Jan-02"
	Start := strings.Split(start, " ")
	End := strings.Split(end, " ")

	Interval, err := interval.MakeTimeIntervalFromStrings(Start[0], End[0], shortForm)
	obj := model.ProposalTimeObject{
		Duration: Interval.Duration(),
		Interval: Interval,
		//		Err:      err,
	}
	return obj, err
}

func CraftProposalFromPayloadString(payload []model.ProposalStringPayload) ([]model.Proposal, error) {
	const shortForm = "2006-Jan-02"

	var proposals []model.Proposal
	for _, p := range payload {
		//		getStart, _ := time.Parse(shortForm, p.StartDate)
		//		getEnd, _ := time.Parse(shortForm, p.EndDate)
		obj, err := CreateTimeStringObject(p.StartDate, p.EndDate)
		newProposal := model.Proposal{
			UserId:     p.UserId,
			StartDate:  p.StartDate,
			EndDate:    p.EndDate,
			Approved:   false,
			Type:       p.Type,
			TimeObject: obj,
		}
		if err != nil {
			timeIntervalErrMsg := errors.New("Error occured in building the time interval for a new proposal")
			return nil, timeIntervalErrMsg
		}
		proposals = append(proposals, newProposal)
	}

	return proposals, nil
}

/*
func CraftProposalFromPayload(payload []model.ProposalPayload) ([]model.Proposal, error) {
	var proposals []model.Proposal
	for _, p := range payload {
		obj, err := CreateTimeObject(p.StartDate, p.EndDate)
		newProposal := model.Proposal{
			UserId:     p.UserId,
			StartDate:  p.StartDate,
			EndDate:    p.EndDate,
			Approved:   false,
			Type:       p.Type,
			TimeObject: obj,
		}
		if err != nil {
			timeIntervalErrMsg := errors.New("Error occured in building the time interval for a new proposal")
			return nil, timeIntervalErrMsg
		}
		proposals = append(proposals, newProposal)
	}

	return proposals, nil
}

*/

func ProposalTimeIntersectsProposals(proposal model.Proposal, Arr []model.Proposal) bool {
	for _, p := range Arr {
		p.TimeObject, _ = CreateTimeStringObject(p.StartDate, p.EndDate)

		if (*p.TimeObject.Interval.Start() == *proposal.TimeObject.Interval.Start()) || (*p.TimeObject.Interval.End() == *proposal.TimeObject.Interval.End()) {
			return true
		}

		if customOverlaps(p, proposal) {
			return true
		}
		/*
			if proposal.TimeObject.Interval.customOverlaps(p.TimeObject.Interval) {
				return true
			}

		*/

		if proposal.TimeObject.Interval.During(p.TimeObject.Interval) {
			return true
		}
		if p.TimeObject.Interval.During(proposal.TimeObject.Interval) {
			return true
		}
		if p.TimeObject.Interval.Equals(proposal.TimeObject.Interval) {
			return true
		}

	}
	return false
}

/*
func StartDateExceedsEndDate(Arr []model.Proposal) bool {
	for _, p := range Arr {
		p.TimeObject, _ = CreateTimeObject(p.StartDate, p.EndDate)
		if p.EndDate.Before(p.StartDate) {
			return true
		}
	}
	return false
}

func CheckOverlappingManually(proposal model.Proposal, Arr []model.Proposal) bool {
	for _, p := range Arr {
		//(StartDate1 <= EndDate2) and (StartDate2 <= EndDate1)
		stat1 := proposal.EndDate.Before(p.StartDate)
		stat2 := p.EndDate.Before(proposal.StartDate)
		if stat1 && stat2 {
			return true
		}
	}
	return false
}

*/

func Test(s string) {
	layout := "2006-01-02T15:04:05.000Z"

	t, err := time.Parse(layout, s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t)
}

/*
func CheckOverlappingManually(proposal model.Proposal, Arr []model.Proposal) bool {
	proposal.TimeObject, _ = CreateTimeStringObject(proposal.StartDate, proposal.EndDate)
	for _, p := range Arr {
		p.TimeObject, _ = CreateTimeStringObject(p.StartDate, p.EndDate)
		//(StartDate1 <= EndDate2) and (StartDate2 <= EndDate1)
		stat1 := proposal.EndDate.Before(p.StartDate)
		stat2 := p.EndDate.Before(proposal.StartDate)
		if stat1 && stat2 {
			return true
		}
	}
	return false
}

*/

func StartDateExceedsEndDate(Arr []model.Proposal) bool {
	for _, p := range Arr {
		p.TimeObject, _ = CreateTimeStringObject(p.StartDate, p.EndDate)
		if p.TimeObject.Interval.End().Before(*p.TimeObject.Interval.Start()) {
			return true
		}
	}
	return false
}

func customOverlaps(p1 model.Proposal, p2 model.Proposal) bool {
	if (*p1.TimeObject.Interval.Start() == *p2.TimeObject.Interval.Start()) && (p1.TimeObject.Interval.End().Before(*p2.TimeObject.Interval.End())) {
		return true
	}

	if (*p2.TimeObject.Interval.Start() == *p1.TimeObject.Interval.Start()) && (p2.TimeObject.Interval.End().Before(*p1.TimeObject.Interval.End())) {
		return true
	}

	return false
}
