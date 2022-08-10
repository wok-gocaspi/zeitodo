package service

import (
	"crypto/sha256"
	"errors"
	"example-project/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (s EmployeeService) CreateUser(usersSignupPayload []model.UserSignupPayload) (interface{}, error) {
	var userList []interface{}
	var errorArray []model.UserSignupPayload
	for _, user := range usersSignupPayload {
		newUser := model.UserSignup{
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Email:      user.Email,
			Username:   user.Username,
			Password:   sha256.Sum256([]byte(user.Password)),
			Permission: "user",
		}
		if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Username == "" || user.Password == "" {
			errorArray = append(errorArray, user)
		}
		userList = append(userList, newUser)
	}
	if len(errorArray) > 0 {
		return errorArray, errors.New("insufficent user data")
	}

	results, err := s.DbService.CreateUser(userList)
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

func (s EmployeeService) LoginUser(username string, password string) (interface{}, error) {
	userObj, err := s.DbService.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("invalid login")
	}
	hashedPassword := sha256.Sum256([]byte(password))
	if hashedPassword != userObj.Password {
		return nil, errors.New("invalid login")
	}
	token := AddUserSession(userObj.ID)
	tokenPayload := gin.H{
		"userid": userObj.ID.Hex(),
		"token":  token,
	}
	var tokenInterface interface{} = tokenPayload
	return tokenInterface, nil
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

func AddUserSession(userObjectID primitive.ObjectID) string {
	for key, _ := range SessionMap {
		if key == userObjectID.Hex() {
			delete(SessionMap, key)
		}
	}
	randomString := (uuid.New()).String()
	SessionMap[userObjectID.Hex()] = randomString
	return randomString
}
