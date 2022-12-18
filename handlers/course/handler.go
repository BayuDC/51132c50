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

	c.Set("course", &course)
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
		if errors.Is(err, models.TeacherNotFound) {
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

	course := (c.MustGet("course")).(*models.Course)

	if body.Name != nil {
		course.Name = *body.Name
	}
	if body.Teacher != nil {
		if *body.Teacher == 0 {
			course.TeacherId = nil
			course.Teacher = nil
		} else {
			course.TeacherId = body.Teacher
		}
	}

	if err := h.db.Omit("Teacher").Save(&course).Error; err != nil {
		if errors.Is(err, models.TeacherNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"course": course})
}

func (h *Handler) Destroy(c *gin.Context) {
	course := (c.MustGet("course")).(*models.Course)

	if err := h.db.Delete(&course).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) AddMember(c *gin.Context) {
	var body ManageCourseMemberSchema

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course := (c.MustGet("course")).(*models.Course)
	course.TeacherId = nil

	result := make([]ManageCourseMemberResult, len(body.Students))
	students := []models.Student{}

	tx := h.db.Begin()
	for i, id := range body.Students {
		var student models.Student
		result[i].Id = id

		if tx.Take(&student, id).Error != nil {
			result[i].Message = "Failed because student was not found"
		} else {
			result[i].Message = "Successfully added student to this course"
			students = append(students, student)
		}
	}
	if err := tx.Model(&course).Omit("Students.*").Association("Students").Append(students); err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()

	c.JSON(http.StatusMultiStatus, result)
}
func (h *Handler) RemoveMember(c *gin.Context) {
	var body ManageCourseMemberSchema

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course := (c.MustGet("course")).(*models.Course)
	course.TeacherId = nil

	result := make([]ManageCourseMemberResult, len(body.Students))
	students := []models.Student{}

	tx := h.db.Begin()
	for i, id := range body.Students {
		var student models.Student
		result[i].Id = id

		if tx.Take(&student, id).Error != nil {
			result[i].Message = "Failed because student was not found"
		} else {
			result[i].Message = "Successfully removed student from this course"
			students = append(students, student)
		}
	}
	if err := tx.Model(&course).Omit("Students.*").Association("Students").Delete(students); err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()

	c.JSON(http.StatusMultiStatus, result)
}
func (h *Handler) ShowMember(c *gin.Context) {
	var students []models.Student
	var course = (c.MustGet("course")).(*models.Course)
	h.db.Model(&course).
		Select("students.id, students.fullname, users.username").
		Joins("left join users on users.id = user_id").
		Association("Students").
		Find(&students)
	c.JSON(http.StatusOK, gin.H{"course": course, "students": students})
}
func (h *Handler) Check(c *gin.Context) {
	course := (c.MustGet("course")).(*models.Course)
	user := (c.MustGet("user")).(*middlewares.User)

	c.Status(http.StatusForbidden)
	switch user.Role {
	case "teacher":
		if user.Userable != *course.TeacherId {
			return
		}
	case "student":
		if h.db.Model(&course).
			Where("students.id = ?", user.Userable).
			Association("Students").
			Count() == 0 {
			return
		}
	default:
		return
	}

	c.Status(http.StatusOK)
	c.Next()
}

func (h *Handler) Setup(r *gin.RouterGroup) {
	router := r.Group("")
	router.Use(middlewares.Guard())
	router.GET("/courses", h.Index)
	router.GET("/courses/:id", h.Load, h.Show)
	router.GET("/courses/:id/students", h.Load, h.ShowMember)
	router.HEAD("/courses/:id/check", h.Load, h.Check)
	router.Use(middlewares.Gate("admin"))
	router.POST("/courses", h.Store)
	router.POST("/courses/:id/students", h.Load, h.AddMember)
	router.PUT("/courses/:id", h.Load, h.Update)
	router.DELETE("/courses/:id", h.Load, h.Destroy)
	router.DELETE("/courses/:id/students", h.Load, h.RemoveMember)
}

func New(db *gorm.DB) *Handler {
	return &Handler{db}
}
