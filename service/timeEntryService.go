package service

import "example-project/model"

func (s EmployeeService) DeleteTimeEntries(id string) (interface{}, error) {
	return s.DbService.DeleteTimeEntryById(id)
}

func (s EmployeeService) UpdateTimeEntries(update model.TimeEntry) (interface{}, error) {
	result, err := s.DbService.UpdateTimeEntryById(update)
	return result, err
}

func (s EmployeeService) CreatTimeEntries(te model.TimeEntry) (interface{}, error) {
	return s.DbService.CreatTimeEntryById(te)

}

func (s EmployeeService) GetTimeEntryByUserID(id string) []model.TimeEntry {
	return s.DbService.GetTimeEntryByUserID(id)
}

func (s EmployeeService) GetAllTimeEntries(id string) model.TimeEntry {
	return s.DbService.GetAllTimeEntriesById(id)
}
