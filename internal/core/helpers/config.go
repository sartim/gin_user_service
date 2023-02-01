package helpers

import "os"

var LogError = Log("ERROR")

func DatabaseConfig() string {
	return os.Getenv("DB_URL")
}
