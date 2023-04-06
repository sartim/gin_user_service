package main

import (
	"fmt"
	"gin-shop-api/internal/models"
	"gin-shop-api/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.ForceConsoleColor()
	repository.ConnectToDb()
}

func Setup() {
	repository.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
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
