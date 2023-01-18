package main

import (
	"flag"
	"fmt"
	"gin-shop-api/app/core"
	"gin-shop-api/app/models"
)

var cliType string

func init() {
	core.LoadEnvVariables()
	core.ConnectToDb()
}

func main() {
	flag.StringVar(&cliType, "cli_type", "", "CLI Type e.g. migrate, drop_all")
	flag.Parse()

	if cliType == "migrate" {
		core.DB.AutoMigrate(&models.User{})
		core.DB.AutoMigrate(&models.Status{})
	}
	if cliType == "drop_all" {
		core.DB.Migrator().DropTable(&models.User{})
		core.DB.Migrator().DropTable(&models.Status{})
		fmt.Println("Finished dropping tables")
	}
}
