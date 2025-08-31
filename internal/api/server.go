package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tajri15/go-pulse-monitoring/internal/db"
	"github.com/tajri15/go-pulse-monitoring/internal/ws"
)

type Server struct {
	store  *db.Store
	hub    *ws.Hub
	router *gin.Engine
}

// Modifikasi NewServer untuk menerima Hub
func NewServer(store *db.Store, hub *ws.Hub) *Server {
	server := &Server{
		store: store,
		hub:   hub,
	}
	router := gin.Default()
    
	// Endpoint WebSocket (tidak perlu dilindungi middleware karena punya auth sendiri)
	router.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(server.hub, c)
	})

	// Rute publik untuk autentikasi
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/register", server.registerUser)
		authRoutes.POST("/login", server.loginUser)
	}

	// Rute yang dilindungi oleh middleware otentikasi
	api := router.Group("/api").Use(authMiddleware())
	{
		api.POST("/sites", server.createSite)
		api.GET("/sites", server.listSites)
		api.DELETE("/sites/:id", server.deleteSite)
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