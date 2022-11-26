package student

import "github.com/gin-gonic/gin"

type Handler struct{}

func (h *Handler) Index(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello Student",
	})
}
func (h *Handler) Show(c *gin.Context) {
}
func (h *Handler) Store(c *gin.Context) {
}
func (h *Handler) Update(c *gin.Context) {
}
func (h *Handler) Destroy(c *gin.Context) {
}

func (h *Handler) Setup(router *gin.RouterGroup) {
	r := router.Group("/students")

	r.GET("/", h.Index)
	r.GET("/:id", h.Show)
	r.POST("/", h.Store)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Destroy)
}

func New() *Handler {
	return &Handler{}
}
