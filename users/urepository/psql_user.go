package urepository

import (
	"../../entity"
	"database/sql"
	"errors"
)

// PsqlCategoryRepository implements the
// menu.CategoryRepository interface
var dbSessions = map[string]string{} // session ID, user ID

type PsqlUserRepository struct {
	conn *sql.DB
}

// NewPsqlCategoryRepository will create an object of PsqlCategoryRepository
func NewPsqlUserRepository(Conn *sql.DB) *PsqlUserRepository {
	return &PsqlUserRepository{conn: Conn}
}

// Categories returns all cateogories from the database
func (ur *PsqlUserRepository) Users() ([]entity.User, error) {

	rows, err := ur.conn.Query("SELECT * FROM users;")
	if err != nil {
		return nil, errors.New("Could not query the database")
	}
	defer rows.Close()

	usrs := []entity.User{}

	for rows.Next() {
		user := entity.User{}
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		usrs = append(usrs, user)
	}

	return usrs, nil
}

// Category returns a category with a given id
func (ur *PsqlUserRepository) User(id int) (entity.User, error) {

	row := ur.conn.QueryRow("SELECT * FROM users WHERE id = ?", id)

	u := entity.User{}

	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		return u, err
	}

	return u, nil
}

// UpdateCategory updates a given object with a new data
func (ur *PsqlUserRepository) UpdateUser(u entity.User) error {

	_, err := ur.conn.Exec("UPDATE users SET name=?,email=?,password=? WHERE id=?",
		u.Name, u.Email, u.Password, u.ID)
	if err != nil {
		return errors.New("Update has failed")
	}

	return nil
}

// DeleteCategory removes a category from a database by its id
func (ur *PsqlUserRepository) DeleteUser(id int) error {

	_, err := ur.conn.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		return errors.New("Delete has failed")
	}

	return nil
}

// StoreCategory stores new category information to database
func (ur *PsqlUserRepository) StoreUser(u entity.User) error {

	_, err := ur.conn.Exec("INSERT INTO users (name,email,phone,password)"+
		" values(?, ?, ?, ?)", u.Name, u.Email, u.Phone, u.Password)

	if err != nil {
		//panic(err)
		return errors.New("Insertion has failed")
	}

	return nil
}
