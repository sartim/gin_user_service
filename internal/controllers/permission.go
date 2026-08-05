package controllers

import (
	"gin-shop-api/internal/helpers/validation"
	"gin-shop-api/internal/models"
	"gin-shop-api/internal/schemas"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PermissionController struct {
	*BaseController[models.Permission]
}

func NewPermissionController(db *gorm.DB) *PermissionController {
	return &PermissionController{NewBaseController[models.Permission](db, map[string]string{
		"name": "name", "description": "description",
	})}
}

func (ctrl *PermissionController) Create(c *gin.Context) {
	var input schemas.PermissionSchema

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("%s: %s", "Field validation failed", err)
		errors := validation.ValidateSchema(err, "body")
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Set the hashed password in the permission model
	permission := models.Permission{
		Name:        input.Name,
		Description: input.Description,
	}

	// Save the permission to the database
	result := ctrl.db.Create(&permission)

	if result.Error != nil {
		handleWriteError(c, result.Error)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": permission})
}

func (ctrl *PermissionController) RegisterPermissionRoutes(
	router *gin.RouterGroup) {
	permissionRouter := router.Group("/permissions")
	{
		permissionRouter.GET("", ctrl.GetAll)
		permissionRouter.GET("/:id", ctrl.Get)
		permissionRouter.POST("", ctrl.Create)
		permissionRouter.PATCH("/:id", ctrl.Update)
		permissionRouter.DELETE("/:id", ctrl.Delete)
	}
}
