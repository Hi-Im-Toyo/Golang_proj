package middleware

import (
	"net/http"

	token "github.com/Hi-Im-Toyo/GO_Proj/tokens"
	"github.com/gin-gonic/gin"
)

func Authhentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientToken := c.Request.Header.Get("Authorization")
		if ClientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Token Found"})
			c.Abort()
			return

		}
		claims, err := token.ValidateToken(ClientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
