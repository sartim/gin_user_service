package main

import (
	"flag"
	"fmt"
	core "gin-shop-api/app/core"
	"gin-shop-api/app/models"
	"gin-shop-api/app/routes"
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

func main() {
	flag.StringVar(&action, "action", "", "action e.g. run_server, make_migrations, drop_tables")
	flag.Parse()
	runServer()
	makeMigrations()
	dropTables()
}
