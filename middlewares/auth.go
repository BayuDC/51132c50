package middlewares

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
)

type User struct {
	Id       int    `json:"id"`
	Userable int    `json:"userable"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Secure   bool   `json:"secure"`
	Role     string `json:"role"`
}

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

		claims, _ := token.Claims.(jwt.MapClaims)
		user := User{
			Id:       int(claims["id"].(float64)),
			Userable: int(claims["userable"].(float64)),
			Username: claims["username"].(string),
			Fullname: claims["fullname"].(string),
			Secure:   claims["secure"].(bool),
			Role:     claims["role"].(string),
		}

		c.Set("user", &user)
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
		userRaw, ok := c.Get("user")
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		user := userRaw.(*User)

		if user.Username == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !bypass && (!user.Secure) {
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
		user := (c.MustGet("user")).(*User)

		if user.Role != role {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
