package entity

type Category struct {
	ID          uint
	Name        string
	Description string
	Image       string
	Items       []Product
}
