package service

import (
	"crypto/sha256"
	"errors"
	"example-project/model"
	"example-project/routes"
	"example-project/utils"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var SessionMap = make(map[string]string)

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

func (s EmployeeService) LoginUser(username string, password string) (string, error) {
	userObj, err := s.DbService.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("invalid login")
	}
	hashedPassword := sha256.Sum256([]byte(password))
	if hashedPassword != userObj.Password {
		return "", errors.New("invalid login")
	}
	token := utils.GenerateToken(userObj.ID)
	return token, nil
}

func (s EmployeeService) LogoutUser(userid string) bool {
	for key := range SessionMap {
		if key == userid {
			delete(SessionMap, key)
			return true
		}
	}
	return false
}

func (s EmployeeService) RefreshToken(token string) (string, error) {
	tkn, claims, err := utils.ValidateToken(token)
	if err != nil {
		return "", err
	}
	var userID string
	for key, val := range claims {
		if key == "userID" {
			userID = fmt.Sprint(val)
		}
	}
	if !tkn.Valid {
		return "", errors.New("invalid token")
	}
	tokenObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", err
	}
	tokenString := utils.GenerateToken(tokenObj)
	return tokenString, nil
}

func (s EmployeeService) AuthenticateUser(requestedURI string, requestMethod string, token string) (bool, error) {
	_, claims, err := utils.ValidateToken(token)
	if err != nil {
		return false, err
	}
	var userID string
	for key, val := range claims {
		if key == "userID" {
			userID = fmt.Sprint(val)
		}
	}

	userObj, err := s.GetUserByID(userID)
	if err != nil {
		return false, err
	}
	_, err = routes.PermissionList.CheckPolicy(requestedURI, requestMethod, userObj.Group, userID)
	if err != nil {
		return false, err
	}
	return true, nil
}
