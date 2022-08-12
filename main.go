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

		//URL:      "mongodb://on4tdb:BMZQuk6pIL39nq46fOQPzygHtrhad5MFQMxs8YBQDW6YsJQSgbsIwO3aeOzlXXEnTjoz7ADVVr9jE1PKzU6GyQ==@on4tdb.mongo.cosmos.azure.com:10255/?ssl=true&replicaSet=globaldb&retrywrites=false&maxIdleTimeMS=120000&appName=@on4tdb@",
		URL:      os.Getenv("MONGO_URL"),
		Database: "zeitodo",
	})
	engine := middleware.SetupEngine([]gin.HandlerFunc{middleware.SetupService(databaseClient)})
	engine.Run(":" + port)
}
