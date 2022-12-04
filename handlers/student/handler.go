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

func (h *Handler) Index(c *gin.Context) {
	var students []models.Student
	h.db.Table("students").
		Select("students.id, students.fullname, users.username").
		Joins("left join users on users.id = user_id").
		Find(&students)

	c.JSON(http.StatusOK, gin.H{
		"students": students,
	})
}
func (h *Handler) Show(c *gin.Context) {
	var student models.Student
	id := c.Param("id")

	if err := h.db.Preload("User").First(&student, id).Error; err != nil {
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

	student.Username = student.User.Username
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
		Fullname: body.Fullname,
		User: models.User{
			Username: body.Username,
			Role:     models.RoleStudent,
		},
	}
	if err := h.db.Create(&student).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	student.Username = student.User.Username
	c.JSON(http.StatusCreated, gin.H{
		"student": student,
	})
}

func (h *Handler) Update(c *gin.Context) {
	var body UpdateStudentSchema
	var student models.Student
	id := c.Param("id")

	if err := h.db.Preload("User").First(&student, id).Error; err != nil {
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

	if body.Fullname != nil {
		student.Fullname = *body.Fullname
	}
	if err := h.db.Omit("User").Save(&student).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	student.Username = student.User.Username
	c.JSON(http.StatusOK, gin.H{
		"student": student,
	})
}

func (h *Handler) Destroy(c *gin.Context) {
	var student models.Student
	id := c.Param("id")

	if err := h.db.Preload("User").First(&student, id).Error; err != nil {
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

	h.db.Delete(&student.User)
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) Setup(r *gin.RouterGroup) {
	r.Use(middlewares.Guard())
	r.GET("/students", h.Index)
	r.GET("/students/:id", h.Show)
	r.POST("/students", h.Store)
	r.PUT("/students/:id", h.Update)
	r.DELETE("/students/:id", h.Destroy)
}

func New(db *gorm.DB) *Handler {
	return &Handler{db}
}
