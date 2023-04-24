package main

import (
	"bufio"
	"flag"
	"fmt"
	"gin-shop-api/internal/controllers"
	"gin-shop-api/internal/helpers"
	"gin-shop-api/internal/middleware"
	"gin-shop-api/internal/models"
	"gin-shop-api/internal/repository"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var action string

func init() {
	gin.ForceConsoleColor()
	helpers.LoadEnvVariables()
	repository.ConnectToDb()
}

func healthCheckRoute(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}

func RegisterRoutes(r *gin.Engine) {
	APIVersion := "/api/v1"

	// Auth API
	controllers.RegisterAuthRoutes(r.Group(APIVersion))

	// User API
	userCtrl := controllers.NewUserController(repository.DB)
	userCtrl.RegisterUserRoutes(r.Group(APIVersion))
}

func runServer() {
	if action == "run_server" {
		r := gin.Default()
		r.Use(gin.LoggerWithFormatter(
			func(param gin.LogFormatterParams) string {
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
		r.Use(gin.Recovery())
		r.Use(middleware.CORSMiddleware())
		healthCheckRoute(r)
		RegisterRoutes(r)

		r.Run()
	}
}

func createUserTable() {
	err := repository.DB.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}
}

func dropUserTable() {
	err := repository.DB.Migrator().DropTable(&models.User{})
	if err != nil {
		panic(err)
	}
}

func makeMigrations() {
	if action == "create_tables" {
		repository.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
		createUserTable()
		fmt.Println("Finished running migrations")
	}
}

func dropTables() {
	if action == "drop_tables" {
		dropUserTable()
		fmt.Println("Finished dropping tables")
	}
}

func createSuperUser() {
	if action == "create_super_user" {
		first_name := StringPrompt("first_name:")
		last_name := StringPrompt("last_name")
		email := StringPrompt("email")
		password := StringPrompt("password")

		user := models.User{
			ID:        uuid.New(),
			FirstName: first_name,
			LastName:  last_name,
			Email:     email,
			Password:  helpers.HashPassword(password),
			IsActive:  true,
		}
		result := repository.DB.Create(&user)

		if result.Error != nil {
			panic(result.Error)
		}
		fmt.Println("Finished creating super user record")
	}
}

func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func main() {
	flag.StringVar(&action, "action", "", "action e.g. run_server, create_tables, drop_tables")
	flag.Parse()
	runServer()
	makeMigrations()
	dropTables()
	createSuperUser()
}
