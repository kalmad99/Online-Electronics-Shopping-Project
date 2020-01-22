package entity

import "time"

// Order represents customer order
var OrderMock = Order{
	ID:        1,
	CreatedAt: time.Time{},
	UserID:    1,
	ItemID:    1,
	Quantity:  1,
}
