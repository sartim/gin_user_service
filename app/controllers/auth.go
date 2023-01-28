package controllers

import (
	"fmt"
	"gin-shop-api/app/core"
	"gin-shop-api/app/models"
	"gin-shop-api/app/schemas"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateJWT(c *gin.Context) {
	var input schemas.AuthSchema

	// Validate fields
	if err := c.ShouldBindJSON(&input); err != nil {
		core.LogError.Printf("%s: %s", "Field validation failed", err)
		core.ValidateSchema(c, err)
		return
	}

	// Lookup user
	var user models.User
	core.DB.First(&user, "email = ?", input.Email)

	fmt.Println(input.Email)
	if user.ID == uuid.Nil {
		core.LogError.Printf("%s", "Email does not exist")
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
		core.LogError.Printf("%s: %s", "Password does not match", err)
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
	var sampleSecretKey = []byte(os.Getenv("SECRET_KEY"))
	tokenString, err := token.SignedString(sampleSecretKey)

	if err != nil {
		core.LogError.Printf("%s: %s", "Failed to create token", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
