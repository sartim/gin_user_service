package middleware

import (
	"fmt"
	"gin-shop-api/internal/helpers/logging"
	"gin-shop-api/internal/models"
	"gin-shop-api/internal/repository"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer logging.HandlePanic()

		authorizationValue := c.Request.Header.Get("Authorization")
		authorizationValues := strings.Fields(authorizationValue)
		if strings.HasPrefix(authorizationValue, "Bearer") {
			if len(authorizationValues) > 0 {
				tokenString := authorizationValues[1]
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(os.Getenv("SECRET_KEY")), nil
				})
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{
						"error": "Invalid Authorization Token",
					})
					c.AbortWithStatus(http.StatusUnauthorized)
				}

				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					// check the exp
					if float64(time.Now().Unix()) > claims["exp"].(float64) {
						c.JSON(http.StatusUnauthorized, gin.H{
							"error": "Expired Token",
						})
						c.AbortWithStatus(http.StatusUnauthorized)
					}

					var user = models.User{ID: uuid.Must(uuid.Parse(fmt.Sprint(claims["sub"])))}

					repository.DB.First(&user)

					if user.ID == uuid.Nil {
						c.AbortWithStatus(http.StatusUnauthorized)
					}

					c.Set("user", user)

					c.Next()
				} else {
					c.AbortWithStatus(http.StatusUnauthorized)
				}
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Missing Authorization Header",
				})
				c.AbortWithStatus(http.StatusUnauthorized)
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Bearer Token Required",
			})
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	}
}
