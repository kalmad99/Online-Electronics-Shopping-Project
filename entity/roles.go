package entity

// Role repesents application user roles
type Role struct {
	ID    uint
	Name  string `gorm:"type:varchar(255)"`
	Users []User
}
