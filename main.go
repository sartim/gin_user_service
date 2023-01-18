package main

import (
	"fmt"
	core "gin-shop-api/app/core"
	"gin-shop-api/app/routes"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	core.LoadEnvVariables()
	core.ConnectToDb()
}

func main() {
	gin.ForceConsoleColor()
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
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	routes.AuthRoutes(r)
	routes.UserRoutes(r)
	r.Run()
}
