package main

import (
	"fmt"
	"gin-shop-api/internal/models"
	"gin-shop-api/internal/repository"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	gin.ForceConsoleColor()
	err := godotenv.Load("../.env")
	if err != nil {
		log.Panicf("%s: %s", "Error loading env vars", err)
	}
	repository.ConnectToDb()
}

func Setup() {
	repository.DB.AutoMigrate(&models.User{})
	repository.DB.AutoMigrate(&models.Status{})
	fmt.Println("Finished running migrations")
}

func TearDown() {
	repository.DB.Migrator().DropTable(&models.User{})
	repository.DB.Migrator().DropTable(&models.Status{})
	fmt.Println("Finished dropping tables")
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())
	return router
}
