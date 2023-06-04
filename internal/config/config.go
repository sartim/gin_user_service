package config

import (
	"fmt"
	"os"
)

const Environment = "ENV"
const Port = "PORT"
const DbUrl = "DB_URL"

func DatabaseConfig() string {
	return os.Getenv("DB_URL")
}

type Config struct {
	EnvVar string
}

func (config Config) Get() string {
	envVar := os.Getenv(config.EnvVar)
	return envVar
}

func CurrDir() string {
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	return mydir
}
