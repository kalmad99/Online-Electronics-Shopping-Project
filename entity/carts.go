package entity

import "time"

// Cart represents customer Cart
type Cart struct {
	ID        uint      `json:"id"`
	UserId    uint      `json:"userId" gorm:"type:varchar(255)"`
	ProductId uint      `json:"ProductId"`
	AddedTime time.Time `json:"added_date" gorm:"type:varchar(100);not null; unique"`
	Price     float64   `json:"price" gorm:"type:varchar(255);not null; unique"`
}
