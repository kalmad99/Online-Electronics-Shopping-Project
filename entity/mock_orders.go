package entity

import "time"

// Order represents customer order
var OrderMock = Order{
	ID:        1,
	UserID:    1,
	CreatedAt: time.Time{},
	ItemsID:   "{1, 2, 3}",
	Total:     200.00,
}
