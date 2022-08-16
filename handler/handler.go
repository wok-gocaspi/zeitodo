package handler

import (
	"encoding/json"
	"example-project/model"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ServiceInterface
type ServiceInterface interface {
	GetUserByID(id string) (model.UserPayload, error)
	GetAllUser() ([]model.UserPayload, error)
	CreateUser(model.UserSignupPayload) (interface{}, error)
	GetTeamMembersByUserID(id string) (interface{}, error)
	UpdateUsers(users []model.User) (interface{}, error)
	GetTeamMembersByName(name string) (interface{}, error)
	DeleteUsers(id string) (interface{}, error)
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
			"errorMessage": "id or name is not given",
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
	var user model.UserSignupPayload
	body := c.Copy().Request.Body
	jsonString, _ := ioutil.ReadAll(body)
	err := json.Unmarshal(jsonString, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}
	result, err := handler.ServiceInterface.CreateUser(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, result)
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
