package service

import (
	"example-project/model"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . DatabaseInterface
type DatabaseInterface interface {
	GetByID(id string) model.Employee
	GetProposals(id string) ([]model.Proposal, error)
	SaveProposals(docs []interface{}) (interface{}, error)
	DeleteProposalByIdAndDate(id string, date string) (*mongo.DeleteResult, error)
	UpdateProposal(update model.Proposal, date string) (*mongo.UpdateResult, error)
}

type EmployeeService struct {
	DbService DatabaseInterface
}

func NewEmployeeService(dbInterface DatabaseInterface) EmployeeService {
	return EmployeeService{
		DbService: dbInterface,
	}
}

func (s EmployeeService) GetEmployeeById(id string) model.Employee {
	return s.DbService.GetByID(id)
}
