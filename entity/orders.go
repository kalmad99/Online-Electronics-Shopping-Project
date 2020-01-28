package entity

import "time"

// Order represents customer order
type Order struct {
	ID        uint
	CreatedAt time.Time
	UserID    uint
	ItemID    uint
	Quantity  uint
}
