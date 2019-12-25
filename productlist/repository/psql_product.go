package repository

import (
	"../../entity"
	"database/sql"
	"errors"
)

// PsqlProductRepository implements the productlist.ProductRepository interface
type PsqlProductRepository struct {
	conn *sql.DB
}

// NewPsqlProductRepository will create an object of PsqlProductRepository
func NewPsqlProductRepository(Conn *sql.DB) *PsqlProductRepository {
	return &PsqlProductRepository{conn: Conn}
}

// Products returns all products from the database
func (pr *PsqlProductRepository) Products() ([]entity.Product, error) {

	rows, err := pr.conn.Query("SELECT * FROM products;")
	if err != nil {
		return nil, errors.New("Could not query the database")
	}
	defer rows.Close()

	prds := []entity.Product{}

	for rows.Next() {
		product := entity.Product{}
		err = rows.Scan(&product.ID, &product.Name, &product.ItemType,
			&product.Quantity, &product.Price, &product.Description, &product.Image,
			&product.Rating, &product.RatersCount)
		if err != nil {
			return nil, err
		}
		prds = append(prds, product)
	}

	return prds, nil
}

// Product returns a product with a given id
func (pr *PsqlProductRepository) Product(id int) (entity.Product, error) {

	row := pr.conn.QueryRow("SELECT * FROM products WHERE id = ?", id)

	p := entity.Product{}

	err := row.Scan(&p.ID, &p.Name, &p.ItemType,
		&p.Quantity, &p.Price, &p.Description, &p.Image, &p.Rating, &p.RatersCount)
	if err != nil {
		return p, err
	}

	return p, nil
}

// UpdateProduct updates a given product with a new data
func (pr *PsqlProductRepository) UpdateProduct(c entity.Product) error {

	_, err := pr.conn.Exec("UPDATE products SET itemname=?,itemtype=?,"+
		"quantity=?,price=?,description=?,image=? WHERE id=?",
		c.Name, c.ItemType, c.Quantity, c.Price, c.Description, c.Image, c.ID)
	if err != nil {
		return errors.New("Update has failed")
	}

	return nil
}

// DeleteProduct removes a product from a database by its id
func (pr *PsqlProductRepository) DeleteProduct(id int) error {

	_, err := pr.conn.Exec("DELETE FROM products WHERE id=$1", id)
	if err != nil {
		return errors.New("Delete has failed")
	}

	return nil
}

// StoreProduct stores new product information to database
func (pr *PsqlProductRepository) StoreProduct(c entity.Product) error {

	_, err := pr.conn.Exec("INSERT INTO products (itemname,itemtype,quantity,price,description,image)"+
		" values(?, ?, ?, ?, ?, ?)", c.Name, c.ItemType, c.Quantity, c.Price, c.Description, c.Image)
	if err != nil {
		panic(err)
		return errors.New("Insertion has failed or Item already Exists")
	}

	return nil
}
