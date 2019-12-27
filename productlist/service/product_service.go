package service

import (
	"../../entity"
	"../../productlist"
)

// ProductService implements productlist.ProductService interface
type ProductService struct {
	productRepo productlist.ProductRepository
}

// NewProductService will create new ProductService object
func NewProductService(ProRepo productlist.ProductRepository) *ProductService {
	return &ProductService{productRepo: ProRepo}
}

// Products returns list of products
func (ps *ProductService) Products() ([]entity.Product, error) {

	products, err := ps.productRepo.Products()

	if err != nil {
		return nil, err
	}

	return products, nil
}

// StoreProduct persists new product information
func (ps *ProductService) StoreProduct(product entity.Product) error {

	err := ps.productRepo.StoreProduct(product)

	if err != nil {
		return err
		//return errors.New("Insertion has failed")
	}

	return nil
}

// Product returns a product object with a given id
func (ps *ProductService) Product(id int) (entity.Product, error) {

	p, err := ps.productRepo.Product(id)

	if err != nil {
		return p, err
	}

	return p, nil
}

// UpdateProduct updates a product with new data
func (ps *ProductService) UpdateProduct(product entity.Product) error {

	err := ps.productRepo.UpdateProduct(product)

	if err != nil {
		return err
	}

	return nil
}

// DeleteProduct delete a product by its id
func (ps *ProductService) DeleteProduct(id int) error {

	err := ps.productRepo.DeleteProduct(id)
	if err != nil {
		return err
	}
	return nil
}
