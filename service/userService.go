package service

import (
	"errors"
	"example-project/model"
	"fmt"
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

func (s EmployeeService) CreateUser(user model.UserSignupPayload) (interface{}, error) {

	getUser, _ := s.DbService.GetUserByUsername(user.Username)
	if len(getUser.Username) > 0 {
		return nil, errors.New("user already exists, please choose another username")
	}

	if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Username == "" || user.Password == "" {
		return nil, errors.New("insufficent user data")
	}
	results, err := s.DbService.CreateUser(user)
	userID := results.(primitive.ObjectID)
	fmt.Println(userID)

	userResult := model.UserSignupResult{ID: userID, FirstName: user.FirstName, LastName: user.LastName, Username: user.Username, Email: user.Email}
	if err != nil {
		return nil, err
	}
	return userResult, nil
}

func (s EmployeeService) UpdateUsers(users []model.User) (interface{}, error) {
	result := s.DbService.UpdateManyUserByID(users)
	var errorUser []model.UserUpdateResult
	for _, updateResult := range result {
		if updateResult.Success == false {
			errorUser = append(errorUser, updateResult)
		}
	}
	if len(errorUser) > 0 {
		return errorUser, errors.New("a few users couldn't be updated")
	}
	return result, nil
}

func (s EmployeeService) DeleteUsers(id string) (interface{}, error) {
	idObject, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	result, err := s.DbService.DeleteUser(idObject)
	return result, err
}
