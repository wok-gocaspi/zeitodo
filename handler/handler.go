package handler

import (
	"encoding/json"
	"example-project/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ServiceInterface
type ServiceInterface interface {
	GetProposalsByID(id string) ([]model.Proposal, error)
	CreateProposals(proposalPayloadArr []model.ProposalPayload, id string) (interface{}, error)
	DeleteProposalsByID(id string, date string) error
	UpdateProposalByDate(update model.Proposal, date string) (*mongo.UpdateResult, error)
	DeleteTimeEntries(id string) (interface{}, error)
	UpdateTimeEntries(update model.TimeEntry) (interface{}, error)
	GetTimeEntryByUserID(id string) []model.TimeEntry
	CreatTimeEntries(te model.TimeEntry) (interface{}, error)
	GetAllTimeEntries(id string) model.TimeEntry
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

const idNotFoundMsg = "id is not given"

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
