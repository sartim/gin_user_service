package controllers

import (
	"errors"
	"net/http"
	"time"

	"gin-shop-api/internal/helpers/validation"
	"gin-shop-api/internal/models"
	"gin-shop-api/internal/schemas"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const tokenIssuer = "gin-user-service"

type AuthController struct {
	db       *gorm.DB
	secret   []byte
	tokenTTL time.Duration
}

func NewAuthController(db *gorm.DB, secret string, tokenTTL time.Duration) *AuthController {
	return &AuthController{db: db, secret: []byte(secret), tokenTTL: tokenTTL}
}

func (ctrl *AuthController) GenerateJWT(c *gin.Context) {
	var input schemas.AuthSchema
	if err := c.ShouldBindJSON(&input); err != nil {
		writeValidationError(c, err)
		return
	}

	var user models.User
	err := ctrl.db.WithContext(c.Request.Context()).First(&user, "email = ?", input.Email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) || err == nil && !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	if err != nil {
		internalError(c)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		Issuer:    tokenIssuer,
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ctrl.tokenTTL)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(ctrl.secret)
	if err != nil {
		internalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": tokenString,
		"token_type":   "Bearer",
		"expires_in":   int64(ctrl.tokenTTL.Seconds()),
		"user":         user,
	})
}

func (ctrl *AuthController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/auth/token", ctrl.GenerateJWT)
}

func writeValidationError(c *gin.Context, err error) {
	errors := validation.ValidateSchema(err, "body")
	if errors == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
}
