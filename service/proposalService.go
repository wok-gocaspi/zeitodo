package service

import (
	"errors"
	"example-project/model"
	"fmt"
	"github.com/retailify/go-interval"
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

func (s EmployeeService) CreateProposals(proposalPayloadArr []model.ProposalPayload, id string) (interface{}, error) {
	proposalArr, err := CraftProposalFromPayload(proposalPayloadArr)
	if err != nil {
		return nil, err
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

func CreateTimeObject(start, end time.Time) (model.ProposalTimeObject, error) {

	Interval, err := interval.MakeTimeInterval(&start, &end)
	obj := model.ProposalTimeObject{
		Duration: Interval.Duration(),
		Interval: Interval,
		//		Err:      err,
	}
	return obj, err
}

func OverlappingProposals(proposal model.Proposal, Arr []model.Proposal) bool {
	for _, p := range Arr {
		if p.TimeObject.Interval != nil && !(proposal.TimeObject.Interval.During(p.TimeObject.Interval) || p.TimeObject.Interval.During(proposal.TimeObject.Interval)) {
			return true
		}
	}
	return false
}

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

func OverlappingProposalsTest(proposal model.Proposal, Arr []model.Proposal) bool {
	for _, p := range Arr {
		p.TimeObject, _ = CreateTimeObject(p.StartDate, p.EndDate)
		if p.TimeObject.Interval != nil {
			intervalP := p.TimeObject.Interval
			intervalProposal := proposal.TimeObject.Interval
			fmt.Println(intervalP, intervalProposal)
			statment := p.TimeObject.Interval.Overlaps(proposal.TimeObject.Interval)
			fmt.Println(statment)
			return true
		}
	}
	return false
}

func ProposalTimeIntersectsProposals(proposal model.Proposal, Arr []model.Proposal) bool {
	for _, p := range Arr {
		p.TimeObject, _ = CreateTimeObject(p.StartDate, p.EndDate)
		if p.TimeObject.Interval.Overlaps(proposal.TimeObject.Interval) {
			return true
		}
		if proposal.TimeObject.Interval.Overlaps(p.TimeObject.Interval) {
			return true
		}
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
