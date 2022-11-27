package auth

import (
	"errors"
	"net/http"
	"time"
	"tink/models"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func (h *Handler) Login(c *gin.Context) {
	var body LoginSchema
	var user models.User

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.db.Where("username = ?", body.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "User not found",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Password incorrect",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString([]byte("mysecretkey"))

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
func (h *Handler) Logout(c *gin.Context) {

}

func (h *Handler) Setup(router *gin.RouterGroup) {
	r := router.Group("/auth")

	r.POST("/login", h.Login)
	r.POST("/logout", h.Logout)
}
func New(db *gorm.DB) *Handler {
	return &Handler{db}
}
