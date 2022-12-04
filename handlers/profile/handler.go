package profile

import (
	"net/http"
	"tink/middlewares"
	"tink/models"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func (h *Handler) Index(c *gin.Context) {
	claims, _ := c.Get("claims")
	claimStrings, _ := claims.(jwt.MapClaims)
	claim, _ := claimStrings["username"].(string)

	var user models.User
	if err := h.db.Where("username = ?", claim).Find(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	if user.Role == models.RoleStudent {
		var student models.Student
		h.db.Where("user_id = ?", user.Id).Find(&student)
		user.Fullname = student.User.Username
	}
	if user.Role == models.RoleTeacher {
		var teacher models.Teacher
		h.db.Where("user_id = ?", user.Id).Find(&teacher)
		user.Fullname = teacher.User.Username
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
func (h *Handler) UpdatePassword(c *gin.Context) {

}

func (h *Handler) Setup(router *gin.RouterGroup) {
	r := router.Group("/profile")

	r.Use(middlewares.Guard())
	r.GET("/", h.Index)
	r.PATCH("/password", h.UpdatePassword)
}

func New(db *gorm.DB) *Handler {
	return &Handler{db}
}
