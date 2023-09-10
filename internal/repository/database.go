package repository

import (
	"fmt"
	"gin-shop-api/internal/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func ConnectToDb() {
	var err error
	dsn := config.DatabaseConfig()
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	fmt.Println(DB)

	if err != nil {
		log.Printf("%s: %s", "Failed to connect to the database", err)
	}
}
