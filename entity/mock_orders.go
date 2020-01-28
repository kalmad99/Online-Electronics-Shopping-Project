package entity

import "time"

// Order represents customer order
var OrderMock = Order{
	ID:        1,
	UserID:    1,
	CreatedAt: time.Time{},
	//ItemsID:   []int64{1, 2, 3},
	ItemsID:  "{1, 2, 3}",
	//ItemsID: 2,
	Total:   200.00,
}
