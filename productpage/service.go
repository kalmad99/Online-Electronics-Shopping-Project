package productpage

import (
	"../entity"
)

// ProductService specifies product services
type ProductService interface {
	Products() ([]entity.Product, error)
	Product(id int) (entity.Product, error)
	UpdateProduct(product entity.Product) error
	DeleteProduct(id int) error
	StoreProduct(product entity.Product) error
	//RateProduct(product entity.Product) (entity.Product, []error)
	SearchProduct(index string) ([]entity.Product, error)
}
