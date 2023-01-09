package routes

import (
	"go-shop-api/app/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/api/v1/generate-jwt", controllers.GenerateJWT)
}

func UserRoutes(r *gin.Engine) {
	r.GET("/api/v1/user", controllers.UserGetAll)
	r.GET("/api/v1/user/:id", controllers.UserGetByID)
	r.POST("/api/v1/user", controllers.UserCreate)
	r.PUT("/api/v1/user/:id", controllers.UserUpdate)
	r.DELETE("/api/v1/user/:id", controllers.UserDelete)
}
