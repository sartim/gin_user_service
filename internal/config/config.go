package config

import "os"

func DatabaseConfig() string {
	return os.Getenv("DB_URL")
}
