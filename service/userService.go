package service

import (
	"crypto/sha256"
	"errors"
	"example-project/model"
	"time"

	"example-project/utilities"

	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
	user := model.UserSignup{Username: userPayload.Username, Password: hashedPassword, Email: userPayload.Email, FirstName: userPayload.FirstName, LastName: userPayload.LastName, Group: "user", EntryTime: time.Now(), VacationDays: 30, HoursPerWeek: 40}
	checkDBEmpty, _ := s.DbService.GetAllUser()
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

func (s EmployeeService) UpdateUsers(users []model.UpdateUserPayload, userID string, userGroup string) (interface{}, error) {
	var totalSuccess = true
	var userUpdateResults []model.UserUpdateResult
	for _, user := range users {
		var UMR model.UserUpdateResult
		searchUser, err := s.DbService.GetUserByID(user.ID)
		if err != nil || len(searchUser.Username) == 0 {
			totalSuccess = false
			UMR.Error = "no correlating user found"
			UMR.Success = false
			UMR.ID = user.ID
			userUpdateResults = append(userUpdateResults, UMR)
			continue
		}
		filter := bson.M{"_id": user.ID}
		if user.ID.Hex() == userID || userGroup == "admin" {
			userSetter, err := utilities.UserUpdateSetter(user, userGroup)
			if err != nil {
				totalSuccess = false
				UMR.Error = err.Error()
				UMR.Success = false
				UMR.UpdateResult = nil
				UMR.ID = user.ID
				userUpdateResults = append(userUpdateResults, UMR)
				continue
			}
			updateResult, err := s.DbService.UpdateUserByID(filter, userSetter)
			if err != nil {
				totalSuccess = false
				UMR.Error = err.Error()
				UMR.Success = false
				UMR.UpdateResult = updateResult
				UMR.ID = user.ID
				userUpdateResults = append(userUpdateResults, UMR)
				continue
			}
			UMR.Success = true
			UMR.UpdateResult = updateResult
			UMR.ID = user.ID
			userUpdateResults = append(userUpdateResults, UMR)
			continue
		} else {
			totalSuccess = false
			UMR.Success = false
			UMR.ID = user.ID
			UMR.Error = "unauthorized user update, users can only update themself"
			userUpdateResults = append(userUpdateResults, UMR)
			continue
		}

	}
	var resultInterface interface{}
	resultInterface = userUpdateResults
	if !totalSuccess {
		return resultInterface, errors.New("users couldn't be updated")
	}
	return resultInterface, nil

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
	token := utilities.GenerateToken(userObj.ID)
	return token, nil
}

/*
func (s EmployeeService) GetUserId(username string token string) (string, error) {

    filter := bson.M{"username": username, "token": token}

	userObj, err := s.DbService.GetUserByUsername(username)

	token, er := utilities.ValidateToken(userObj)

	result,err :=s.DbService.GetUserById(username , token)
	if err != nil {
		return "", errors.New("no user found to that username")
	}

	return userObj.ID.Hex(), result, nil
}

*/

func (s EmployeeService) GetUserId(username string) (string, error) {

	userObj, err := s.DbService.GetUserByUsername(username)

	if err != nil {
		return "", errors.New("no user found to that username")
	}

	return userObj.ID.Hex(), nil
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

	tkn, claims, err := utilities.ValidateToken(token)

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
	tokenString := utilities.GenerateToken(tokenObj)
	return tokenString, nil
}

func (s EmployeeService) AuthenticateUser(token string) (string, string, error) {

	_, claims, err := utilities.ValidateToken(token)
	if err != nil {
		return "", "", err
	}
	var userID string

	for key, val := range claims {
		if key == "userID" {
			userID = fmt.Sprint(val)
		}
	}
	userObj, err := s.GetUserByID(userID)
	if err != nil {
		return userID, "", err
	}

	return userID, userObj.Group, nil
}

//****************************************

func (s EmployeeService) CheckUserPolicy(c *gin.Context, pl model.PermissionList) error {
	_, err := pl.CheckPolicy(c)

	if err != nil {
		return err
	}
	return nil
}

func (s EmployeeService) CheckIsSameUser(c *gin.Context, pl model.PermissionList, userid string) error {

	result := pl.IsSameUser(c, userid)

	if result == true {
		return nil
	}
	return errors.New("requesting user data of other users is not allowed")
}
