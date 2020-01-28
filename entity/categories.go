package entity

// Category represents Food Menu Category
type Category struct {
	ID          uint
	Name        string `gorm:"type:varchar(255);not null"`
	Description string
	Image       string    `gorm:"type:varchar(255)"`
	Products    []Product `gorm:"many2many:product_categories"`
}
