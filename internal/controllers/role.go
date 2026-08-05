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

type RoleController struct {
	*BaseController[models.Role]
}

func NewRoleController(db *gorm.DB) *RoleController {
	return &RoleController{NewBaseController[models.Role](db, map[string]string{
		"name": "name", "description": "description",
	})}
}

func (ctrl *RoleController) Create(c *gin.Context) {
	var input schemas.RoleSchema

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("%s: %s", "Field validation failed", err)
		errors := validation.ValidateSchema(err, "body")
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Set the hashed password in the role model
	role := models.Role{
		Name:        input.Name,
		Description: input.Description,
	}

	// Save the role to the database
	result := ctrl.db.Create(&role)

	if result.Error != nil {
		handleWriteError(c, result.Error)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": role})
}

func (ctrl *RoleController) RegisterRoleRoutes(
	router *gin.RouterGroup) {
	roleRouter := router.Group("/roles")
	{
		roleRouter.GET("", ctrl.GetAll)
		roleRouter.GET("/:id", ctrl.Get)
		roleRouter.POST("", ctrl.Create)
		roleRouter.PATCH("/:id", ctrl.Update)
		roleRouter.DELETE("/:id", ctrl.Delete)
	}
}
