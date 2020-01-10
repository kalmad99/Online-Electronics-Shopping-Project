package service

import (
	"../../entity"
	"../../productpage"
)

// ProductService implements productpage.ProductService interface
type ProductService struct {
	productRepo productpage.ProductRepository
}

// NewProductService will create new ProductService object
func NewProductService(ProRepo productpage.ProductRepository) *ProductService {
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

func (ps *ProductService) SearchProduct(index string) ([]entity.Product, error) {
	products, err := ps.productRepo.SearchProduct(index)

	if err != nil {
		return nil, err
	}

	return products, nil
}

//func (ps *ProductService) RateProduct(pro entity.Product) (entity.Product, []error) {
//
//	prowithrate, err := ps.productRepo.RateProduct(pro)
//	if err != nil {
//		return prowithrate, err
//	}
//	return prowithrate, nil
//}
