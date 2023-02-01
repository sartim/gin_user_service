package routes

import (
	"gin-shop-api/app/controllers"
	"gin-shop-api/app/core"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/api/v1/generate-jwt", controllers.GenerateJWT)
}

func UserRoutes(r *gin.Engine) {
	r.GET("/api/v1/user", core.RequireAuth, controllers.UserGetAll)
	r.GET("/api/v1/user/:id", core.RequireAuth, controllers.UserGetByID)
	r.POST("/api/v1/user", core.RequireAuth, controllers.UserCreate)
	r.PUT("/api/v1/user/:id", core.RequireAuth, controllers.UserUpdate)
	r.DELETE("/api/v1/user/:id", core.RequireAuth, controllers.UserDelete)
}
