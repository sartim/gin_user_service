package controllers

import (
	"gin-shop-api/app/core"
	"gin-shop-api/app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserGetAll(c *gin.Context) {
	var users []models.User
	core.DB.Find(&users)

	c.JSON(200, gin.H{
		"data": users,
	})

}

func UserGetByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	core.DB.First(&user, id)
	c.JSON(200, gin.H{
		"data": user,
	})
}

type UserSchema struct {
	FirstName string  `json:"first_name" binding:"required"`
	LastName  string  `json:"last_name" binding:"required"`
	Email     *string `json:"email" binding:"required"`
	Password  string  `json:"password" binding:"required"`
	IsActive  bool    `json:"is_active" binding:"required"`
}

func UserCreate(c *gin.Context) {

	// Validate input
	var input UserSchema
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Bind(&input)

	user := models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  core.HashPassword(input.Password),
		IsActive:  input.IsActive,
	}
	result := core.DB.Create(&user)

	if result.Error != nil {
		// c.Status(400)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Record was not saved"})
		return
	}
}

func UserUpdate(c *gin.Context) {
	id := c.Param("id")
	var input UserSchema
	c.Bind(&input)
	var user models.User
	core.DB.First(&user, id)
	core.DB.Model(&user).Updates(models.User{
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
	id := c.Param("id")
	core.DB.Delete(&models.User{}, id)
	c.Status(200)
}
