package assignment

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
	var assignment models.Assignment
	var id = c.Param("id")

	if err := h.db.Preload("Course").Take(&assignment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Assignment not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Set("assignment", &assignment)
	c.Next()
}
func (h *Handler) Show(c *gin.Context) {
	var assignment = c.MustGet("assignment").(*models.Assignment)
	c.JSON(http.StatusOK, gin.H{"assignment": assignment})
}
func (h *Handler) Check(c *gin.Context) {
	assignment := c.MustGet("assignment").(*models.Assignment)
	user := (c.MustGet("user")).(*middlewares.User)

	if !assignment.Course.Check(h.db, user) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	c.Status(http.StatusOK)
	c.Next()
}

func (h *Handler) Setup(r *gin.RouterGroup) {
	router := r.Group("")
	router.Use(h.Load)
	router.Use(h.Check)
	router.GET("/assignments/:id", h.Show)
}

func New(db *gorm.DB) *Handler {
	return &Handler{db}
}
