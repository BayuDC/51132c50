package group

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
	var group models.Group

	if err := h.db.First(&group, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Group not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Set("group", group)
	c.Next()
}

func (h *Handler) Index(c *gin.Context) {
	var groups []models.Group
	if err := h.db.Order("id ASC").Find(&groups).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

func (h *Handler) Show(c *gin.Context) {
	group := c.MustGet("group")
	c.JSON(http.StatusOK, gin.H{"group": group})
}

func (h *Handler) Store(c *gin.Context) {
	var body CreateGroupSchema

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group := models.Group{
		Name: body.Name,
	}

	if err := h.db.Create(&group).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"group": group})
}

func (h *Handler) Update(c *gin.Context) {
	var body UpdateGroupSchema

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group := (c.MustGet("group")).(models.Group)

	if body.Name != nil {
		group.Name = *body.Name
	}

	if err := h.db.Save(&group).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"group": group})
}

func (h *Handler) Delete(c *gin.Context) {
	group := (c.MustGet("group")).(models.Group)

	if err := h.db.Delete(&group).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) Setup(r *gin.RouterGroup) {
	router := r.Group("")
	router.Use(middlewares.Guard())
	router.GET("/groups", h.Index)
	router.GET("/groups/:id", h.Load, h.Show)
	router.Use(middlewares.Gate("admin"))
	router.POST("/groups", h.Store)
	router.PUT("/groups/:id", h.Load, h.Update)
	router.DELETE("/groups/:id", h.Load, h.Delete)
}
func New(db *gorm.DB) *Handler {
	return &Handler{db}
}
