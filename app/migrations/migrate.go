package main

import (
	"go-shop-admin/app/core"
	"go-shop-admin/app/models"
)

func init() {
	core.LoadEnvVariables()
	core.ConnectToDb()
}

func main() {
	core.DB.AutoMigrate(&models.User{})
}
