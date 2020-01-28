package service

import (
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/cart"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
	//"context"
	//"errors"
	//"time"
)

// OrderGormRepo implements the menu.OrderRepository interface
type CartService struct {
	CartRepo cart.CartRepository
}

// NewOrderGormRepo returns new object of OrderGormRepo
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
//func (cs *CartService) GetUserCart(user *entity.User) (*entity.Cart, []error) {
//	crt, errs := cs.CartRepo.GetUserCart(user)
//	if len(errs) > 0 {
//		return nil, errs
//	}
//	return crt, nil
//
//}

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

///////////////////////////////////////////////////////////////////////////////////////////////////////////
//func (cs *CartService) CartByID(ctx context.Context, id uint) (*entity.Cart, error) {
//	ccart, err := cs.CartRepo.GetCart(id)
//	if err != nil {
//		return nil, errors.New("Cant Create cart")
//	}
//	return ccart, nil
//}
//
//func (cs *CartService) CreateCart(ctx context.Context, cart *entity.Cart) (*entity.Cart, error) {
//	now := time.Now()
//
//	newCart := &entity.Cart{
//		UserID: cart.UserID,
//		PlacedAt: now,
//	}
//
//	cart, err := cs.CartRepo.CreateCart(ctx, newCart)
//	if err != nil {
//		panic(err.Error())
//	}
//	return cart, nil
//}
//
//func (cs *CartService) DeleteCart(ctx context.Context, cartID uint) error {
//	err := cs.DeleteCart(ctx, cartID)
//	if err != nil {
//		panic(err.Error())
//	}
//	return nil
//}
//
//func (cs *CartService) AddtoCart(cart *entity.Cart, proID uint) []error {
//	c := &entity.Cart{}
//	errs := cs.CartRepo.AddtoCart(c, proID)
//	if len(errs) > 0{
//		return errs
//	}else{
//		return nil
//	}
//}
//
//func (cs *CartService) ItemsinCart(cart *entity.Cart) ([]entity.Product, []error) {
//	products, errs := cs.CartRepo.ItemsinCart(cart)
//	if len(errs) > 0{
//		return nil, errs
//	}else{
//		return products, nil
//	}
//}
//
//func (cs *CartService) RemovefromCart(pro *entity.Product) (*entity.Cart, []error) {
//	products, errs := cs.CartRepo.RemovefromCart(pro)
//	if len(errs) > 0{
//		return nil, errs
//	}else{
//		return products, nil
//	}
//}
//
//func (cs *CartService) GetCart(id uint) (*entity.Cart, []error){
//	cart, errs := cs.CartRepo.GetCart(id)
//	if len(errs) > 0{
//		return nil, errs
//	}else {
//		return cart, nil
//	}
//}
