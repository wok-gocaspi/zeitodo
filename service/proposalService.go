package service

import (
	"errors"
	"example-project/model"
	"example-project/utilities"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	proposalArr, err := utilities.CraftProposalFromPayload(proposalPayloadArr)
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

//***************************************

func (s EmployeeService) UpdateProposalByDate(update model.Proposal, date string, ctx *gin.Context) (*mongo.UpdateResult, error) {

	if update.UserId == ctx.GetString("userid") || ctx.GetString("group") == "admin" {
		result, err := s.DbService.UpdateProposal(update, date)
		return result, err
	} else {
		return nil, errors.New("user can not update ")
	}
}

func (s EmployeeService) GetAllProposals() ([]model.ProposalsByUser, error) {
	var proposalUserArray []model.ProposalsByUser
	users, err := s.DbService.GetAllUser()
	if err != nil {
		return proposalUserArray, err
	}
	for _, user := range users {
		var proposalUserItem model.ProposalsByUser
		proposalUserItem.Userid = user.ID
		proposalUserItem.Username = user.Username
		proposalUserItem.Email = user.Email
		proposalUserItem.FirstName = user.FirstName
		proposalUserItem.LastName = user.LastName

		proposals, err := s.DbService.GetProposals(user.ID.Hex())
		if err != nil {
			return proposalUserArray, err
		}
		if len(proposals) == 0 {
			continue
		}
		proposalUserItem.Proposals = append(proposalUserItem.Proposals, proposals...)
		proposalUserArray = append(proposalUserArray, proposalUserItem)

	}
	return proposalUserArray, nil
}

func (s EmployeeService) GetTotalAbsence(userid string) (model.AbsenceObject, error) {
	var absenceTime model.AbsenceObject
	objectID, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return absenceTime, err
	}
	user, err := s.DbService.GetUserByID(objectID)
	if err != nil {
		return absenceTime, err
	}
	proposals, err := s.DbService.GetProposals(userid)
	if err != nil {
		return model.AbsenceObject{}, err
	}

	var totalSicknessDays int = 0
	var totalVacationDays int = 0
	var userTotalVacationDays float32 = float32(user.VacationDays)
	var daysPerMonth float32 = userTotalVacationDays / 12
	fmt.Println(daysPerMonth)
	if user.EntryTime.Year() == time.Now().Year() {
		var lastDate = time.Date(user.EntryTime.Year(), 12, 31, 0, 0, 0, 0, time.UTC)
		lastMonths := float32(lastDate.Month() - user.EntryTime.Month())
		fmt.Println(lastMonths)
		userTotalVacationDays = lastMonths * daysPerMonth
	}

	for _, proposal := range proposals {
		const layout = "2006-Jan-02"
		startDate, err := time.Parse(layout, proposal.StartDate)
		if err != nil {
			return model.AbsenceObject{}, err
		}
		endDate, err := time.Parse(layout, proposal.EndDate)
		if err != nil {
			return model.AbsenceObject{}, err
		}
		days := endDate.Sub(startDate).Hours() / 24
		if proposal.Type == "sickness" {
			totalSicknessDays = totalSicknessDays + int(days)
		}
		if proposal.Type == "vacation" {
			totalVacationDays = totalVacationDays + int(days)
		}
	}

	absenceTime.VacationDays = totalVacationDays
	absenceTime.SicknessDays = totalSicknessDays
	absenceTime.TotalVacationDays = int(userTotalVacationDays)
	return absenceTime, nil
}
