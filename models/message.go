package models

import "time"

type ChatMessage struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Content   string    `json:"content"`
    SenderID  int       `json:"senderId"`
    CreatedAt time.Time `json:"createdAt"`
} 