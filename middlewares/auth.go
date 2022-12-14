package middlewares

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, _ := c.Cookie("token")

		if tokenString == "" {
			c.Next()
			return
		}

		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		claims, ok := token.Claims.(jwt.MapClaims)

		if !(ok && token.Valid) {
			c.Next()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func Guard(args ...interface{}) gin.HandlerFunc {
	bypass := false
	for _, arg := range args {
		switch t := arg.(type) {
		case bool:
			bypass = t
		default:
			panic("Unknown argument")
		}
	}

	return func(c *gin.Context) {
		claims, ok := c.Get("claims")
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claimStrings, ok := claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		username, ok := claimStrings["username"].(string)
		if !ok || username == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		secure, ok := claimStrings["secure"].(bool)
		if !bypass && (!ok || !secure) {
			c.AbortWithStatusJSON(http.StatusTooEarly, gin.H{
				"message": "Please set your password first",
			})
			return
		}
		c.Next()
	}
}
func Gate(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		claimStrings, _ := claims.(jwt.MapClaims)
		claim, _ := claimStrings["role"].(string)

		if claim != role {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
