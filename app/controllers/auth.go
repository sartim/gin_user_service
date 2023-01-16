package controllers

import (
	"gin-shop-api/app/core"
	"gin-shop-api/app/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthSchema struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func GenerateJWT(c *gin.Context) {
	var input AuthSchema
	// Validate fields
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if c.Bind(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read",
		})
		return

	}

	// Lookup user
	var user models.User
	core.DB.First(&user, "email = ?", input.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read",
		})
		return

	}
	// Check password
	hashed_password := []byte(user.Password)
	password := []byte(input.Password)
	err := bcrypt.CompareHashAndPassword(hashed_password, password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid email or password",
		})
		return
	}
	// Getnerate token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // Expired in 30 days
	})
	// Sign and get encoded string
	tokenString, err := token.SignedString(os.Getenv("SECRET"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create token",
		})
		return
	}
	// Send it back
	// c.SetSameSite(http.SameSiteLaxMode)
	// c.SetCookie("Authorization", tokenString, 3600 * 24 * 30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
