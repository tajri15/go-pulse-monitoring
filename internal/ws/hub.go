package ws

import "log"

type Hub struct {
	clients    map[int64]*Client
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		clients:    make(map[int64]*Client),
	}
}

// Send mengirimkan pesan ke user ID tertentu jika terhubung
func (h *Hub) Send(userID int64, message []byte) {
	// Cek apakah client untuk userID ini ada dan terhubung
	if client, ok := h.clients[userID]; ok {
		// Kirim pesan ke channel Send milik client tersebut
		// Gunakan select untuk mencegah blocking jika channel penuh
		select {
		case client.Send <- message:
			log.Printf("Message sent to user ID: %d", userID)
		default:
			// Jika channel Send penuh, anggap client bermasalah dan tutup koneksi
			close(client.Send)
			delete(h.clients, client.UserID)
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if oldClient, ok := h.clients[client.UserID]; ok {
				close(oldClient.Send)
				delete(h.clients, client.UserID)
			}
			h.clients[client.UserID] = client
			log.Printf("Client registered for user ID: %d", client.UserID)

		case client := <-h.Unregister:
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
				log.Printf("Client unregistered for user ID: %d", client.UserID)
			}
		}
	}
}