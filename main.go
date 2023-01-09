package main

import (
	core "go-shop-api/app/core"
	"go-shop-api/app/routes"

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
