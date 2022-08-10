package service

import (
	"errors"
	"example-project/model"
	"example-project/utilities"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
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

func (s EmployeeService) DeleteProposalsByID(id string, date string) error {
	result, err := s.DbService.DeleteProposalByIdAndDate(id, date)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		deleterror := errors.New("the Employee id is not existing")
		return deleterror
	}

	return nil
}

func (s EmployeeService) CreateProposals(proposalPayloadArr []model.ProposalPayload, id string) (interface{}, error) {
	proposalArr, err := CraftProposalFromPayload(proposalPayloadArr)
	if err != nil {
		return nil, err
	}

	if utilities.StartDateExceedsEndDate(proposalArr) {
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
		newTIme, _ := utilities.CreateTimeObject(Start[0], End[0])
		newP.TimeObject = newTIme
		actualProposalsString = append(actualProposalsString, newP)
	}

	var proposals []interface{}
	for _, p := range proposalArr {
		if !utilities.ProposalTimeIntersectsProposals(p, actualProposalsString) {
			proposals = append(proposals, p)
		} else {
			return nil, overlappingErrMsg
		}
	}

	result, err := s.DbService.SaveProposals(proposals)
	return result, err
}

func CraftProposalFromPayload(payload []model.ProposalPayload) ([]model.Proposal, error) {

	var proposals []model.Proposal
	for _, p := range payload {
		obj, err := utilities.CreateTimeObject(p.StartDate, p.EndDate)
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

func (s EmployeeService) UpdateProposalByDate(update model.Proposal, date string) (*mongo.UpdateResult, error) {
	result, err := s.DbService.UpdateProposal(update, date)
	return result, err
}
