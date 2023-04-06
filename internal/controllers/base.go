package controllers

import (
	"context"
	"fmt"
	"gin-shop-api/internal/models"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BaseController struct {
	db     *gorm.DB
	model  interface{}
	schema interface{}
}

type PaginationParams struct {
	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=100"`
}

func NewBaseController(db *gorm.DB, model interface{}, schema interface{}) *BaseController {
	return &BaseController{db, model, schema}
}

func (ctrl *BaseController) GetAll(c *gin.Context) {
	var params PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// check if page query parameter is provided and parse it
	if pageParam := c.Query("page"); pageParam != "" {
		parsedPage, err := strconv.Atoi(pageParam)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid page parameter",
			})
			return
		}
		params.Page = parsedPage
	}

	// check if limit query parameter is provided and parse it
	if limitParam := c.Query("limit"); limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid limit parameter",
			})
			return
		}
		params.Limit = parsedLimit
	}

	// use reflection to create a new slice of the correct type
	sliceType := reflect.SliceOf(reflect.TypeOf(ctrl.model))
	records := reflect.New(sliceType).Interface()

	// calculate offset based on page and limit
	offset := (params.Page - 1) * params.Limit

	// pass ctx to database queries
	// use WithContext() method of gorm.DB to pass the context
	ctrl.db = ctrl.db.WithContext(ctx)

	// pass a pointer to the slice to Offset() and Limit() methods
	ctrl.db.Offset(offset).Limit(params.Limit).Find(records)

	// check if records are empty and return 404 if true
	if reflect.ValueOf(records).Elem().Len() == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "No results",
		})
		return
	}

	count := int64(reflect.ValueOf(records).Elem().Len())

	// convert slice of user models to slice of interfaces
	var interfaceSlice []interface{}
	for _, record := range reflect.ValueOf(records).Elem().Interface().([]models.User) {
		interfaceSlice = append(interfaceSlice, record)
	}

	// add full url and count to response
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, c.Request.URL.String())
	response := gin.H{
		"count": count,
		"url":   baseURL,
		"data":  interfaceSlice,
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *BaseController) Get(c *gin.Context) {
	id := c.Param("id")

	// use reflection to create a new slice of the correct type
	sliceType := reflect.SliceOf(reflect.TypeOf(ctrl.model))
	record := reflect.New(sliceType).Interface()

	if err := ctrl.db.First(record, "id = ?", id).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, record)
}

func (ctrl *BaseController) Create(c *gin.Context) {
	model := reflect.New(reflect.TypeOf(ctrl.model)).Interface()

	if err := c.ShouldBindJSON(&model); err != nil {
		log.Printf("%s: %s", "Field validation failed", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	tx := ctrl.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := tx.Create(model).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	tx.Commit()

	c.JSON(http.StatusCreated, model)
}

func (ctrl *BaseController) Update(c *gin.Context) {
	id := c.Param("id")
	model := reflect.New(reflect.TypeOf(ctrl.model)).Interface()

	if err := ctrl.db.First(model, id).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err := c.ShouldBindJSON(&ctrl.model); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctrl.db.Save(&ctrl.model)
	c.JSON(http.StatusOK, &ctrl.model)
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
