package server

import (
	"net/http"
	"os"
	"tink/handlers/auth"
	"tink/handlers/profile"
	"tink/handlers/student"
	"tink/handlers/teacher"
	"tink/middlewares"

	"github.com/gin-contrib/cors"
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
	group.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})
	group.Use(middlewares.Auth())

	auth.New(s.db).Setup(group)
	profile.New(s.db).Setup(group)
	student.New(s.db).Setup(group)
	teacher.New(s.db).Setup(group)
}

func New(db *gorm.DB) *Server {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("ORIGIN")}
	config.AllowCredentials = true
	router.Use(cors.New(config))
	return &Server{
		router: router,
		db:     db,
	}
}
