package chat

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID     string
	UserID int
	Conn   *websocket.Conn
	Send   chan []byte
}

type Message struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	SenderID  int    `json:"senderId"`
	Timestamp string `json:"timestamp"`
}

type Manager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	Mutex      sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (m *Manager) Run() {
	for {
		select {
		case client := <-m.Register:
			m.Mutex.Lock()
			m.Clients[client] = true
			m.Mutex.Unlock()

		case client := <-m.Unregister:
			if _, ok := m.Clients[client]; ok {
				m.Mutex.Lock()
				delete(m.Clients, client)
				close(client.Send)
				m.Mutex.Unlock()
			}

		case message := <-m.Broadcast:
			m.Mutex.Lock()
			for client := range m.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(m.Clients, client)
				}
			}
			m.Mutex.Unlock()
		}
	}
}
