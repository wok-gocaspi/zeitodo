package main

import (
	"example-project/datasource"
	"example-project/middleware"
	"example-project/model"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	var port string
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	env := os.Getenv("environment")
	if env == "production" {
		port = "80"
	} else {
		port = "9090"
	}
	databaseClient := datasource.NewDbClient(model.DbConfig{
		URL:      os.Getenv("MONGO_URL"),
		Database: "zeitodo",
	})
	engine := middleware.SetupEngine([]gin.HandlerFunc{middleware.SetupService(databaseClient)})
	engine.Run(":" + port)
}
