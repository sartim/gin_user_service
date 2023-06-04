package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	passwordToByte := []byte(password)
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordToByte, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}
