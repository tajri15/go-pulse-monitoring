package api

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tajri15/go-pulse-monitoring/internal/db"
	"github.com/tajri15/go-pulse-monitoring/internal/ws"
)

// Server adalah struct utama yang menampung semua dependensi server.
type Server struct {
	store  *db.Store
	hub    *ws.Hub
	router *gin.Engine
}

// NewServer membuat instance server baru dan mengatur semua rute.
func NewServer(store *db.Store, hub *ws.Hub) *Server {
	server := &Server{
		store: store,
		hub:   hub,
	}
	router := gin.Default()

	// --- Konfigurasi Middleware CORS ---
	// Ini penting agar browser tidak memblokir permintaan dari frontend Vue Anda.
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"} // Izinkan frontend dev server
	config.AllowMethods = []string{"GET", "POST", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))
	// --------------------------------

	// Endpoint WebSocket (tidak perlu dilindungi middleware API karena punya auth sendiri)
	router.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(server.hub, c)
	})

	// Rute publik untuk autentikasi
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/register", server.registerUser)
		authRoutes.POST("/login", server.loginUser)
	}

	// Grup rute yang dilindungi oleh middleware otentikasi JWT
	api := router.Group("/api").Use(authMiddleware())
	{
		api.POST("/sites", server.createSite)
		api.GET("/sites", server.listSites)
		api.DELETE("/sites/:id", server.deleteSite)
	}

	server.router = router
	return server
}

// Start menjalankan server HTTP pada alamat yang diberikan.
func (server *Server) Start(address string) error {
	log.Printf("Starting server on %s", address)
	return server.router.Run(address)
}

// errorResponse adalah helper untuk membuat response error yang konsisten.
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}