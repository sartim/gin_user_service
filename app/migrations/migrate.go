package main

import (
	"go-shop-api/app/core"
	"go-shop-api/app/models"
)

func init() {
	core.LoadEnvVariables()
	core.ConnectToDb()
}

func main() {
	core.DB.AutoMigrate(&models.User{})
}
