package assignment

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

func (h *Handler) SetSchedule(c *gin.Context) {
	var body SetScheduleSchema

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if body.OpenAt != nil && body.CloseAt != nil {
		if (*body.CloseAt).Sub(*body.OpenAt) < 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Close time must be after open time"})
			return
		}
	}

	assignment := c.MustGet("assignment").(*models.Assignment)
	schedule := models.Schedule{
		AssignmentId: assignment.Id,
		GroupId:      body.Group,
	}

	if err := h.db.FirstOrCreate(&schedule, &schedule).Error; err != nil {
		if errors.Is(err, models.GroupNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	if body.OpenAt != nil && body.CloseAt != nil {
		schedule.OpenAt = *body.OpenAt
		schedule.CloseAt = *body.CloseAt
	}
	if err := h.db.Save(&schedule).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
func (h *Handler) SetAttachment(c *gin.Context) {
	assignment := c.MustGet("assignment").(*models.Assignment)

	form, _ := c.MultipartForm()
	files := form.File["attachments[]"]

	filedir := fmt.Sprintf("./data/assignment%d/_", assignment.Id)
	if _, err := os.Stat(filedir); os.IsNotExist(err) {
		os.MkdirAll(filedir, os.ModePerm)
	}

	for _, file := range files {
		c.SaveUploadedFile(file, filepath.Join(filedir, filepath.Base(file.Filename)))
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) Setup(r *gin.RouterGroup) {
	router := r.Group("")
	router.Use(h.Load)
	router.Use(h.Check)
	router.GET("/assignments/:id", h.Show)
	router.PATCH("/assignments/:id/schedules", middlewares.Gate("teacher"), h.SetSchedule)
	router.PATCH("/assignments/:id/attachments", middlewares.Gate("teacher"), h.SetAttachment)
}

func New(db *gorm.DB) *Handler {
	return &Handler{db}
}
