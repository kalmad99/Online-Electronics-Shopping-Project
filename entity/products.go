package entity

type Product struct {
	ID          int
	Name        string
	ItemType    string
	Quantity    int
	Price       float64
	Seller      string
	Description string
	Image       string
	Rating      float64
	RatersCount float64
}
