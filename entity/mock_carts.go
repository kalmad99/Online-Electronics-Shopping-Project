package entity

import "time"

var AddToCartMock = Cart{
	ID:        1,
	AddedTime: time.Time{},
	UserId:    1,
	ProductId: 1,
}
