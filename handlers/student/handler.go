package student

import (
	"errors"
	"net/http"
	"tink/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func (h *Handler) Index(c *gin.Context) {
	var students []models.Student
	h.db.Find(&students)
	c.JSON(http.StatusOK, gin.H{
		"students": students,
	})
}
func (h *Handler) Show(c *gin.Context) {
	var student models.Student
	id := c.Param("id")

	if err := h.db.First(&student, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Student not found",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"student": student,
	})
}
func (h *Handler) Store(c *gin.Context) {
	var body CreateStudentSchema

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	student := models.Student{
		Name: body.Name,
	}
	h.db.Create(&student)
	c.JSON(http.StatusCreated, gin.H{
		"student": student,
	})
}
func (h *Handler) Update(c *gin.Context) {
	var body UpdateStudentSchema
	var student models.Student
	id := c.Param("id")

	if err := h.db.First(&student, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Student not found",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		return
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if body.Name != nil {
		student.Name = *body.Name
	}

	h.db.Save(&student)
	c.JSON(http.StatusOK, gin.H{
		"student": student,
	})
}
func (h *Handler) Destroy(c *gin.Context) {
	var student models.Student
	id := c.Param("id")

	if err := h.db.First(&student, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Student not found",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		return
	}

	h.db.Delete(&student)
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) Setup(router *gin.RouterGroup) {
	r := router.Group("/students")

	r.GET("/", h.Index)
	r.GET("/:id", h.Show)
	r.POST("/", h.Store)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Destroy)
}

func New(db *gorm.DB) *Handler {
	return &Handler{db}
}
