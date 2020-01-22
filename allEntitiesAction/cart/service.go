package cart

import (
	"../../entity"
	//"context"
)

type CartService interface {
	//CartByID(ctx context.Context, id uint) (*entity.Cart, error)
	//CreateCart(ctx context.Context, cart *entity.Cart) (*entity.Cart, error)
	//DeleteCart(ctx context.Context, cartID uint) error
	//AddtoCart(cart *entity.Cart, proID uint) []error
	//ItemsinCart(cart *entity.Cart) ([]entity.Product, []error)
	//RemovefromCart (pro *entity.Product) (*entity.Cart, []error)
	//GetCart(id uint) (*entity.Cart, []error)
	GetCarts() ([]entity.Cart, []error)
	GetSingleCart(id uint) (*entity.Cart, []error)
	GetUserCart(user *entity.User) (*entity.Cart, []error)
	AddtoCart(cart *entity.Cart) (*entity.Cart, []error)
	DeleteCart(id uint) (*entity.Cart, []error)
	UpdateCart(cart *entity.Cart) (*entity.Cart, []error)
}
