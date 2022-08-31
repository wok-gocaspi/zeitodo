package middleware

import (
	"example-project/handler"
	"example-project/routes"
	"example-project/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupEngine(functions []gin.HandlerFunc) *gin.Engine {
	engine := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("Authorization")
	engine.Use(cors.New(config))
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
