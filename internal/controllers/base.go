package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BaseController struct {
	db    *gorm.DB
	model any
}

func NewBaseController(db *gorm.DB, model any) *BaseController {
	return &BaseController{db, model}
}

func (ctrl *BaseController) GetAll(c *gin.Context) {
	var records []interface{}
	ctrl.db.Find(&records)
	c.JSON(http.StatusOK, records)
}

func (ctrl *BaseController) Get(c *gin.Context) {
	id := c.Param("id")
	var record interface{}
	if err := ctrl.db.First(&record, id).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, record)
}

func (ctrl *BaseController) Create(c *gin.Context) {
	if err := c.ShouldBindJSON(&ctrl.model); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctrl.db.Create(&ctrl.model)
	c.JSON(http.StatusCreated, ctrl.model)
}

func (ctrl *BaseController) Update(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.db.First(&ctrl.model, id).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err := c.ShouldBindJSON(&ctrl.model); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctrl.db.Save(&ctrl.model)
	c.JSON(http.StatusOK, ctrl.model)
}

func (ctrl *BaseController) Delete(c *gin.Context) {
	id := c.Param("id")
	var record interface{}
	if err := ctrl.db.First(&record, id).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	ctrl.db.Delete(&record)
	c.Status(http.StatusNoContent)
}
