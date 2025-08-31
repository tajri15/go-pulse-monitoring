package ws

import "log"

// Hub mengelola semua client dan menyiarkan pesan.
type Hub struct {
	// Map client yang terdaftar, dengan key userID.
	clients map[int64]*Client

	// Channel untuk menerima pesan yang akan disiarkan.
	// Kita akan gunakan ini di langkah berikutnya.
	Broadcast chan []byte

	// Channel untuk mendaftarkan client baru.
	Register chan *Client

	// Channel untuk membatalkan pendaftaran client.
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		clients:    make(map[int64]*Client),
	}
}

// Run adalah event loop untuk Hub. Harus dijalankan sebagai goroutine.
func (h *Hub) Run() {
	for {
		select {
		// Jika ada client baru yang mendaftar...
		case client := <-h.Register:
			// Hanya izinkan satu koneksi per user.
			// Jika user sudah ada, tutup koneksi lama.
			if oldClient, ok := h.clients[client.UserID]; ok {
				close(oldClient.Send)
				delete(h.clients, client.UserID)
			}
			h.clients[client.UserID] = client
			log.Printf("Client registered for user ID: %d", client.UserID)

		// Jika ada client yang keluar...
		case client := <-h.Unregister:
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
				log.Printf("Client unregistered for user ID: %d", client.UserID)
			}
		
		// Jika ada pesan untuk disiarkan... (akan kita implementasikan nanti)
		case message := <-h.Broadcast:
			// Logika broadcast akan ditambahkan di langkah 8
			log.Printf("Broadcasting message: %s", message)
		}
	}
}