package entity

// Item represents food menu items
type Product struct {
	ID         uint
	Name       string `gorm:"type:varchar(255);not null"`
	CategoryID uint   `gorm:"many2many:product_categories"`
	//Category []Category `gorm:"many2many:product_categories"`
	Quantity int
	Price    float64
	//Seller string
	Description string
	Image       string `gorm:"type:varchar(255)"`
	Rating      float64
	RatersCount float64
}
