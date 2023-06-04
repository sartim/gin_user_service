package env

import (
	"gin-shop-api/internal/helpers/logging"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		msg := "Error loading .env file"
		log.Panicf("%s: %s", msg, err)
	}
}

func GetEnv(envVar string) string {
	res, exists := os.LookupEnv(envVar)
	if exists {
		return res
	} else {
		logging.LogError.Printf("%s does not exist", envVar)
		return ""
	}
}
