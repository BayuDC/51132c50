package core

import (
	"net/http"
	"tink/handlers/auth"
	"tink/handlers/group"
	"tink/handlers/profile"
	"tink/handlers/student"
	"tink/handlers/teacher"
	"tink/middlewares"

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
	g := s.router.Group("/api")
	g.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})
	g.Use(middlewares.Auth())

	auth.New(s.db).Setup(g)
	profile.New(s.db).Setup(g)
	student.New(s.db).Setup(g)
	teacher.New(s.db).Setup(g)
	group.New(s.db).Setup(g)
}

func CreateApp(db *gorm.DB) *Server {
	router := gin.Default()
	return &Server{
		router: router,
		db:     db,
	}
}
