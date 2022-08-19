package service

import "example-project/model"

func (s EmployeeService) CreatTimeEntries(te model.TimeEntry) (interface{}, error) {
	return s.DbService.CreatTimeEntryById(te)

}
func (s EmployeeService) UpdateTimeEntries(update model.TimeEntry) (interface{}, error) {
	result, err := s.DbService.UpdateTimeEntryById(update)
	return result, err
}
func (s EmployeeService) GetTimeEntries(id string) []model.TimeEntry {
	return s.DbService.GetTimeEntryByID(id)
}
