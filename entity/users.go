package entity

// User represents application user
type User struct {
	ID          uint
	Name        string `gorm:"type:varchar(255);not null"`
	Email       string `gorm:"type:varchar(255);not null; unique"`
	Phone       string `gorm:"type:varchar(100);not null; unique"`
	Password    string `gorm:"type:varchar(255)"`
	RoleID      uint
	ItemsInCart []Cart
}
