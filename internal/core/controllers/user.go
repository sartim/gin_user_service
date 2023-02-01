package controllers

import (
	"gin-shop-api/internal/core/helpers"
	"gin-shop-api/internal/core/models"
	"gin-shop-api/internal/core/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	helpers.DB.Find(&users)
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
	helpers.DB.First(&user, id)
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
	result := helpers.DB.Create(&user)

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
	helpers.DB.First(&user, id)
	helpers.DB.Model(&user).Updates(models.User{
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
	helpers.DB.Delete(&models.User{}, id)
	c.Status(200)
}
