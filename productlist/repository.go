package productlist

import (
	"../entity"
)

// ProductRepository specifies menu product related database operations
type ProductRepository interface {
	Products() ([]entity.Product, error)
	Product(id int) (entity.Product, error)
	UpdateProduct(product entity.Product) error
	DeleteProduct(id int) error
	StoreProduct(product entity.Product) error
}
