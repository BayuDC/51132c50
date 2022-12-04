package profile

import (
	"tink/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func (h *Handler) Index(c *gin.Context) {

}
func (h *Handler) UpdatePassword(c *gin.Context) {

}

func (h *Handler) Setup(router *gin.RouterGroup) {
	r := router.Group("/profile")

	r.Use(middlewares.Guard())
	r.GET("/", h.Index)
	r.PATCH("/password", h.UpdatePassword)
}

func New(db *gorm.DB) *Handler {
	return &Handler{db}
}
