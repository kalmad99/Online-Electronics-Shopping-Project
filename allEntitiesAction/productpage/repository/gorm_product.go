package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/productpage"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
	"log"
	"math"
)

// ItemGormRepo implements the menu.ItemRepository interface
type ItemGormRepo struct {
	conn *gorm.DB
}

// NewItemGormRepo will create a new object of ItemGormRepo
func NewItemGormRepo(db *gorm.DB) productpage.ItemRepository {
	return &ItemGormRepo{conn: db}
}

// Items returns all items stored in the database
func (itemRepo *ItemGormRepo) Items() ([]entity.Product, []error) {
	items := []entity.Product{}
	errs := itemRepo.conn.Find(&items).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return items, errs
}

// Item retrieves a an item by its id from the database
func (itemRepo *ItemGormRepo) Item(id uint) (*entity.Product, []error) {
	item := entity.Product{}
	errs := itemRepo.conn.First(&item, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &item, errs
}

// UpdateItem updates a given item in the database
func (itemRepo *ItemGormRepo) UpdateItem(item *entity.Product) (*entity.Product, []error) {
	itm := item
	errs := itemRepo.conn.Save(itm).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return itm, errs
}

// DeleteItem deletes a given item from the database
func (itemRepo *ItemGormRepo) DeleteItem(id uint) (*entity.Product, []error) {
	itm, errs := itemRepo.Item(id)

	if len(errs) > 0 {
		return nil, errs
	}

	errs = itemRepo.conn.Delete(itm, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return itm, errs
}

// StoreItem stores a given item in the database
func (itemRepo *ItemGormRepo) StoreItem(item *entity.Product) (*entity.Product, []error) {
	itm := item
	errs := itemRepo.conn.Create(itm).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return itm, errs
}
func (itemRepo *ItemGormRepo) SearchProduct(index string) ([]entity.Product, error) {
	items := []entity.Product{}

	err := itemRepo.conn.Where("name ILIKE ?", "%"+index+"%").Find(&items).GetErrors()

	if len(err) != 0 {
		//return nil, err
		errors.New("Search Product Repo not working")
	}
	return items, nil
}
func (itemRepo *ItemGormRepo) RateProduct(pro *entity.Product) (*entity.Product, []error) {

	u := entity.Product{}
	item := entity.Product{}
	row := itemRepo.conn.Select("rating").First(&item).Where("id = ?", pro.ID).Scan(&u)
	log.Println("Old rate", u.Rating)
	if row.RecordNotFound() {
		panic(row.Error)
	}

	row = itemRepo.conn.Select("raters_count").First(&item).Where("id = ?", pro.ID).Scan(&u)
	log.Println("Old count", u.RatersCount)
	if row.RecordNotFound() {
		panic(row.Error)
	}

	newratings := ((u.Rating * u.RatersCount) + pro.Rating) / (u.RatersCount + 1)
	log.Println(newratings)
	log.Println("Pro ", pro.Rating)

	row = itemRepo.conn.Model(&pro).Updates(entity.Product{Rating: float64(math.Round((newratings * 2))) / 2, RatersCount: u.RatersCount + 1})

	if row.RowsAffected < 1 {
		return &item, []error{errors.New("Error")}
	}
	return &item, nil
}

func (itemRepo *ItemGormRepo) StoreItemCateg(product *entity.Product) []error {
	pro := product

	err := itemRepo.conn.Exec("Insert into product_categories (product_id, category_id) values (?, ?)", pro.ID, pro.CategoryID).GetErrors()
	if err != nil {
		return err
	}
	return nil
}
