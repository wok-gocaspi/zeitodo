package handler

import (
	"example-project/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ServiceInterface
type ServiceInterface interface {
	CreateEmployees(employees []model.Employee) interface{}
	GetEmployeeById(id string) model.Employee
	DeleteTimeEntries(id string) (interface{}, error)
	UpdateTimeEntries(update model.TimeEntry) (interface{}, error)
	GetTimeEntries(id string) model.TimeEntry
	CreatTimeEntries(id string) (interface{}, error)
	GetAllTimeEntries(id string) model.TimeEntry
}

type Handler struct {
	ServiceInterface ServiceInterface
}

func NewHandler(serviceInterface ServiceInterface) Handler {
	return Handler{
		ServiceInterface: serviceInterface,
	}
}
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

/*
func (handler Handler) UpdateTimeEntry(context *gin.Context) {
	id, ok := context.Params.Get("id")
	if !ok {
		context.AbortWithStatusJSON(401, "No Time was submitted")
		return
	}
	response := handler.ServiceInterface.GetTimeEntryByUserID(id)
	if response.UserId == "" {
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
	update := model.TimeEntry{}
	result, err := handler.ServiceInterface.UpdateTimeEntries(update)
	if err != nil {
		context.AbortWithStatusJSON(400, err.Error())
		return
	}
	context.JSON(200, result)
}
*/

func (handler Handler) CreatTimeEntry(c *gin.Context) {
	var timeEntry model.TimeEntry
	err := c.BindJSON(&timeEntry)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": "Time is not created ",
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

func (handler Handler) GetTimeEntryByUserID(c *gin.Context) {
	pathParam, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": "time is not given",
		})
		return
	}

	response := handler.ServiceInterface.GetTimeEntryByUserID(pathParam)
	dt := time.Now()
	{
		fmt.Println("Current date and time is : ", dt.String())
	}
	c.JSON(http.StatusOK, response)

}

func (handler Handler) GetAllTimeEntry(c *gin.Context) {
	pathParam, ok := c.Params.Get("time")
	pages, pageOk := c.GetQuery("page")
	limit, limitOk := c.GetQuery("limit")
	_, pageErr := strconv.Atoi(pages)
	_, limitErr := strconv.Atoi(limit)

	if pageOk && limitOk {
		if pageOk && limitOk && pageErr == nil && limitErr == nil {

			response := handler.ServiceInterface.GetAllTimeEntries(pathParam)
			if !ok {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"errorMessage": "Alltime not given",
				})
				return
			}
			c.JSON(http.StatusOK, response)
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errorMessage": "queries are invalid, please check or remove them",
			})
			return
		}
	} else {

		_ = 1
		_ = 1000000 * 100000

		response := handler.ServiceInterface.GetAllTimeEntries(pathParam)

		c.JSON(http.StatusOK, response)
	}

}
