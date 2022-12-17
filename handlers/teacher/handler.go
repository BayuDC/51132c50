package teacher

import (
	"errors"
	"net/http"
	"tink/middlewares"
	"tink/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func (h *Handler) Load(c *gin.Context) {
	var id = c.Param("id")
	var teacher models.Teacher

	if err := h.db.Preload("User").First(&teacher, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Teacher not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	teacher.Username = teacher.User.Username

	c.Set("teacher", &teacher)
	c.Next()
}

func (h *Handler) Index(c *gin.Context) {
	var teachers []models.Teacher
	h.db.Table("teachers").
		Select("teachers.id, teachers.fullname, users.username").
		Joins("left join users on users.id = user_id").
		Find(&teachers)
	c.JSON(http.StatusOK, gin.H{
		"teachers": teachers,
	})
}
func (h *Handler) Show(c *gin.Context) {
	teacher := c.MustGet("teacher")
	c.JSON(http.StatusOK, gin.H{"teacher": teacher})
}
func (h *Handler) Store(c *gin.Context) {
	var body CreateTeacherSchema

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacher := models.Teacher{
		Fullname: body.Fullname,
		User: models.User{
			Username: body.Username,
			Role:     models.RoleTeacher,
		},
	}
	if err := h.db.Create(&teacher).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	teacher.Username = teacher.User.Username
	c.JSON(http.StatusCreated, gin.H{"teacher": teacher})
}
func (h *Handler) Update(c *gin.Context) {
	var body UpdateTeacherSchema

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacher := (c.MustGet("teacher")).(*models.Teacher)

	if body.Fullname != nil {
		teacher.Fullname = *body.Fullname
	}
	if err := h.db.Omit("User").Save(&teacher).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacher.Username = teacher.User.Username
	c.JSON(http.StatusOK, gin.H{"teacher": teacher})
}

func (h *Handler) Destroy(c *gin.Context) {
	teacher := (c.MustGet("teacher")).(*models.Teacher)

	if err := h.db.Delete(&teacher.User).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Setup(r *gin.RouterGroup) {
	router := r.Group("")
	router.Use(middlewares.Guard())
	router.GET("/teachers", h.Index)
	router.GET("/teachers/:id", h.Load, h.Show)
	router.Use(middlewares.Gate("admin"))
	router.POST("/teachers", h.Store)
	router.PUT("/teachers/:id", h.Load, h.Update)
	router.DELETE("/teachers/:id", h.Load, h.Destroy)
}

func New(db *gorm.DB) *Handler {
	return &Handler{db}
}
