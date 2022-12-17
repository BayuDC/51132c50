package profile

import (
	"net/http"
	"tink/middlewares"
	"tink/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func (h *Handler) Index(c *gin.Context) {
	username := (c.MustGet("user")).(*middlewares.User).Username

	var user models.User
	if err := h.db.Where("username = ?", username).Find(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	if user.Role == models.RoleStudent {
		var student models.Student
		h.db.Where("user_id = ?", user.Id).Find(&student)
		user.Fullname = student.Fullname
	}
	if user.Role == models.RoleTeacher {
		var teacher models.Teacher
		h.db.Where("user_id = ?", user.Id).Find(&teacher)
		user.Fullname = teacher.Fullname
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
func (h *Handler) UpdatePassword(c *gin.Context) {
	username := (c.MustGet("user")).(*middlewares.User).Username

	var body UpdatePasswordSchema

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := h.db.Model(&models.User{}).Where("username = ?", username).Update("password", string(passwordHashed)).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "token",
		Path:   "/",
		MaxAge: -1,
	})
	c.JSON(http.StatusOK, gin.H{
		"message": "Please login again with your new password",
	})
}

func (h *Handler) Setup(r *gin.RouterGroup) {
	router := r.Group("")
	router.GET("/profile", middlewares.Guard(), h.Index)
	router.PATCH("/profile/password", middlewares.Guard(true), h.UpdatePassword)
}

func New(db *gorm.DB) *Handler {
	return &Handler{db}
}
