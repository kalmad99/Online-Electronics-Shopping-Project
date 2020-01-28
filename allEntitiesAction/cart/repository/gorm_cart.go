package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/cart"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
)

// OrderGormRepo implements the menu.OrderRepository interface
type CartGormRepo struct {
	conn *gorm.DB
}

// NewOrderGormRepo returns new object of OrderGormRepo
func NewCartGormRepo(db *gorm.DB) cart.CartRepository {
	return &CartGormRepo{conn: db}
}

func (cartRepo *CartGormRepo) GetCarts() ([]entity.Cart, []error) {
	carts := []entity.Cart{}
	errs := cartRepo.conn.Find(&carts).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return carts, errs
}

func (cartRepo *CartGormRepo) GetSingleCart(id uint) (*entity.Cart, []error) {
	crt := entity.Cart{}
	errs := cartRepo.conn.First(&crt, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &crt, errs
}

func (cartRepo *CartGormRepo) GetUserCart(user *entity.User) (*entity.Cart, []error) {
	userCart := entity.Cart{}
	errs := cartRepo.conn.Where("userid = ?", user.ID).First(&userCart).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &userCart, errs
}

func (cartRepo *CartGormRepo) AddtoCart(cart *entity.Cart) (*entity.Cart, []error) {
	crt := cart
	errs := cartRepo.conn.Create(crt).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return crt, errs
}

func (cartRepo *CartGormRepo) DeleteCart(id uint) (*entity.Cart, []error) {
	crt, errs := cartRepo.GetSingleCart(id)

	if len(errs) > 0 {
		return nil, errs
	}

	errs = cartRepo.conn.Delete(crt, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return crt, errs
}

func (cartRepo *CartGormRepo) UpdateCart(cart *entity.Cart) (*entity.Cart, []error) {
	crt := cart
	errs := cartRepo.conn.Save(crt).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return crt, errs
}

//////////////////////////////////////////////////////////////////////////
//func (cartRepo *CartGormRepo) CartByID(ctx context.Context, id uint) (*Cart, error) {
//	panic("implement me")
//}
//
//func (cartRepo *CartGormRepo) CreateCart(ctx context.Context, cart *Cart) (*Cart, error) {
//	panic("implement me")
//}
//
//func (cartRepo *CartGormRepo) DeleteCart(ctx context.Context, cartID uint) error {
//	panic("implement me")
//}
//
//// Orders returns all customer orders stored in the database
//func (cartRepo *CartGormRepo) AddtoCart(cart *entity.Cart, proID uint) []error {
//	c := cart
//	var result int64
//	errs := cartRepo.conn.Table("add_to_carts").Where("item_id = ? AND user_id = ?", proID, c.UserID).Count(&result).GetErrors()
//	if len(errs) > 0{
//		return errs
//	}else if result > 0{
//		//
//	}
//	proRepo := repository.ItemGormRepo{}
//	prod, _ := proRepo.Item(proID)
//	c = &entity.Cart{
//		PlacedAt:time.Now(),
//		UserID:c.UserID,
//		Items: []entity.Product{{
//			ID:proID,
//			Name:prod.Name,
//			Image:prod.Image,
//			Description:prod.Description,
//		}},
//	}
//	cartRepo.conn.Create(c)
//	return nil
//}
//
//// Order retrieve customer order by order id
//func (cartRepo *CartGormRepo) RemovefromCart(pro *entity.Product) (*entity.Cart, []error) {
//	car := &entity.Cart{}
//	errs := cartRepo.conn.Delete(pro, pro.ID).GetErrors()
//	if len(errs) > 0 {
//		return nil, errs
//	}
//	return car, nil
//}
//
//// UpdateOrder updates a given customer order in the database
//func (cartRepo *CartGormRepo) ItemsinCart(cart *entity.Cart) ([]entity.Product, []error) {
//	pro := []entity.Product{}
//	errs := cartRepo.conn.Model(cart).Related(&pro, "ItemsInCart").GetErrors()
//	if len(errs) > 0 {
//		return nil, errs
//	}
//	return pro, errs
//}
//
//func (cartRepo *CartGormRepo) GetCart(id uint) (*entity.Cart, []error){
//	crt := entity.Cart{}
//	errs := cartRepo.conn.First(&crt, id).GetErrors()
//	if len(errs) > 0 {
//		return nil, errs
//	}
//	return &crt, errs
//}
