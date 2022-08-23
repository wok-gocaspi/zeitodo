package service

import (
	"errors"
	"example-project/model"
)

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

func (s EmployeeService) GetTimeEntries(id string) []model.TimeEntry {
	return s.DbService.GetTimeEntryByID(id)
}
func (s EmployeeService) DeleteTimeEntries(id string) (interface{}, error) {
	return s.DbService.DeleteTimeEntryById(id)
}
func (s EmployeeService) GetAllTimeEntries() ([]model.TimeEntry, error) {
	return s.DbService.GetAllTimeEntry()
}
func (s EmployeeService) CollideTimeEntry(timevor, timenach model.TimeEntry) bool {

	if timevor.Start.Before(timenach.Start) && timevor.End.After(timenach.Start) {
		return true
	}

	if timevor.Start.Before(timenach.End) && timevor.End.After(timenach.End) {
		return true
	}

	if timenach.Start.Before(timevor.Start) && timenach.End.After(timevor.Start) {
		return true
	}

	if timenach.Start.Before(timevor.End) && timenach.End.After(timevor.End) {
		return true
	}
	if timevor.Start == timenach.Start {
		return true
	}
	if timevor.End == timenach.End {
		return true
	}
	return false
}
func (s EmployeeService) CalcultimeEntry(userid string) (map[string]float64, error) {

	m := make(map[string]float64)
	timeentries, err := s.DbService.GetAllTimeEntry()

	if err != nil {
		return nil, err
	}
	for _, timeentry := range timeentries {
		if timeentry.UserId != userid {
			continue
		}
		dur := timeentry.End.Sub(timeentry.Start)
		if _, prs := m[timeentry.Project]; prs {

			//Project in Map

			m[timeentry.Project] += dur.Hours()

		} else {
			m[timeentry.Project] = dur.Hours()
		}

	}
	return m, nil
}
