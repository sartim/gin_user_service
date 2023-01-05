package api

import (
	"go-shop-admin/app/core"
	"go-shop-admin/app/models"

	"github.com/gin-gonic/gin"
)

func UserCreateApi(c *gin.Context) {
	var data struct {
		FirstName string
		LastName  string
	}

	c.Bind(&data)

	user := models.User{FirstName: data.FirstName, LastName: data.LastName}
	result := core.DB.Create(&user)

	if result.Error != nil {
		c.Status(400)
		return
	}
}
