package main

import (
	"example-project/datasource"
	"example-project/middleware"
	"example-project/model"
	"github.com/gin-gonic/gin"
)

func main() {
	databaseClient := datasource.NewDbClient(model.DbConfig{
		URL:      "mongodb://localhost:27017",
		Database: "office",
	})
	engine := middleware.SetupEngine([]gin.HandlerFunc{middleware.SetupService(databaseClient)})
	engine.Run(":9090")
}
