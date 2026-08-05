package controllers

import (
	"errors"
	"net/http"
	"strings"

	"gin-shop-api/internal/helpers/crypto"
	"gin-shop-api/internal/models"
	"gin-shop-api/internal/schemas"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	*BaseController[models.User]
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{NewBaseController[models.User](db, map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"is_active":  "is_active",
		"is_admin":   "is_admin",
	})}
}

func (ctrl *UserController) Create(c *gin.Context) {
	var input schemas.UserSchema
	if err := c.ShouldBindJSON(&input); err != nil {
		writeValidationError(c, err)
		return
	}

	hashedPassword, err := crypto.HashPassword(input.Password)
	if err != nil {
		internalError(c)
		return
	}
	user := models.User{
		FirstName: strings.TrimSpace(input.FirstName),
		LastName:  strings.TrimSpace(input.LastName),
		Email:     strings.ToLower(strings.TrimSpace(input.Email)),
		Password:  hashedPassword,
		IsActive:  input.IsActive,
		IsAdmin:   input.IsAdmin,
	}
	if err := ctrl.db.WithContext(c.Request.Context()).Create(&user).Error; err != nil {
		handleWriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func (ctrl *UserController) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	users.GET("", ctrl.GetAll)
	users.GET("/:id", ctrl.Get)
	users.POST("", ctrl.Create)
	users.PATCH("/:id", ctrl.Update)
	users.DELETE("/:id", ctrl.Delete)
}

func isNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
