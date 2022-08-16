package service

import (
	"example-project/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . DatabaseInterface
type DatabaseInterface interface {
	GetUserByID(id primitive.ObjectID) (model.UserPayload, error)
	GetUserByUsername(username string) (model.UserPayload, error)
	GetAllUser() ([]model.UserPayload, error)
	CreateUser(docs interface{}) (interface{}, error)
	GetUserTeamMembersByID(id primitive.ObjectID) (interface{}, error)
	GetUserTeamMembersByName(name string) (interface{}, error)
	UpdateManyUserByID(docs []model.User) []model.UserUpdateResult
	DeleteUser(id primitive.ObjectID) (interface{}, error)
}

type EmployeeService struct {
	DbService DatabaseInterface
}

func NewEmployeeService(dbInterface DatabaseInterface) EmployeeService {
	return EmployeeService{
		DbService: dbInterface,
	}
}
