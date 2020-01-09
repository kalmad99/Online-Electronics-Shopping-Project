package repos

import (
	"../../entity"
	"database/sql"
	"errors"
)

// PsqlCategoryRepository implements the
// menu.CategoryRepository interface
type PsqlUserRepository struct {
	conn *sql.DB
}

// NewPsqlCategoryRepository will create an object of PsqlCategoryRepository
func NewPsqlUserRepository(Conn *sql.DB) *PsqlUserRepository {
	return &PsqlUserRepository{conn: Conn}
}

// StoreCategory stores new category information to database
func (ur *PsqlUserRepository) Regisration(u entity.User) error {

	_, err := ur.conn.Exec("INSERT INTO users (name,email,password)"+
		" values(?, ?, ?)", u.Name, u.Email, u.Password)
	if err != nil {
		//panic(err)
		return errors.New("Insertion has failed")
	}

	return nil
}
