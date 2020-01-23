package order

import (
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
)

// OrderRepository specifies customer menu order related database operations
type OrderRepository interface {
	Orders() ([]entity.Order, []error)
	Order(id uint) (*entity.Order, []error)
	CustomerOrders(customer *entity.User) ([]entity.Order, []error)
	UpdateOrder(order *entity.Order) (*entity.Order, []error)
	DeleteOrder(id uint) (*entity.Order, []error)
	StoreOrder(order *entity.Order) (*entity.Order, []error)
}
