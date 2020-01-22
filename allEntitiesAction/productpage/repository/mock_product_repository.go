package repository

import (
	"errors"

	"../../../entity"
	"../../productpage"
	"github.com/jinzhu/gorm"
)

// MockCategoryRepo implements the menu.CategoryRepository interface
type MockProductRepo struct {
	conn *gorm.DB
}

// NewMockCategoryRepo will create a new object of MockCategoryRepo
func NewMockProductRepo(db *gorm.DB) productpage.ItemRepository {
	return &MockProductRepo{conn: db}
}

// Categories returns all fake categories
func (mProRepo *MockProductRepo) Items() ([]entity.Product, []error) {
	prds := []entity.Product{entity.ProductMock}
	return prds, nil
}

// Category retrieve a fake category with id 1
func (mProRepo *MockProductRepo) Item(id uint) (*entity.Product, []error) {
	pro := entity.ProductMock
	if id == 1 {
		return &pro, nil
	}
	return nil, []error{errors.New("Not found")}
}

// UpdateCategory updates a given fake category
func (mProRepo *MockProductRepo) UpdateItem(product *entity.Product) (*entity.Product, []error) {
	pro := entity.ProductMock
	return &pro, nil
}

// DeleteCategory deletes a given category from the database
func (mProRepo *MockProductRepo) DeleteItem(id uint) (*entity.Product, []error) {
	pro := entity.ProductMock
	if id != 1 {
		return nil, []error{errors.New("Not found")}
	}
	return &pro, nil
}

// StoreCategory stores a given mock category
func (mProRepo *MockProductRepo) StoreItem(product *entity.Product) (*entity.Product, []error) {
	pro := product
	return pro, nil
}

func (mProRepo *MockProductRepo) SearchProduct(index string) ([]entity.Product, error) {
	prds := []entity.Product{entity.ProductMock}
	return prds, nil
}

func (mProRepo *MockProductRepo) RateProduct(product *entity.Product) (*entity.Product, []error) {
	pro := product
	return pro, nil
}

func (mProRepo *MockProductRepo) StoreItemCateg(product *entity.Product) []error {
	return nil
}
