package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tajri15/go-pulse-monitoring/internal/db"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Definisikan rute API di sini
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/register", server.registerUser)
		authRoutes.POST("/login", server.loginUser)
	}

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	log.Printf("Starting server on %s", address)
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}