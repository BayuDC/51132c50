package student

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
	var student models.Student

	if err := h.db.Preload("User").First(&student, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Student not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	student.Username = student.User.Username

	c.Set("student", &student)
	c.Next()
}
func (h *Handler) Index(c *gin.Context) {
	var students []models.Student
	if err := h.db.Table("students").
		Select("students.id, students.fullname, users.username").
		Joins("left join users on users.id = user_id").
		Find(&students).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"students": students})
}
func (h *Handler) Show(c *gin.Context) {
	student := c.MustGet("student").(*models.Student)
	c.JSON(http.StatusOK, gin.H{"student": student})
}

func (h *Handler) Store(c *gin.Context) {
	var body CreateStudentSchema

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student := models.Student{
		Fullname: body.Fullname,
		User: models.User{
			Username: body.Username,
			Role:     models.RoleStudent,
		},
	}
	if err := h.db.Create(&student).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	student.Username = student.User.Username
	c.JSON(http.StatusCreated, gin.H{"student": student})
}

func (h *Handler) Update(c *gin.Context) {
	var body UpdateStudentSchema

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student := c.MustGet("student").(*models.Student)

	if body.Fullname != nil {
		student.Fullname = *body.Fullname
	}
	if err := h.db.Omit("User").Save(&student).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student.Username = student.User.Username
	c.JSON(http.StatusOK, gin.H{"student": student})
}

func (h *Handler) Destroy(c *gin.Context) {
	student := c.MustGet("student").(*models.Student)

	if err := h.db.Delete(&student.User).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Setup(r *gin.RouterGroup) {
	router := r.Group("")
	router.Use(middlewares.Guard())
	router.GET("/students", h.Index)
	router.GET("/students/:id", h.Load, h.Show)
	router.Use(middlewares.Gate("admin"))
	router.POST("/students", h.Store)
	router.PUT("/students/:id", h.Load, h.Update)
	router.DELETE("/students/:id", h.Load, h.Destroy)
}

func New(db *gorm.DB) *Handler {
	return &Handler{db}
}
