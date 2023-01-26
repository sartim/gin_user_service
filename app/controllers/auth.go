package controllers

import (
	"errors"
	"fmt"
	"gin-shop-api/app/core"
	"gin-shop-api/app/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var LogError = core.Log("ERROR")

type AuthSchema struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ApiError struct {
	Field string
	Msg   string
}

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	}
	return ""
}

func GenerateJWT(c *gin.Context) {
	var input AuthSchema

	// Validate fields
	if err := c.ShouldBindJSON(&input); err != nil {
		LogError.Printf("%s: %s", "Field validation failed", err)

		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ApiError, len(ve))
			for i, fe := range ve {
				out[i] = ApiError{fe.Field(), msgForTag(fe.Tag())}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": out})
		}
		return
	}

	// Lookup user
	var user models.User
	core.DB.First(&user, "email = ?", input.Email)

	fmt.Println(input.Email)
	if user.ID == uuid.Nil {
		LogError.Printf("%s", "Email does not exist")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	// Check password
	hashedPassword := []byte(user.Password)
	password := []byte(input.Password)
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		LogError.Printf("%s: %s", "Password does not match", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	// Getnerate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // Expired in 30 days
	})

	// Sign and get encoded string
	var sampleSecretKey = []byte(os.Getenv("SECRET"))
	tokenString, err := token.SignedString(sampleSecretKey)

	if err != nil {
		LogError.Printf("%s: %s", "Failed to create token", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
