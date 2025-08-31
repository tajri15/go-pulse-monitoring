package ws

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Izinkan koneksi dari origin manapun (untuk development)
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Client adalah perantara antara koneksi WebSocket dan Hub.
type Client struct {
	Hub    *Hub
	UserID int64
	Conn   *websocket.Conn
	Send   chan []byte
}

// readPump membaca pesan dari koneksi WebSocket (untuk mendeteksi disconnect)
func (c *Client) readPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	// Atur batas waktu baca
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })

	for {
		// Cukup baca dan abaikan pesan dari client
		if _, _, err := c.Conn.ReadMessage(); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

// writePump mengirim pesan dari Hub ke koneksi WebSocket
func (c *Client) writePump() {
	ticker := time.NewTicker(50 * time.Second) // Ping ticker
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Hub menutup channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, message)
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}


// ServeWs menangani request WebSocket dari client.
func ServeWs(hub *Hub, c *gin.Context) {
	// 1. Ambil token dari query parameter
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token is required"})
		return
	}
	
	// 2. Validasi token JWT
	secretKey := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	
	var userID int64
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDFloat := claims["sub"].(float64)
		userID = int64(userIDFloat)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// 3. Upgrade koneksi ke WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// 4. Buat client baru dan daftarkan ke Hub
	client := &Client{Hub: hub, UserID: userID, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register <- client

	// Jalankan read dan write pump sebagai goroutine
	go client.writePump()
	go client.readPump()
}