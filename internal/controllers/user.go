package controllers

import (
	"gin-shop-api/internal/helpers/crypto"
	"gin-shop-api/internal/helpers/validation"
	"gin-shop-api/internal/models"
	"gin-shop-api/internal/schemas"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	*BaseController
}

func NewUserController(db *gorm.DB) *UserController {
	var user models.User
	var schema schemas.UserSchema
	return &UserController{NewBaseController(db, user, schema)}
}

func (ctrl *UserController) Create(c *gin.Context) {
	var input schemas.UserSchema

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("%s: %s", "Field validation failed", err)
		validation.ValidateSchema(c, err, "body")
		return
	}

	// Set the hashed password in the user model
	user := models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  crypto.HashPassword(input.Password),
		IsActive:  false,
	}

	// Save the user to the database
	result := ctrl.db.Create(&user)

	if result.Error != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"user": user})
}

func (ctrl *UserController) RegisterUserRoutes(router *gin.RouterGroup) {
	userRouter := router.Group("/users")
	{
		userRouter.GET("", ctrl.GetAll)
		userRouter.GET(":id", ctrl.Get)
		userRouter.POST("", ctrl.Create)
		userRouter.PUT(":id", ctrl.Update)
		userRouter.DELETE(":id", ctrl.Delete)
	}
}
