package cart

import (
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
)

type CartService interface {
	GetCarts() ([]entity.Cart, []error)
	GetSingleCart(id uint) ([]entity.Cart, []error)
	GetUserCart(user *entity.User) ([]entity.Product, []error)
	AddtoCart(cart *entity.Cart) (*entity.Cart, []error)
	DeleteCart(user *entity.User) (*entity.Cart, []error)
	UpdateCart(cart *entity.Cart) (*entity.Cart, []error)
}
