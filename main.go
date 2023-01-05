package main

import (
	"go-shop-admin/app/api"
	core "go-shop-admin/app/core"

	"github.com/gin-gonic/gin"
)

func init() {
	core.LoadEnvVariables()
	core.ConnectToDb()
}

func main() {
	r := gin.Default()
	r.GET("/", api.RootApi)
	r.POST("/api/v1/user", api.UserCreateApi)
	r.Run()
}
