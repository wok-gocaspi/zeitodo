package service

import "example-project/model"

func (s EmployeeService) DeleteTimeEntries(id string) (interface{}, error) {
	return s.DbService.DeleteTimeEntryById(id)
}

func (s EmployeeService) UpdateTimeEntries(update model.TimeEntry) (interface{}, error) {
	result, err := s.DbService.UpdateTimeEntryById(update)
	return result, err
}

func (s EmployeeService) CreatTimeEntries(id string) (interface{}, error) {
	return s.DbService.CreatTimeEntryById(id)

}

func (s EmployeeService) GetTimeEntries(id string) model.TimeEntry {
	return s.DbService.GetTimeEntryById(id)
}

func (s EmployeeService) GetAllTimeEntries(id string) model.TimeEntry {
	return s.DbService.GetAllTimeEntriesById(id)
}
