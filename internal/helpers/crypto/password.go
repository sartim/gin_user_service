package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passwordToByte := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordToByte, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
