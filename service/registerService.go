package service

import (
	"example-project/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . DatabaseInterface

type DatabaseInterface interface {
	GetUserByID(id primitive.ObjectID) (model.UserPayload, error)
	GetUserByUsername(username string) (model.User, error)
	GetAllUser() ([]model.UserPayload, error)
	CreateUser(docs interface{}) (interface{}, error)
	GetUserTeamMembersByID(id primitive.ObjectID) (interface{}, error)
	GetUserTeamMembersByName(name string) (interface{}, error)
	UpdateUserByID(filter bson.M, setter bson.D) (*mongo.UpdateResult, error)
	DeleteUser(id primitive.ObjectID) (interface{}, error)
	GetProposals(id string) ([]model.Proposal, error)
	SaveProposals(docs []interface{}) (interface{}, error)
	DeleteProposalByIdAndDate(id string, date string) (*mongo.DeleteResult, error)
	UpdateProposal(update model.Proposal, date string) (*mongo.UpdateResult, error)
	CreatTimeEntryById(te model.TimeEntry) (interface{}, error)
	UpdateTimeEntryById(update model.TimeEntry) (*mongo.UpdateResult, error)
	GetTimeEntryByID(id string) []model.TimeEntry
	DeleteTimeEntryById(userId string, starttime time.Time) (interface{}, error)
	GetAllTimeEntry() ([]model.TimeEntry, error)
}

type EmployeeService struct {
	DbService DatabaseInterface
}

func NewEmployeeService(dbInterface DatabaseInterface) EmployeeService {
	return EmployeeService{
		DbService: dbInterface,
	}
}
