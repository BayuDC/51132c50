package auth

import (
	"errors"
	"net/http"
	"os"
	"time"
	"tink/middlewares"
	"tink/models"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func (h *Handler) Index(c *gin.Context) {
	claims, _ := c.Get("claims")
	claimStrings, _ := claims.(jwt.MapClaims)

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"username": claimStrings["username"].(string),
			"fullname": claimStrings["fullname"].(string),
			"secure":   claimStrings["secure"].(bool),
			"role":     claimStrings["role"].(string),
		},
	})
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
	switch user.Role {
	case models.RoleStudent:
		var student models.Student
		h.db.Where("user_id = ?", user.Id).Find(&student)
		user.Fullname = student.Fullname
	case models.RoleTeacher:
		var teacher models.Teacher
		h.db.Where("user_id = ?", user.Id).Find(&teacher)
		user.Fullname = teacher.Fullname
	}

	defaultPassword := false
	if user.Password == "" && body.Password == "" {
		defaultPassword = true
	} else if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Password incorrect",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"fullname": user.Fullname,
		"username": user.Username,
		"secure":   !defaultPassword,
		"role":     string(user.Role),
		"exp":      time.Now().Add(time.Hour * 24 * 3).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: true,
		MaxAge:   60 * 60 * 24 * 3,
		Path:     "/",
	})
	c.JSON(http.StatusNoContent, nil)
}
func (h *Handler) Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "token",
		Path:   "/",
		MaxAge: -1,
	})
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) Setup(r *gin.RouterGroup) {
	r.GET("/auth", middlewares.Guard(), h.Index)
	r.POST("/auth/login", h.Login)
	r.POST("/auth/logout", middlewares.Guard(true), h.Logout)
}
func New(db *gorm.DB) *Handler {
	return &Handler{db}
}
