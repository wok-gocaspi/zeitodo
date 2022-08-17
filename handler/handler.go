package handler

import (
	"encoding/json"
	"example-project/model"
	"fmt"
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
	LoginUser(username string, password string) (http.Cookie, error)
	LogoutUser(userid string) bool
	RefreshToken(token string) (string, error)
	AuthenticateUser(url string, method string, token string) (bool, error)
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
	var user model.UserSignupPayload
	body := c.Copy().Request.Body
	jsonString, _ := ioutil.ReadAll(body)
	err := json.Unmarshal(jsonString, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": err.Error(),
		})
	}
	result, err := handler.ServiceInterface.CreateUser(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, result)
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
		return
	}
	c.SetCookie(response.Name, response.Value, 3600, response.Path, response.Domain, response.Secure, response.HttpOnly)
	c.Status(200)
}

func (handler Handler) LogoutUserHandler(c *gin.Context) {
	_, err := c.Request.Cookie("token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}
	c.SetCookie("token", "", 0, "", "", true, true)
	c.Status(200)
	return
}

func (handler Handler) RefreshTokenHandler(c *gin.Context) {
	cookie, err := c.Request.Cookie("token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}
	token, err := handler.ServiceInterface.RefreshToken(cookie.Value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}
	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.Status(200)
}

func (handler Handler) PermissionMiddleware(c *gin.Context) {
	fmt.Println(c.Request.RequestURI)
	tokenCookie, err := c.Request.Cookie("token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}
	ok, err := handler.ServiceInterface.AuthenticateUser(c.Request.RequestURI, c.Request.Method, tokenCookie.Value)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}
	c.Next()
	return
}
