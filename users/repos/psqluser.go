package repos

import (
	"../../entity"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type PsqlUserRepository struct {
	conn *sql.DB
}

// NewPsqlCategoryRepository will create an object of PsqlCategoryRepository
func NewPsqlUserRepository(Conn *sql.DB) *PsqlUserRepository {
	return &PsqlUserRepository{conn: Conn}
}

func (ur *PsqlUserRepository) Registration(u entity.User) error {
	var useremail string

	err := ur.conn.QueryRow("SELECT email FROM users WHERE email=?", u.Email).Scan(&useremail)

	switch {
	case err == sql.ErrNoRows:
		hashedpass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			//http.Error(w, "Server error, unable to create your account.", 500)
			//return errors.New("server error, unable to create your account")
			panic(err.Error())
		}
		_, err = ur.conn.Exec("INSERT INTO users (name,email,phone,password)"+
			" values(?, ?, ?, ?)", u.Name, u.Email, u.Phone, hashedpass)

		if err != nil {
			//			http.Error(w, "Server error, unable to create your account.", 500)
			//return errors.New("server error, unable to create your account")
			panic(err.Error())
		}
	case err != sql.ErrNoRows:
		//http.Error(w, "This email is already taken", 500)
		return errors.New("this address is already taken")

	case err != nil:
		//http.Error(w, "Server error, unable to create your account.", 500)
		panic(err.Error())
		//return errors.New("server error, unable to create your account")
	}
	return nil
}
