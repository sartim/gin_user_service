package main

import (
	"gin-shop-api/app/core"
	"gin-shop-api/app/models"
)

func init() {
	core.LoadEnvVariables()
	core.ConnectToDb()
}

func main() {
	core.DB.AutoMigrate(&models.User{})
}
