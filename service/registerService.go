package service

import (
	"example-project/model"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . DatabaseInterface
type DatabaseInterface interface {

	GetProposals(id string) ([]model.Proposal, error)
	SaveProposals(docs []interface{}) (interface{}, error)
	DeleteProposalByIdAndDate(id string, date string) (*mongo.DeleteResult, error)
	UpdateProposal(update model.Proposal, date string) (*mongo.UpdateResult, error)
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