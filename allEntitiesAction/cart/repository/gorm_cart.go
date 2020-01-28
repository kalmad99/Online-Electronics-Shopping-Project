package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/cart"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
)

// CartGormRepo implements the cart.CartRepository interface
type CartGormRepo struct {
	conn *gorm.DB
}

// NewCartGormRepo returns new object of CartGormRepo
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

func (cartRepo *CartGormRepo) GetSingleCart(id uint) ([]entity.Cart, []error) {
	rows, err := cartRepo.conn.Raw("Select * from carts where user_id=$1", id).Rows()

	defer rows.Close()

	crts := []entity.Cart{}

	for rows.Next() {
		cart := entity.Cart{}
		err = rows.Scan(&cart.ID, &cart.UserId, &cart.ProductId, &cart.AddedTime, &cart.Price)
		if err!=nil{
			panic(err)
		}
		crts = append(crts, cart)
	}
	return crts, nil
}

func (cartRepo *CartGormRepo) GetUserCart(user *entity.User) ([]entity.Product, []error) {
	prds := []entity.Product{}
	rows, err := cartRepo.conn.Raw("select product_id from carts where user_id = ?", user.ID).Rows()// (*sql.Rows, error)
	if err!=nil{
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		product := entity.Product{}
		err = rows.Scan(&product.ID)
		if err != nil {
			panic(err)
		}
		prds = append(prds, product)
	}
	return prds, nil
}
func (cartRepo *CartGormRepo) AddtoCart(cart *entity.Cart) (*entity.Cart, []error) {
	crt := cart
	errs := cartRepo.conn.Create(crt).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return crt, errs
}

func (cartRepo *CartGormRepo) DeleteCart(usr *entity.User) (*entity.Cart, []error) {
	errs := cartRepo.conn.Exec("Delete from carts where user_id=$1", usr.ID).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return nil, errs
}

func (cartRepo *CartGormRepo) UpdateCart(cart *entity.Cart) (*entity.Cart, []error) {
	crt := cart
	errs := cartRepo.conn.Exec("Delete from carts where user_id=$1 and product_id=$2;", crt.UserId, crt.ProductId).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return crt, errs
}
