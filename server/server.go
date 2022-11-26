package server

import (
	"net/http"
	"tink/handlers/student"
	"tink/handlers/teacher"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
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

	student.New().Setup(group)
	teacher.New().Setup(group)
}

func New() *Server {
	router := gin.Default()
	return &Server{
		router: router,
	}
}
