package service

import (
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/productpage"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
)

// CategoryService implements menu.CategoryService interface
type CategoryService struct {
	categoryRepo productpage.CategoryRepository
}

// NewCategoryService will create new CategoryService object
func NewCategoryService(CatRepo productpage.CategoryRepository) productpage.CategoryService {
	return &CategoryService{categoryRepo: CatRepo}
}

// Categories returns list of categories
func (cs *CategoryService) Categories() ([]entity.Category, []error) {

	categories, errs := cs.categoryRepo.Categories()

	if len(errs) > 0 {
		return nil, errs
	}

	return categories, nil
}

// StoreCategory persists new category information
func (cs *CategoryService) StoreCategory(category *entity.Category) (*entity.Category, []error) {

	cat, errs := cs.categoryRepo.StoreCategory(category)

	if len(errs) > 0 {
		return nil, errs
	}

	return cat, nil
}

// Category returns a category object with a given id
func (cs *CategoryService) Category(id uint) (*entity.Category, []error) {
	//func (cs *CategoryService) Category(id uint) (*entity.Category, []error) {

	c, err := cs.categoryRepo.Category(id)

	if err != nil {
		return c, err
	}

	return c, nil
}

// UpdateCategory updates a cateogory with new data
func (cs *CategoryService) UpdateCategory(category *entity.Category) (*entity.Category, []error) {

	cat, errs := cs.categoryRepo.UpdateCategory(category)

	if len(errs) > 0 {
		return nil, errs
	}

	return cat, nil
}

// DeleteCategory delete a category by its id
func (cs *CategoryService) DeleteCategory(id uint) (*entity.Category, []error) {

	cat, errs := cs.categoryRepo.DeleteCategory(id)

	if len(errs) > 0 {
		return nil, errs
	}

	return cat, nil
}

// ItemsInCategory returns list of menu items in a given category
func (cs *CategoryService) ItemsInCategory(category *entity.Category) ([]entity.Product, []error) {

	cts, errs := cs.categoryRepo.ItemsInCategory(category)

	if len(errs) > 0 {
		return nil, errs
	}

	return cts, nil

}

//
//
//// CategoryService implements menu.CategoryService interface
//type CategoryService struct {
//	categoryRepo productpage.CategoryRepository
//}
//
//// NewCategoryService will create new CategoryService object
//func NewCategoryService(CatRepo productpage.CategoryRepository) productpage.CategoryService {
//	return &CategoryService{categoryRepo: CatRepo}
//}
//
//// Categories returns list of categories
//func (cs *CategoryService) Categories() ([]entity.Category, error) {
//
//	categories, err := cs.categoryRepo.Categories()
//
//	if err != nil {
//		return nil, err
//	}
//
//	return categories, nil
//}
//
//// StoreCategory persists new category information
//func (cs *CategoryService) StoreCategory(category entity.Category) (entity.Category, error) {
//
//	cat, err := cs.categoryRepo.StoreCategory(category)
//
//	if err != nil {
//		return cat, err
//	}
//
//	return cat, nil
//}
//
//// Category returns a category object with a given id
//func (cs *CategoryService) Category(id uint) (entity.Category, error) {
//	//func (cs *CategoryService) Category(id uint) (*entity.Category, []error) {
//
//	c, err := cs.categoryRepo.Category(id)
//
//	if err != nil {
//		return c, err
//	}
//
//	return c, nil
//}
//
//// UpdateCategory updates a cateogory with new data
//func (cs *CategoryService) UpdateCategory(category entity.Category) error {
//
//	err := cs.categoryRepo.UpdateCategory(category)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// DeleteCategory delete a category by its id
//func (cs *CategoryService) DeleteCategory(id uint) error {
//
//	err := cs.categoryRepo.DeleteCategory(id)
//
//	if err != nil{
//		return err
//	}
//
//	return nil
//}
//
//// ItemsInCategory returns list of menu items in a given category
//func (cs *CategoryService) ItemsInCategory(category entity.Category) ([]entity.Product, error) {
//
//	cts, err := cs.categoryRepo.ItemsInCategory(category)
//
//	if err != nil{
//		return cts, err
//	}
//
//	return cts, nil
//
//}
