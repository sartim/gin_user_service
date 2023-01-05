package api

import (
	"go-shop-api/app/core"
	"go-shop-api/app/models"

	"github.com/gin-gonic/gin"
)

func UserGetApi(c *gin.Context) {

}

func UserCreateApi(c *gin.Context) {
	var data struct {
		FirstName string
		LastName  string
		Email     string
		Password  string
		IsActive  bool
	}

	c.Bind(&data)

	user := models.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  core.HashPassword(data.Password),
		IsActive:  data.IsActive,
	}
	result := core.DB.Create(&user)

	if result.Error != nil {
		c.Status(400)
		return
	}
}

func UserUpdateApi(c *gin.Context) {

}

func UserDeleteApi(c *gin.Context) {

}
