package service

import (
	"../../../entity"
	"../../productpage"
)

// ItemService implements menu.ItemService interface
type ItemService struct {
	itemRepo productpage.ItemRepository
}

// NewItemService returns new ItemService object
func NewItemService(itemRepository productpage.ItemRepository) productpage.ItemService {
	return &ItemService{itemRepo: itemRepository}
}

// Items returns all stored food menu items
func (is *ItemService) Items() ([]entity.Product, []error) {
	itms, errs := is.itemRepo.Items()
	if len(errs) > 0 {
		return nil, errs
	}
	return itms, nil
}

// Item retrieves a food menu item by its id
func (is *ItemService) Item(id uint) (*entity.Product, []error) {
	itm, errs := is.itemRepo.Item(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return itm, nil
}

// UpdateItem updates a given food menu item
func (is *ItemService) UpdateItem(item *entity.Product) (*entity.Product, []error) {
	itm, errs := is.itemRepo.UpdateItem(item)
	if len(errs) > 0 {
		return nil, errs
	}
	return itm, nil
}

// DeleteItem deletes a given food menu item
func (is *ItemService) DeleteItem(id uint) (*entity.Product, []error) {
	itm, errs := is.itemRepo.DeleteItem(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return itm, nil
}

// StoreItem stores a given food menu item
func (is *ItemService) StoreItem(item *entity.Product) (*entity.Product, []error) {
	itm, errs := is.itemRepo.StoreItem(item)
	if len(errs) > 0 {
		return nil, errs
	}
	return itm, nil
}
func (is *ItemService) SearchProduct(index string) ([]entity.Product, error) {
	products, err := is.itemRepo.SearchProduct(index)

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (is *ItemService) RateProduct(pro *entity.Product) (*entity.Product, []error) {

	prowithrate, err := is.itemRepo.RateProduct(pro)
	if err != nil {
		return prowithrate, err
	}
	return prowithrate, nil
}

func (is *ItemService) StoreItemCateg(product *entity.Product) []error {

	err := is.itemRepo.StoreItemCateg(product)
	if err != nil {
		return err
	}
	return nil
}
