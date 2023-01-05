package main

import (
	"go-shop-api/app/api"
	core "go-shop-api/app/core"

	"github.com/gin-gonic/gin"
)

func init() {
	core.LoadEnvVariables()
	core.ConnectToDb()
}

func main() {
	r := gin.Default()
	r.GET("/", api.RootApi)
	r.GET("/api/v1/user", api.UserGetApi)
	r.POST("/api/v1/user", api.UserCreateApi)
	r.PUT("/api/v1/user", api.UserUpdateApi)
	r.DELETE("/api/v1/user", api.UserDeleteApi)
	r.Run()
}
