package middleware

import (
	"example-project/handler"
	"example-project/routes"
	"example-project/service"
	"github.com/gin-gonic/gin"
)

func SetupEngine(functions []gin.HandlerFunc) *gin.Engine {
	engine := gin.Default()
	for _, f := range functions {
		engine.Use(f)
	}

	routes.CreateRoutes(&engine.RouterGroup)
	return engine
}

func SetupService(dbClient service.DatabaseInterface) gin.HandlerFunc {

	service := service.NewEmployeeService(dbClient)
	routes.Handler = handler.NewHandler(service)

	return func(c *gin.Context) {
		c.Next()
	}
}
