package core

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	password_to_byte := []byte(password)
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password_to_byte, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}
