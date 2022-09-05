package service

import (
	"crypto/sha256"
	"errors"
	"example-project/model"
	"example-project/routes"
	"example-project/utilities"
	"fmt"
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
	user := model.UserSignup{Username: userPayload.Username, Password: hashedPassword, Email: userPayload.Email, FirstName: userPayload.FirstName, LastName: userPayload.LastName, Group: "user"}
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
			UMR.Error = "no correlating user found"
			UMR.Success = false
			UMR.ID = user.ID
			userUpdateResults = append(userUpdateResults, UMR)
			continue
		}
		filter := bson.M{"_id": user.ID}
		var setElements bson.D
		if user.ID.Hex() == userID || userGroup == "admin" {
			if user.FirstName != "" {
				setElements = append(setElements, bson.E{Key: "firstname", Value: user.FirstName})
			}
			if user.LastName != "" {
				setElements = append(setElements, bson.E{Key: "lastname", Value: user.LastName})
			}
			if user.Email != "" {
				setElements = append(setElements, bson.E{Key: "email", Value: user.Email})
			}
			if user.Team != "" {
				setElements = append(setElements, bson.E{Key: "team", Value: user.Team})
			}
			if user.TotalWorkingHours != 0 {
				setElements = append(setElements, bson.E{Key: "totalWorkingHours", Value: user.TotalWorkingHours})
			}
			if user.VacationDays != 0 {
				setElements = append(setElements, bson.E{Key: "vacationDays", Value: user.VacationDays})
			}
			if user.Username != "" {
				setElements = append(setElements, bson.E{Key: "username", Value: user.Username})
			}
			if user.Password != "" {
				setElements = append(setElements, bson.E{Key: "password", Value: sha256.Sum256([]byte(user.Password))})
			}
			if userGroup == "admin" && user.Group != "" {
				setElements = append(setElements, bson.E{Key: "group", Value: user.Group})
			}
			if len(user.Projects) != 0 {
				setElements = append(setElements, bson.E{Key: "projects", Value: user.Projects})
			} else if len(setElements) == 0 {
				UMR.Success = false
				UMR.ID = user.ID
				UMR.Error = "no items changed"
				userUpdateResults = append(userUpdateResults, UMR)
				continue
			}
			setMap := bson.D{
				{"$set", setElements},
			}
			updateResult, err := s.DbService.UpdateUserByID(filter, setMap)
			if err != nil {
				UMR.Error = err.Error()
				UMR.Success = false
				UMR.UpdateResult = updateResult
				UMR.ID = user.ID
				userUpdateResults = append(userUpdateResults, UMR)
				continue
			} else {
				UMR.Success = true
				UMR.UpdateResult = updateResult
				UMR.ID = user.ID
				userUpdateResults = append(userUpdateResults, UMR)
				continue
			}
		} else {
			UMR.Success = false
			UMR.ID = user.ID
			UMR.Error = "unauthorized user update, users can only update themself"
			userUpdateResults = append(userUpdateResults, UMR)
			continue
		}

	}
	var resultInterface interface{}
	resultInterface = userUpdateResults
	for _, ur := range userUpdateResults {
		if ur.Success == false {
			totalSuccess = false
		}
	}
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

func (s EmployeeService) AuthenticateUser(requestedURI string, requestMethod string, token string) (string, string, error) {
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
	_, err = routes.PermissionList.CheckPolicy(requestedURI, requestMethod, userObj.Group, userID)
	if err != nil {
		return userID, userObj.Group, err
	}
	return userID, userObj.Group, nil
}
