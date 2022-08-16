package service

import (
	"crypto/sha256"
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

func (s EmployeeService) CreateUser(userPayload model.UserSignupPayload) (interface{}, error) {
	hashedPassword := sha256.Sum256([]byte(userPayload.Password))
	user := model.UserSignup{Username: userPayload.Username, Password: hashedPassword, Email: userPayload.Email, FirstName: userPayload.FirstName, LastName: userPayload.LastName, Group: "user"}
	checkDBEmpty, _ := s.DbService.GetAllUser()
	fmt.Println(len(checkDBEmpty))
	if len(checkDBEmpty) == 0 {
		user.Group = "admin"
	}

	getUserByUsername, _ := s.DbService.GetUserByUsername(user.Username)
	if len(getUserByUsername.Username) > 0 {
		return nil, errors.New("user already exists, please choose another username")
	}

	if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Username == "" || userPayload.Password == "" {
		return nil, errors.New("insufficent user data")
	}
	results, err := s.DbService.CreateUser(user)
	if err != nil {
		return nil, err
	}
	userID := results.(primitive.ObjectID)
	userResult := model.UserSignupResult{ID: userID, FirstName: user.FirstName, LastName: user.LastName, Username: user.Username, Email: user.Email}
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
