package service

import (
	"errors"
	"example-project/model"
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

func (s EmployeeService) CreateProposals(proposalArr []model.Proposal, id string) (interface{}, error) {
	actualProposals, err := s.GetProposalsByID(id)
	var proposals []interface{}
	for _, p := range proposalArr {
		if !DuplicatedStart(p, actualProposals) {
			proposals = append(proposals, p)
		}
	}

	result, err := s.DbService.SaveProposals(proposals)
	return result, err
}

func DuplicatedStart(proposal model.Proposal, Arr []model.Proposal) bool {
	for _, p := range Arr {
		if p.StartDate == proposal.StartDate {
			return true
		}
	}
	return false
}
