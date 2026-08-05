package controllers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

type BaseController[T any] struct {
	db           *gorm.DB
	updateFields map[string]string
}

func NewBaseController[T any](db *gorm.DB, updateFields map[string]string) *BaseController[T] {
	return &BaseController[T]{db: db, updateFields: updateFields}
}

func (ctrl *BaseController[T]) GetAll(c *gin.Context) {
	page, ok := positiveQueryInt(c, "page", 1, 0)
	if !ok {
		return
	}
	limit, ok := positiveQueryInt(c, "limit", defaultPageSize, maxPageSize)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	db := ctrl.db.WithContext(ctx)

	var total int64
	if err := db.Model(new(T)).Count(&total).Error; err != nil {
		internalError(c)
		return
	}

	records := make([]T, 0)
	if err := db.Offset((page - 1) * limit).Limit(limit).Find(&records).Error; err != nil {
		internalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       records,
		"pagination": gin.H{"page": page, "limit": limit, "total": total},
	})
}

func (ctrl *BaseController[T]) Get(c *gin.Context) {
	var record T
	err := ctrl.db.WithContext(c.Request.Context()).First(&record, "id = ?", c.Param("id")).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
		return
	}
	if err != nil {
		internalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": record})
}

func (ctrl *BaseController[T]) Update(c *gin.Context) {
	var input map[string]any
	if err := c.ShouldBindJSON(&input); err != nil || len(input) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	updates := make(map[string]any, len(input))
	for key, value := range input {
		column, allowed := ctrl.updateFields[key]
		if !allowed {
			c.JSON(http.StatusBadRequest, gin.H{"error": "field cannot be updated: " + key})
			return
		}
		updates[column] = value
	}

	var record T
	result := ctrl.db.WithContext(c.Request.Context()).Model(&record).
		Where("id = ?", c.Param("id")).Updates(updates)
	if result.Error != nil {
		handleWriteError(c, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
		return
	}
	ctrl.Get(c)
}

func (ctrl *BaseController[T]) Delete(c *gin.Context) {
	var record T
	result := ctrl.db.WithContext(c.Request.Context()).Where("id = ?", c.Param("id")).Delete(&record)
	if result.Error != nil {
		internalError(c)
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func positiveQueryInt(c *gin.Context, key string, fallback, maximum int) (int, bool) {
	raw := c.Query(key)
	if raw == "" {
		return fallback, true
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 || maximum > 0 && value > maximum {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + key + " parameter"})
		return 0, false
	}
	return value, true
}

func handleWriteError(c *gin.Context, err error) {
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23505" {
		c.JSON(http.StatusConflict, gin.H{"error": "resource already exists"})
		return
	}
	internalError(c)
}

func internalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
