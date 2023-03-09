package controllers

import (
	"gin-shop-api/internal/helpers"
	"gin-shop-api/internal/models"
	"gin-shop-api/internal/repository"
	"gin-shop-api/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	*BaseController
}

func NewUserController(db *gorm.DB) *UserController {
	var user models.User
	return &UserController{NewBaseController(db, user)}
}

func (ctrl *UserController) RegisterRoutes(router *gin.RouterGroup) {
	userRouter := router.Group("/users")
	{
		userRouter.GET("", ctrl.GetAll)
		userRouter.GET(":id", ctrl.Get)
		userRouter.POST("", ctrl.Create)
		userRouter.PUT(":id", ctrl.Update)
		userRouter.DELETE(":id", ctrl.Delete)
	}
}

func UserGetAll(c *gin.Context) {
	// Validate headers
	var header schemas.HeaderSchema
	if err := c.ShouldBindHeader(&header); err != nil {
		helpers.LogError.Printf("%s: %s", "Missing header", err)
		helpers.ValidateSchema(c, err, "header")
		return
	}

	// Validate Content-Type
	contentType := c.Request.Header.Get("Content-Type")
	if contentType != "application/json" {
		c.JSON(400, gin.H{
			"errors": []helpers.Dict{{"Content-Type": "Not application/json"}},
		})
		return
	}

	var users []models.User
	repository.DB.Find(&users)
	c.JSON(200, gin.H{
		"data": users,
	})

}

func UserGetByID(c *gin.Context) {
	// Validate headers
	var header schemas.HeaderSchema
	if err := c.ShouldBindHeader(&header); err != nil {
		helpers.LogError.Printf("%s: %s", "Missing header", err)
		helpers.ValidateSchema(c, err, "header")
		return
	}

	// Validate Content-Type
	contentType := c.Request.Header.Get("Content-Type")
	if contentType != "application/json" {
		c.JSON(400, gin.H{
			"errors": []helpers.Dict{{"Content-Type": "Not application/json"}},
		})
		return
	}

	id := c.Param("id")
	var user models.User
	repository.DB.First(&user, id)
	c.JSON(200, gin.H{
		"data": user,
	})
}

func UserCreate(c *gin.Context) {
	// Validate headers
	var header schemas.HeaderSchema
	if err := c.ShouldBindHeader(&header); err != nil {
		helpers.LogError.Printf("%s: %s", "Missing header", err)
		helpers.ValidateSchema(c, err, "header")
		return
	}

	// Validate Content-Type
	contentType := c.Request.Header.Get("Content-Type")
	if contentType != "application/json" {
		c.JSON(400, gin.H{
			"errors": []helpers.Dict{{"Content-Type": "Not application/json"}},
		})
		return
	}

	// Validate input
	var input schemas.UserSchema
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Bind(&input)

	id := helpers.GenerateUUID()

	user := models.User{
		ID:        id,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  helpers.HashPassword(input.Password),
		IsActive:  input.IsActive,
		Deleted:   false,
	}
	result := repository.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Record was not saved"})
		return
	}
}

func UserUpdate(c *gin.Context) {
	// Validate headers
	var header schemas.HeaderSchema
	if err := c.ShouldBindHeader(&header); err != nil {
		helpers.LogError.Printf("%s: %s", "Missing header", err)
		helpers.ValidateSchema(c, err, "header")
		return
	}

	// Validate Content-Type
	contentType := c.Request.Header.Get("Content-Type")
	if contentType != "application/json" {
		c.JSON(400, gin.H{
			"errors": []helpers.Dict{{"Content-Type": "Not application/json"}},
		})
		return
	}

	id := c.Param("id")
	var input schemas.UserSchema
	c.Bind(&input)
	var user models.User
	repository.DB.First(&user, id)
	repository.DB.Model(&user).Updates(models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		IsActive:  input.IsActive,
	})
	c.JSON(200, gin.H{
		"data": user,
	})
}

func UserDelete(c *gin.Context) {
	// Validate headers
	var header schemas.HeaderSchema
	if err := c.ShouldBindHeader(&header); err != nil {
		helpers.LogError.Printf("%s: %s", "Missing header", err)
		helpers.ValidateSchema(c, err, "header")
		return
	}

	// Validate Content-Type
	contentType := c.Request.Header.Get("Content-Type")
	if contentType != "application/json" {
		c.JSON(400, gin.H{
			"errors": []helpers.Dict{{"Content-Type": "Not application/json"}},
		})
		return
	}

	id := c.Param("id")
	repository.DB.Delete(&models.User{}, id)
	c.Status(200)
}
