package entity

import (
	"time"
)

// Order represents customer order
type Order struct {
	ID        uint
	UserID    uint      `json:"userId" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null;"`
	ItemsID   string    `json:"ProductId" gorm:"not null"`
	Total     float64   `json:"total" gorm:"type:float;not null;"`
}
