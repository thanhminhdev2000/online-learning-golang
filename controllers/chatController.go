package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"online-learning-golang/chat"
	"online-learning-golang/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// @Summary      WebSocket Chat Connection
// @Description  Thiết lập kết nối WebSocket cho chat realtime
// @Tags         chat
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      101  {string}  string    "Switching Protocols to WebSocket"
// @Failure      400  {object}  models.Error
// @Failure      401  {object}  models.Error
// @Router       /ws [get]
func HandleWebSocket(manager *chat.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userId")

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade connection: %v", err)
			return
		}

		client := &chat.Client{
			ID:     conn.RemoteAddr().String(),
			UserID: userID,
			Conn:   conn,
			Send:   make(chan []byte, 256),
		}

		manager.Register <- client

		go client.WritePump()
		go client.ReadPump(manager)
	}
}

// @Summary      Get Chat History
// @Description  Lấy lịch sử chat với phân trang
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        limit   query     int  false  "Số lượng tin nhắn mỗi trang"  default(50)
// @Param        offset  query     int  false  "Vị trí bắt đầu"               default(0)
// @Success      200    {array}    chat.Message
// @Failure      500    {object}   models.Error
// @Router       /history [get]
func GetChatHistory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := 50
		offset := 0

		if limitStr := c.Query("limit"); limitStr != "" {
			limit, _ = strconv.Atoi(limitStr)
		}
		if offsetStr := c.Query("offset"); offsetStr != "" {
			offset, _ = strconv.Atoi(offsetStr)
		}

		query := `
			SELECT id, content, senderId, createdAt 
			FROM chat_messages 
			ORDER BY createdAt DESC 
			LIMIT ? OFFSET ?
		`

		rows, err := db.Query(query, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
			return
		}
		defer rows.Close()

		var messages []chat.Message
		for rows.Next() {
			var msg chat.Message
			err := rows.Scan(&msg.ID, &msg.Content, &msg.SenderID, &msg.Timestamp)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
				return
			}
			messages = append(messages, msg)
		}

		c.JSON(http.StatusOK, messages)
	}
}
