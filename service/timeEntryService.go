package service

import (
	"errors"
	"example-project/model"
	"example-project/utilities"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (s EmployeeService) GetTimeEntries(id string) []model.TimeEntry {
	return s.DbService.GetTimeEntryByID(id)
}
func (s EmployeeService) DeleteTimeEntries(userId string, starttime time.Time) (interface{}, error) {
	return s.DbService.DeleteTimeEntryById(userId, starttime)
}
func (s EmployeeService) GetAllTimeEntries() ([]model.TimeEntry, error) {
	return s.DbService.GetAllTimeEntry()
}

func (s EmployeeService) CollideTimeEntry(timevor, timenach model.TimeEntry) bool {

	if timevor.End == timenach.End {
		return true
	}
	if timevor.Start == timenach.Start {
		return true
	}

	if timevor.Start.Before(timenach.Start) && timevor.End.After(timenach.Start) {
		return true
	}

	if timevor.Start.Before(timenach.End) && timevor.End.After(timenach.End) {
		return true
	}

	if timenach.Start.Before(timevor.Start) && timenach.End.After(timevor.Start) {
		return true
	}
	//if timenach.Start.Before(timevor.End) && timenach.End.After(timevor.End) {
	//	return true
	//}
	return false
}
func (s EmployeeService) CreatTimeEntries(te model.TimeEntry) (interface{}, error) {

	var timeentriesDb []model.TimeEntry

	timeentriesDb, err := s.DbService.GetAllTimeEntry()

	if err != nil {
		return nil, err
	}
	for _, timeentry := range timeentriesDb {

		if timeentry.UserId != te.UserId {
			continue
		}
		if s.CollideTimeEntry(te, timeentry) {
			err = errors.New("check if TimeEntries collide")
			return nil, err
		}
	}
	return s.DbService.CreatTimeEntryById(te)
}
func (s EmployeeService) UpdateTimeEntries(update model.TimeEntry) (interface{}, error) {

	var timeentriesDb []model.TimeEntry

	timeentriesDb, err := s.DbService.GetAllTimeEntry()

	if err != nil {
		return nil, err
	}

	for _, timeentry := range timeentriesDb {
		if timeentry.UserId != update.UserId {
			continue
		}
		if s.CollideTimeEntry(update, timeentry) && timeentry.Start != update.Start {
			err = errors.New("check if update TimeEntries collide")
			return nil, err
		}
	}
	return s.DbService.UpdateTimeEntryById(update)
}

func (s EmployeeService) CalculateTimeEntries(ctx *gin.Context) (model.WorkingHoursPayload, error) {
	var workingPayload model.WorkingHoursPayload

	userID, userIDOK := ctx.GetQuery("userid")
	if !userIDOK {
		return workingPayload, errors.New("no user id supplied")
	}
	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return workingPayload, err
	}
	user, err := s.DbService.GetUserByID(userIDObj)

	proposalFilter, proposalSort := utilities.FormGetAllProposalsFilter(user, ctx)
	proposals, err := s.DbService.GetProposalsByFilter(proposalFilter, proposalSort)
	if err != nil {
		return workingPayload, err
	}

	requiredWorkingHours, err := utilities.CalculateRequiredWorkingHours(user, proposals, ctx)
	if err != nil {
		return workingPayload, err
	}
	workingPayload.Required = requiredWorkingHours

	filter, err := utilities.FormGetTimeEntryFilter(ctx)
	if err != nil {
		return workingPayload, err
	}
	timeEntries, err := s.DbService.GetTimeEntriesByFilter(filter)

	if err != nil {
		return workingPayload, err
	}
	m := make(map[string]float64)
	for _, timeEntry := range timeEntries {
		dur := timeEntry.End.Sub(timeEntry.Start)
		breakDur := timeEntry.BreakEnd.Sub(timeEntry.BreakStart)
		dur -= breakDur
		workingPayload.Actual += dur.Hours()
		if _, prs := m[timeEntry.Project]; prs {
			m[timeEntry.Project] += dur.Hours()
		} else {
			m[timeEntry.Project] = dur.Hours()
		}

	}
	workingPayload.Projects = m
	return workingPayload, nil
}
