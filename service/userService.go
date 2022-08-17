package service

import (
	"crypto/sha256"
	"errors"
	"example-project/model"
	"example-project/routes"
	"example-project/utils"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
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

func (s EmployeeService) CreateUser(user model.UserSignupPayload) (interface{}, error) {
	result, _ := s.DbService.GetUserByUsername(user.Username)
	if len(result.ID) < 0 {
		return nil, errors.New("user already exists with this username")
	}

	newUser := model.UserSignup{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Username:  user.Username,
		Password:  sha256.Sum256([]byte(user.Password)),
		Group:     "user",
	}
	if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Username == "" || user.Password == "" {
		return nil, errors.New("invalid register form")
	}

	results, err := s.DbService.CreateUser(newUser)
	if err != nil {
		return nil, err
	}
	return results, nil
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

func (s EmployeeService) LoginUser(username string, password string) (http.Cookie, error) {
	userObj, err := s.DbService.GetUserByUsername(username)
	if err != nil {
		return http.Cookie{}, errors.New("invalid login")
	}
	hashedPassword := sha256.Sum256([]byte(password))
	if hashedPassword != userObj.Password {
		return http.Cookie{}, errors.New("invalid login")
	}
	token := utils.GenerateToken(userObj.ID)
	expDate := time.Now().Add(time.Minute * 5)
	tokenPayload := http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expDate,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	}
	return tokenPayload, nil
}

func (s EmployeeService) LogoutUser(userid string) bool {
	for key, _ := range SessionMap {
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
