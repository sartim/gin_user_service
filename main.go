package main

import (
	core "gin-shop-api/app/core"
	"gin-shop-api/app/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	core.LoadEnvVariables()
	core.ConnectToDb()
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	routes.AuthRoutes(r)
	routes.UserRoutes(r)
	r.Run()
}
