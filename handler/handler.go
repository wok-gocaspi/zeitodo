package handler

import (
	"example-project/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ServiceInterface
type ServiceInterface interface {
	//	CreateEmployees(employees []model.Employee) interface{}
	GetEmployeeById(id string) model.Employee
	GetProposalsByID(id string) ([]model.Proposal, error)
	CreateProposals(proposalPayloadArr []model.ProposalPayload, id string) (interface{}, error)
}

type Handler struct {
	ServiceInterface ServiceInterface
}

func NewHandler(serviceInterface ServiceInterface) Handler {
	return Handler{
		ServiceInterface: serviceInterface,
	}
}

/*
func (handler Handler) CreateEmployeeHandler(c *gin.Context) {
	var payLoad model.Payload
	err := c.BindJSON(&payLoad)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": "invalid payload",
		})
		return
	}

	response := handler.ServiceInterface.CreateEmployees(payLoad.Employees)
	c.JSON(200, response)
}

*/

func (handler Handler) CreateProposalsHandler(c *gin.Context) {
	pathParam, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(404, gin.H{
			"errorMessage": "id is not given",
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

func (handler Handler) GetEmployeeHandler(c *gin.Context) {
	pathParam, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": "id is not given",
		})
		return
	}

	response := handler.ServiceInterface.GetEmployeeById(pathParam)
	fmt.Println(response)
	c.JSON(http.StatusOK, response)
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
