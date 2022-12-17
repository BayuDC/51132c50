package course

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
	var course models.Course

	if err := h.db.Preload("Teacher").First(&course, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Course not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Set("course", course)
	c.Next()
}

func (h *Handler) Index(c *gin.Context) {
	var courses []models.Course

	if err := h.db.Preload("Teacher").Order("id ASC").Find(&courses).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"courses": courses})
}

func (h *Handler) Show(c *gin.Context) {
	course := c.MustGet("course")
	c.JSON(http.StatusOK, gin.H{"course": course})
}

func (h *Handler) Store(c *gin.Context) {
	var body CreateCourseSchema

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course := models.Course{
		Name:      body.Name,
		TeacherId: body.Teacher,
	}

	if err := h.db.Create(&course).Error; err != nil {
		if errors.Is(err, models.CourseTeacherNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"group": course})
}

func (h *Handler) Update(c *gin.Context) {
	var body UpdateCourseSchema

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course := (c.MustGet("course")).(models.Course)

	if body.Name != nil {
		course.Name = *body.Name
	}

	if err := h.db.Save(&course).Error; err != nil {
		if errors.Is(err, models.CourseTeacherNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"course": course})
}
func (h *Handler) Destroy(c *gin.Context) {
	course := (c.MustGet("course")).(models.Course)

	if err := h.db.Delete(&course).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Setup(r *gin.RouterGroup) {
	router := r.Group("")
	router.Use(middlewares.Guard())
	router.GET("/courses", h.Index)
	router.GET("/courses/:id", h.Load, h.Show)
	router.Use(middlewares.Gate("admin"))
	router.POST("/courses", h.Store)
	router.PUT("/courses/:id", h.Load, h.Update)
	router.DELETE("/courses/:id", h.Load, h.Destroy)
}

func New(db *gorm.DB) *Handler {
	return &Handler{db}
}
