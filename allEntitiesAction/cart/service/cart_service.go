package service

import (
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/cart"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
	//"context"
	//"errors"
	//"time"
)

// CartrGormRepo implements the cart.CartRepository interface
type CartService struct {
	CartRepo cart.CartRepository
}

// NewCartGormRepo returns new object of CartGormRepo
func NewCartService(cartRepo cart.CartRepository) cart.CartService {
	return &CartService{CartRepo: cartRepo}
}

func (cs *CartService) GetCarts() ([]entity.Cart, []error) {
	ords, errs := cs.CartRepo.GetCarts()
	if len(errs) > 0 {
		return nil, errs
	}
	return ords, nil
}

func (cs *CartService) GetSingleCart(id uint) ([]entity.Cart, []error) {
	crts, errs := cs.CartRepo.GetSingleCart(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return crts, nil
}

func (cs *CartService) GetUserCart(user *entity.User) ([]entity.Product, []error) {
	prds, errs := cs.CartRepo.GetUserCart(user)
	if len(errs) > 0 {
		return nil, errs
	}
	return prds, nil

}

func (cs *CartService) AddtoCart(cart *entity.Cart) (*entity.Cart, []error) {
	crt, errs := cs.CartRepo.AddtoCart(cart)
	if len(errs) > 0 {
		return nil, errs
	}
	return crt, nil
}

func (cs *CartService) DeleteCart(usr *entity.User) (*entity.Cart, []error) {
	crt, errs := cs.CartRepo.DeleteCart(usr)
	if len(errs) > 0 {
		return nil, errs
	}
	return crt, nil
}

func (cs *CartService) UpdateCart(cart *entity.Cart) (*entity.Cart, []error) {
	crt, errs := cs.CartRepo.UpdateCart(cart)
	if len(errs) > 0 {
		return nil, errs
	}
	return crt, nil
}
