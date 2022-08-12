package service

import (
	"example-project/model"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . DatabaseInterface
type DatabaseInterface interface {
	UpdateMany(docs []interface{}) interface{}
	GetByID(id string) model.Employee
	DeleteTimeEntryById(id string) (interface{}, error)
	UpdateTimeEntryById(update model.TimeEntry) (*mongo.UpdateResult, error)
	CreatTimeEntryById(te model.TimeEntry) (interface{}, error)
	GetTimeEntryByUserID(id string) []model.TimeEntry
	GetAllTimeEntriesById(id string) model.TimeEntry
}

type EmployeeService struct {
	DbService DatabaseInterface
}

func NewEmployeeService(dbInterface DatabaseInterface) EmployeeService {
	return EmployeeService{
		DbService: dbInterface,
	}
}

func (s EmployeeService) CreateEmployees(employees []model.Employee) interface{} {

	var emp []interface{}
	for _, employee := range employees {
		emp = append(emp, employee)

	}
	return s.DbService.UpdateMany(emp)
}

func (s EmployeeService) GetEmployeeById(id string) model.Employee {
	return s.DbService.GetByID(id)
}
