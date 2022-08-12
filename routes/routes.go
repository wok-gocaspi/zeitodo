package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . HandlerInterface
type HandlerInterface interface {
	CreateEmployeeHandler(c *gin.Context)
	GetEmployeeHandler(c *gin.Context)
	GetProposalsById(c *gin.Context)
	//	CreateProposalsHandler(c *gin.Context)
	CreateProposalsHandler(c *gin.Context)
	DeleteProposalHandler(c *gin.Context)
	UpdateProposalsHandler(c *gin.Context)
	DeleteTimeEntry(c *gin.Context)
	CreatTimeEntry(c *gin.Context)
	GetTimeEntryByUserID(c *gin.Context)
	GetAllTimeEntry(c *gin.Context)
	GetUserHandler(c *gin.Context)
	GetAllUserHandler(c *gin.Context)
	CreateUserHandler(c *gin.Context)
	GetTeamMemberHandler(c *gin.Context)
	UpdateUserHandler(c *gin.Context)
	DeleteUserHandler(c *gin.Context)
}

var Handler HandlerInterface

func CreateRoutes(group *gin.RouterGroup) {
	group.Use(CORS)
	user := group.Group("/user")
	user.GET("/:id/get", Handler.GetUserHandler)
	user.GET("/get", Handler.GetAllUserHandler)
	user.POST("/create", Handler.CreateUserHandler)
	user.GET("/team/get", Handler.GetTeamMemberHandler)
	user.PUT("/update", Handler.UpdateUserHandler)
	user.DELETE("/:id/delete", Handler.DeleteUserHandler)
	group.Use(CORS)
	timeentry := group.Group("/timeentry")
	timeentry.DELETE("/:id", Handler.DeleteTimeEntry)
	timeentry.POST("/", Handler.CreatTimeEntry)
	timeentry.GET("/:id", Handler.GetTimeEntryByUserID)
	timeentry.GET("/:id/getalltime", Handler.GetAllTimeEntry)

	proposal := group.Group("/proposals")
	proposal.GET("/:id", Handler.GetProposalsById)
	proposal.POST("/:id", Handler.CreateProposalsHandler)
	proposal.DELETE("/:id", Handler.DeleteProposalHandler)
	proposal.PATCH("/", Handler.UpdateProposalsHandler)
}
func CORS(c *gin.Context) {

	// First, we add the headers with need to enable CORS
	// Make sure to adjust these headers to your needs
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")

	// Second, we handle the OPTIONS problem
	if c.Request.Method != "OPTIONS" {

		c.Next()

	} else {

		// Everytime we receive an OPTIONS request,
		// we just return an HTTP 200 Status Code
		// Like this, Angular can now do the real
		// request using any other method than OPTIONS
		c.AbortWithStatus(http.StatusOK)
	}
}
