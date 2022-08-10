package handler

import (
	"encoding/json"
	"example-project/model"
	"example-project/service"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ServiceInterface
type ServiceInterface interface {
	GetUserByID(id string) (model.UserPayload, error)
	GetAllUser() ([]model.UserPayload, error)
	CreateUser([]model.UserSignupPayload) (interface{}, error)
	GetTeamMembersByUserID(id string) (interface{}, error)
	UpdateUsers(users []model.User) (interface{}, error)
	GetTeamMembersByName(name string) (interface{}, error)
	DeleteUsers(id string) (interface{}, error)
	LoginUser(username string, password string) (interface{}, error)
	LogoutUser(userid string) bool
}

type Handler struct {
	ServiceInterface ServiceInterface
}

func NewHandler(serviceInterface ServiceInterface) Handler {
	return Handler{
		ServiceInterface: serviceInterface,
	}
}

func (handler Handler) GetUserHandler(c *gin.Context) {
	pathParam, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": "id is not given",
		})
		return
	}
	response, err := handler.ServiceInterface.GetUserByID(pathParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)
	return

}

func (handler Handler) GetAllUserHandler(c *gin.Context) {
	response, err := handler.ServiceInterface.GetAllUser()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)
	return

}

func (handler Handler) GetTeamMemberHandler(c *gin.Context) {
	idParam, idOK := c.GetQuery("id")
	nameParam, nameOK := c.GetQuery("name")
	if !idOK && !nameOK {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": "id is not given",
		})
		return
	}
	if idOK {
		result, err := handler.ServiceInterface.GetTeamMembersByUserID(idParam)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"errorMessage": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, result)
	}
	if nameOK {
		result, err := handler.ServiceInterface.GetTeamMembersByName(nameParam)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"errorMessage": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, result)
	}

}

func (handler Handler) CreateUserHandler(c *gin.Context) {
	var users []model.UserSignupPayload
	var user model.UserSignupPayload
	body := c.Copy().Request.Body
	jsonString, _ := ioutil.ReadAll(body)
	err := json.Unmarshal(jsonString, &user)
	if err != nil {
		err := json.Unmarshal(jsonString, &users)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errorMessage": err.Error(),
			})
			return
		}
		result, err := handler.ServiceInterface.CreateUser(users)
		if err != nil && err.Error() == "insufficent user data" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"errorMessage": err.Error(),
				"errorUser":    result,
			})
			return
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"errorMessage": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, result)
		return
	}
	users = append(users, user)
	result, err := handler.ServiceInterface.CreateUser(users)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, result)
	return
}

func (handler Handler) UpdateUserHandler(c *gin.Context) {
	var user model.User
	var users []model.User
	body := c.Copy().Request.Body
	jsonString, _ := ioutil.ReadAll(body)
	err := json.Unmarshal(jsonString, &user)
	if err != nil {
		err := json.Unmarshal(jsonString, &users)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errorMessage": err.Error(),
			})
			return
		}
		result, err := handler.ServiceInterface.UpdateUsers(users)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"errorMessage": err.Error(),
				"errorUser":    result,
			})
			return
		}
		c.JSON(http.StatusOK, result)
		return
	}
	users = append(users, user)
	result, err := handler.ServiceInterface.UpdateUsers(users)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"errorMessage": err.Error(),
			"errorUser":    result,
		})
		return
	}
	c.JSON(http.StatusOK, result)
	return
}

func (handler Handler) DeleteUserHandler(c *gin.Context) {
	pathParam, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": "no id given...",
		})
		return
	}
	result, err := handler.ServiceInterface.DeleteUsers(pathParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (handler Handler) LoginUserHandler(c *gin.Context) {
	var authObj model.UserAuthPayload
	body := c.Copy().Request.Body
	jsonString, _ := ioutil.ReadAll(body)
	err := json.Unmarshal(jsonString, &authObj)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": "invalid data",
		})
		return
	}
	response, err := handler.ServiceInterface.LoginUser(authObj.Username, authObj.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"errorMessage": err.Error(),
		})
	}
	c.JSON(http.StatusOK, response)
	return
}

func (handler Handler) LogoutUserHandler(c *gin.Context) {
	userID, ok := c.GetQuery("userid")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": "no user id delivered",
		})
		return
	}
	logoutOK := handler.ServiceInterface.LogoutUser(userID)
	if !logoutOK {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"errorMessage": "user already logged out",
		})
		return
	}
	c.Status(http.StatusOK)
	return
}

func (handler Handler) PermissionMiddleware(c *gin.Context, permissionLevel string) {
	authHeader := c.GetHeader("Authorization")
	var userID string
	if len(authHeader) < 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"errorMessage": "no authorization header provided. please check if user is logged in",
		})
	}
	for key, value := range service.SessionMap {
		if value == authHeader {
			userID = key
		}
	}
	if len(userID) < 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"errorMessage": "user is currently not logged in, please log in again",
		})
	}
	userObj, err := handler.ServiceInterface.GetUserByID(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"errorMessage": "user dosent exist",
		})
	}
	if userObj.Permission == permissionLevel {
		c.Next()
		return
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"errorMessage": "user does not have permission to access this endpoint",
		})
	}
}
