package entity

import "time"

// Cart represents customer Cart
type Cart struct {
	ID        uint      `json:"id"`
	UserId    uint      `json:"userId" gorm:"not null"`
	ProductId uint      `json:"ProductId" gorm:"not null"`
	AddedTime time.Time `json:"added_date" gorm:"type:timestamp;not null;"`
	//Quantity  uint      `json:"quantity" gorm:"not null"`
	Price     float64   `json:"price" gorm:"type:float;not null;"`
}
