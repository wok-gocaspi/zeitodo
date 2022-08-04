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
