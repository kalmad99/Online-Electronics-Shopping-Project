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
	//errs := cartRepo.conn.First(&crt, id).GetErrors()
	//if len(errs) > 0 {
	//	return nil, errs
	//}
	//return &crt, errs
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
//func (cartRepo *CartGormRepo) GetUserCart(user *entity.User) (*entity.Cart, []error) {
//	userCart := entity.Cart{}
//	errs := cartRepo.conn.Where("userid = ?", user.ID).Find(&userCart).GetErrors()
//	if len(errs) > 0 {
//		return nil, errs
//	}
//	return &userCart, errs
//}

func (cartRepo *CartGormRepo) AddtoCart(cart *entity.Cart) (*entity.Cart, []error) {
	//var quantity uint = 1
	//row := cartRepo.conn.Raw("Select quantity from carts where user_id=$1 and product_id=$2", cart.UserId, cart.ProductId)
	//if row.RowsAffected == 0{
	//	errs := cartRepo.conn.Exec("Insert into carts (user_id, product_id, added_time, quantity, price) values ($1, $2, $3, $4, $5)",
	//		cart.UserId, cart.ProductId, cart.AddedTime, quantity, cart.Price).GetErrors()
	//	return nil, errs
	//}else{
	//	errs := cartRepo.conn.Exec("Update carts set quantity=$1 where user_id=$2 and product_id=$3", quantity+1, cart.UserId, cart.ProductId).GetErrors()

	crt := cart
	errs := cartRepo.conn.Create(crt).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return crt, errs
}

func (cartRepo *CartGormRepo) DeleteCart(usr *entity.User) (*entity.Cart, []error) {
	//crt, errs := cartRepo.GetSingleCart(usr.ID)
	//
	//if len(errs) > 0 {
	//	return nil, errs
	//}

	errs := cartRepo.conn.Exec("Delete from carts where user_id=$1", usr.ID).GetErrors()
	//errs = cartRepo.conn.Delete(crt, usr.ID).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return nil, errs
}

func (cartRepo *CartGormRepo) UpdateCart(cart *entity.Cart) (*entity.Cart, []error) {
	crt := cart
	//row := cartRepo.conn.Raw("with cte(array1, array2) as (values(Select items_id from carts where user_id=$1), array($2)))" +
	//	"select array_agg(elem) from cte, unnest(array1) elem where elem <> all(array2)", cart.UserId, cart.ProductId)
	//errs := cartRepo.conn.Save(row).GetErrors()
	//errs := cartRepo.conn.Delete(crt, cart.ProductId, cart.UserId).GetErrors()
	errs := cartRepo.conn.Exec("Delete from carts where user_id=$1 and product_id=$2;", crt.UserId, crt.ProductId).GetErrors()
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
