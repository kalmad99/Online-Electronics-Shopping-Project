package entity

type Bank struct {
	ID        uint
	AccountNo string  `gorm:"type:varchar(255);not null"`
	Balance   float64 `gorm:"not null"`
}
