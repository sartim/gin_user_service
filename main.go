package main

import (
	"bufio"
	"flag"
	"fmt"
	core "gin-shop-api/app/core"
	"gin-shop-api/app/models"
	"gin-shop-api/app/routes"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var action string

func init() {
	gin.ForceConsoleColor()
	core.LoadEnvVariables()
	core.ConnectToDb()
}

func registerRoutes(r *gin.Engine) {
	routes.AuthRoutes(r)
	routes.UserRoutes(r)
}

func healthCheckRoute(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}

func runServer() {
	if action == "run_server" {
		r := gin.Default()
		r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
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
		healthCheckRoute(r)
		registerRoutes(r)
		r.Run()
	}
}

func makeMigrations() {
	if action == "make_migrations" {
		core.DB.AutoMigrate(&models.User{})
		core.DB.AutoMigrate(&models.Status{})
		fmt.Println("Finished running migrations")
	}
}

func dropTables() {
	if action == "drop_tables" {
		core.DB.Migrator().DropTable(&models.User{})
		core.DB.Migrator().DropTable(&models.Status{})
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
			FirstName: first_name,
			LastName:  last_name,
			Email:     email,
			Password:  core.HashPassword(password),
			IsActive:  true,
		}
		result := core.DB.Create(&user)

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
	flag.StringVar(&action, "action", "", "action e.g. run_server, make_migrations, drop_tables")
	flag.Parse()
	runServer()
	makeMigrations()
	dropTables()
	createSuperUser()
}
