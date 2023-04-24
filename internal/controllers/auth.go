package controllers

import (
	"gin-shop-api/internal/helpers"
	"gin-shop-api/internal/models"
	"gin-shop-api/internal/repository"
	"gin-shop-api/internal/schemas"
	"log"
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
		log.Printf("%s: %s", "Field validation failed", err)
		helpers.ValidateSchema(c, err, "body")
		return
	}

	// Lookup user
	var user models.User
	repository.DB.First(&user, "email = ?", input.Email)

	if user.ID == uuid.Nil {
		log.Printf("%s", "Email does not exist")
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
		log.Printf("%s: %s", "Password does not match", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // Expired in 30 days
	})

	// Sign and get encoded string
	var sampleSecretKey = []byte(os.Getenv("SECRET_KEY"))
	tokenString, err := token.SignedString(sampleSecretKey)

	if err != nil {
		log.Printf("%s: %s", "Failed to create token", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user":  &user,
	})
}

func RegisterAuthRoutes(router *gin.RouterGroup) {
	authRouter := router.Group("/auth")
	{
		authRouter.POST("/generate-jwt", GenerateJWT)
	}
}
