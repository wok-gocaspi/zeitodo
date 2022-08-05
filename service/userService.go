package service

import (
	"crypto/sha256"
	"example-project/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s EmployeeService) GetUserByID(id string) (model.UserPayload, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.UserPayload{}, err
	}
	result, err := s.DbService.GetUserByID(objectID)
	if err != nil {
		return model.UserPayload{}, err
	}
	return result, nil

}

func (s EmployeeService) GetAllUser() ([]model.UserPayload, error) {
	result, err := s.DbService.GetAllUser()
	if err != nil {
		return []model.UserPayload{}, err
	}
	return result, nil
}

func (s EmployeeService) GetTeamMembersByUserID(id string) (interface{}, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	result, err := s.DbService.GetUserTeamMembersByID(objectID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s EmployeeService) GetTeamMembersByName(team string) (interface{}, error) {

	result, err := s.DbService.GetUserTeamMembersByName(team)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s EmployeeService) CreateUser(usersSignupPayload []model.UserSignupPayload) (interface{}, error) {
	var userList []interface{}
	for _, user := range usersSignupPayload {
		newUser := model.UserSignup{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Username:  user.Username,
			Password:  sha256.Sum256([]byte(user.Password)),
		}
		userList = append(userList, newUser)
	}

	results, err := s.DbService.CreateUser(userList)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s EmployeeService) UpdateUsers(users []model.User) interface{} {
	result := s.DbService.UpdateManyUserByID(users)
	return result
}

func (s EmployeeService) DeleteUsers(id string) (interface{}, error) {
	idObject, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	result, err := s.DbService.DeleteUser(idObject)
	return result, err
}
