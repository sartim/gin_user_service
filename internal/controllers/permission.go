package controllers

import (
	"gin-shop-api/internal/helpers/validation"
	"gin-shop-api/internal/models"
	"gin-shop-api/internal/schemas"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PermissionController struct {
	*BaseController
}

func NewPermissionController(db *gorm.DB) *PermissionController {
	var permission models.Permission
	var schema schemas.PermissionSchema
	return &PermissionController{NewBaseController(db, permission, schema)}
}

func (ctrl *PermissionController) Create(c *gin.Context) {
	var input schemas.PermissionSchema

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("%s: %s", "Field validation failed", err)
		validation.ValidateSchema(c, err, "body")
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
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"permission": permission})
}

func (ctrl *PermissionController) RegisterPermissionRoutes(
	router *gin.RouterGroup) {
	permissionRouter := router.Group("/permissions")
	{
		permissionRouter.GET("", ctrl.GetAll)
		permissionRouter.GET(":id", ctrl.Get)
		permissionRouter.POST("", ctrl.Create)
		permissionRouter.PUT(":id", ctrl.Update)
		permissionRouter.DELETE(":id", ctrl.Delete)
	}
}
