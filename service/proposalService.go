package service

import (
	"errors"
	"example-project/model"
	"example-project/utilities"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"strings"
	"time"
)

func (s EmployeeService) GetProposalsByID(id string) ([]model.Proposal, error) {

	result, err := s.DbService.GetProposalsByUserID(id)
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

func (s EmployeeService) GetAllProposals(ctx *gin.Context) ([]model.ProposalsByUser, error) {
	var proposalUserArray []model.ProposalsByUser
	var users []model.UserPayload
	var err error

	userid, useridok := ctx.GetQuery("userid")
	users, err = s.DbService.GetAllUser()
	if useridok {
		userid, _ := primitive.ObjectIDFromHex(userid)
		users, err = s.DbService.GetUserTeamMembersByID(userid)
	}

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

		filter, sort := utilities.FormGetAllProposalsFilter(user.ID.Hex(), ctx)
		proposals, err := s.DbService.GetProposalsByFilter(filter, sort)
		if err != nil {
			return proposalUserArray, err
		}
		if len(proposals) == 0 {
			continue
		}
		for _, prop := range proposals {
			if prop.Type == "sickness" {
				proposalUserItem.SicknessProposals = append(proposalUserItem.SicknessProposals, prop)
			}
			if prop.Type == "vacation" {
				proposalUserItem.VacationProposals = append(proposalUserItem.VacationProposals, prop)
			}
		}
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
	proposals, err := s.DbService.GetProposalsByUserID(userid)
	if err != nil {
		return model.AbsenceObject{}, err
	}

	var totalSicknessDays = 0
	var totalVacationDays = 0
	var VacationResult int
	var daysPerMonth = float64(user.VacationDays) / 12
	VacationResult = user.VacationDays
	if user.EntryTime.Year() == time.Now().Year() {
		var lastDate = time.Date(user.EntryTime.Year(), 12, 31, 0, 0, 0, 0, time.UTC)
		lastMonths := float64(lastDate.Month() - user.EntryTime.Month())
		VacationResult = int(math.RoundToEven(lastMonths * daysPerMonth))
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
		days := utilities.GetWeekdaysBetween(startDate, endDate)
		if proposal.Type == "sickness" {
			totalSicknessDays = totalSicknessDays + int(days) + 1
		}
		if proposal.Type == "vacation" {
			totalVacationDays = totalVacationDays + int(days) + 1
		}
	}

	absenceTime.VacationDays = totalVacationDays
	absenceTime.SicknessDays = totalSicknessDays
	absenceTime.TotalVacationDays = VacationResult
	return absenceTime, nil
}
