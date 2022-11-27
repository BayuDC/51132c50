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

func (h *Handler) Index(c *gin.Context) {
	var teachers []models.Teacher
	h.db.Find(&teachers)
	c.JSON(http.StatusOK, gin.H{
		"teachers": teachers,
	})
}
func (h *Handler) Show(c *gin.Context) {
	var teacher models.Student
	id := c.Param("id")

	if err := h.db.First(&teacher, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Teacher not found",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"teacher": teacher,
	})
}
func (h *Handler) Store(c *gin.Context) {
	var body CreateTeacherSchema

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	teacher := models.Teacher{
		Name: body.Name,
	}
	h.db.Create(&teacher)
	c.JSON(http.StatusCreated, gin.H{
		"teacher": teacher,
	})
}
func (h *Handler) Update(c *gin.Context) {
	var body UpdateTeacherSchema
	var teacher models.Teacher
	id := c.Param("id")

	if err := h.db.First(&teacher, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Teacher not found",
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
		teacher.Name = *body.Name
	}

	h.db.Save(&teacher)
	c.JSON(http.StatusOK, gin.H{
		"teacher": teacher,
	})
}
func (h *Handler) Destroy(c *gin.Context) {
	var teacher models.Teacher
	id := c.Param("id")

	if err := h.db.First(&teacher, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Teacher not found",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		return
	}

	h.db.Delete(&teacher)
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) Setup(router *gin.RouterGroup) {
	r := router.Group("/teachers")

	r.Use(middlewares.Guard())
	r.GET("/", h.Index)
	r.GET("/:id", h.Show)
	r.POST("/", h.Store)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Destroy)
}

func New(db *gorm.DB) *Handler {
	return &Handler{db}
}
