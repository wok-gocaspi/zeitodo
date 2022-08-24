package handler

import (
	"encoding/json"
	"example-project/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"net/http"
	"time"
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
	GetProposalsByID(id string) ([]model.Proposal, error)
	CreateProposals(proposalPayloadArr []model.ProposalPayload, id string) (interface{}, error)
	DeleteProposalsByID(id string, date string) error
	UpdateProposalByDate(update model.Proposal, date string) (*mongo.UpdateResult, error)
	CreatTimeEntries(te model.TimeEntry) (interface{}, error)
	UpdateTimeEntries(update model.TimeEntry) (interface{}, error)
	GetTimeEntries(id string) []model.TimeEntry
	DeleteTimeEntries(id string) (interface{}, error)
	GetAllTimeEntries() ([]model.TimeEntry, error)
	CollideTimeEntry(a, b model.TimeEntry) bool
	CalcultimeEntry(userid string) (map[string]float64, error)
}

type Handler struct {
	ServiceInterface ServiceInterface
}

func NewHandler(serviceInterface ServiceInterface) Handler {
	return Handler{
		ServiceInterface: serviceInterface,
	}
}

//Creat TimeEntries
func (handler Handler) CreatTimeEntry(c *gin.Context) {

	var timeEntry model.TimeEntry
	err := c.BindJSON(&timeEntry)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": "Time User is not created ",
		})
		return
	}
	response, err := handler.ServiceInterface.CreatTimeEntries(timeEntry)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)
}
func (handler Handler) CalcultimeEntry(context *gin.Context) {
	pathParam, ok := context.Params.Get("id")
	if !ok {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": "time is not given",
		})
		return
	}

	response, err := handler.ServiceInterface.CalcultimeEntry(pathParam)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": err.Error(),
		})
	}
	dt := time.Now()
	{
		fmt.Println("Current date and time is : ", dt.String())
	}
	context.JSON(http.StatusOK, response)

}

func (handler Handler) UpdateTimeEntry(context *gin.Context) {
	id, ok := context.Params.Get("id")
	if !ok {
		context.AbortWithStatusJSON(401, "No Time was submitted")
		return
	}
	response := handler.ServiceInterface.GetTimeEntries(id)
	if response == nil {
		context.AbortWithStatusJSON(400, "Time user ist not existing ")
		return
	}
	var payLoad model.TimeEntry
	err := context.ShouldBindJSON(&payLoad)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": "invalid payload",
		})
		return
	}
	result, err := handler.ServiceInterface.UpdateTimeEntries(payLoad)
	if err != nil {
		context.AbortWithStatusJSON(400, err.Error())
		return
	}
	context.JSON(200, result)
}
func (handler Handler) GetTimeEntry(c *gin.Context) {
	pathParam, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": "time is not given",
		})
		return
	}

	response := handler.ServiceInterface.GetTimeEntries(pathParam)
	dt := time.Now()
	{
		fmt.Println("Current date and time is : ", dt.String())
	}
	c.JSON(http.StatusOK, response)

}

//Delete TimeEntry
func (handler Handler) DeleteTimeEntry(c *gin.Context) {
	pathParam, ok := c.Params.Get("id")

	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{

			"errorMessage": "TimeEntry is not existing",
		})
		return
	}
	response, err := handler.ServiceInterface.DeleteTimeEntries(pathParam)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

//GetAllTimeEntry

func (handler Handler) GetAllTimeEntry(c *gin.Context) {

	response, err := handler.ServiceInterface.GetAllTimeEntries()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)
	return

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

const idNotFoundMsg = "id is not given"

func (handler Handler) DeleteProposalHandler(c *gin.Context) {
	id, idOk := c.Params.Get("id")
	if !idOk {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": idNotFoundMsg,
		})
		return
	}

	date, dateOk := c.GetQuery("date")
	if !dateOk {
		noQueryError := "No date was given in the query parameter!"
		c.AbortWithStatusJSON(404, gin.H{
			"errorMessage": noQueryError,
		})
		return
	}

	responseErr := handler.ServiceInterface.DeleteProposalsByID(id, date)
	if responseErr != nil {
		c.AbortWithStatusJSON(404, gin.H{
			"errorMessage": responseErr.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, "")
}

func (handler Handler) GetProposalsById(context *gin.Context) {
	id, idOk := context.Params.Get("id")
	if !idOk {
		noQueryError := "No department was given in the query parameter!"
		context.AbortWithStatusJSON(404, gin.H{
			"errorMessage": noQueryError,
		})
		return
	}

	response, err := handler.ServiceInterface.GetProposalsByID(id)
	if err != nil {
		context.AbortWithStatusJSON(404, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}

	context.JSON(200, response)
	return

}

func (handler Handler) CreateProposalsHandler(c *gin.Context) {
	pathParam, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(404, gin.H{
			"errorMessage": idNotFoundMsg,
		})
		return
	}

	var payLoad []model.ProposalPayload
	err := c.BindJSON(&payLoad)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": "invalid payload",
		})
		return
	}

	response, err := handler.ServiceInterface.CreateProposals(payLoad, pathParam)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}
	c.JSON(200, response)
}

func (handler Handler) UpdateProposalsHandler(c *gin.Context) {
	/*
		pathParam, ok := c.Params.Get("id")
		if !ok {
			c.AbortWithStatusJSON(404, gin.H{
				"errorMessage": "id is not given",
			})
			return
		}

	*/

	date, dateOk := c.GetQuery("date")
	if !dateOk {
		noQueryError := "No date was given in the query parameter!"
		c.AbortWithStatusJSON(404, gin.H{
			"errorMessage": noQueryError,
		})
		return
	}

	var payLoad model.Proposal
	err := c.BindJSON(&payLoad)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": "invalid payload",
		})
		return
	}

	response, err := handler.ServiceInterface.UpdateProposalByDate(payLoad, date)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}
	c.JSON(200, response)
}
