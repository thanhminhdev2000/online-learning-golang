package chat

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func (client *Client) WritePump() {
	defer func() {
		client.Conn.Close()
	}()

	for {
		message, ok := <-client.Send
		if !ok {
			client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		err := client.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
	}
}

func (client *Client) ReadPump(manager *Manager) {
	defer func() {
		manager.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := Message{
			Type:      "message",
			Content:   string(message),
			SenderID:  client.UserID,
			Timestamp: time.Now().Format(time.RFC3339),
		}

		jsonMessage, _ := json.Marshal(msg)
		manager.Broadcast <- jsonMessage
	}
}
