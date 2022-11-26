package server

import (
	"net/http"
	"tink/handlers/student"
	"tink/handlers/teacher"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	router *gin.Engine
	db     *gorm.DB
}

func (s *Server) Run() {
	s.Setup()
	s.router.Run()
}
func (s *Server) Setup() {
	group := s.router.Group("/api")
	group.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})

	student.New(s.db).Setup(group)
	teacher.New(s.db).Setup(group)
}

func New(db *gorm.DB) *Server {
	router := gin.Default()
	return &Server{
		router: router,
		db:     db,
	}
}
