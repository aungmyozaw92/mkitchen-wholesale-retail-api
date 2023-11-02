package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/utils/token"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.TokenValid(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
