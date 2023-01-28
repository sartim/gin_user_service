package core

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func ConnectToDb() {
	var err error
	dsn := DatabaseConfig()
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	fmt.Println(DB)

	if err != nil {
		LogError.Printf("%s: %s", "Failed to connect to the datbase", err)
	}
}
