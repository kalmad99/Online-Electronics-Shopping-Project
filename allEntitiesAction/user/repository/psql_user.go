package repository

import (
	"database/sql"
	"errors"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// PsqlCategoryRepository implements the
// menu.CategoryRepository interface
var dbSessions = map[string]string{} // session ID, user ID

type UserRepositoryImpl struct {
	conn *sql.DB
}

// NewPsqlCategoryRepository will create an object of PsqlCategoryRepository
func NewUserRepositoryImpl(Conn *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{conn: Conn}
}

// Categories returns all cateogories from the database
func (ur *UserRepositoryImpl) Users() ([]entity.User, error) {

	rows, err := ur.conn.Query("SELECT * FROM users;")
	if err != nil {
		return nil, errors.New("Could not query the database")
	}
	defer rows.Close()

	usrs := []entity.User{}

	for rows.Next() {
		user := entity.User{}
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Password)
		if err != nil {
			return nil, err
		}
		usrs = append(usrs, user)
	}

	return usrs, nil
}

func (ur *UserRepositoryImpl) Login(email string) (entity.User, error) {

	row := ur.conn.QueryRow("SELECT * FROM users WHERE email = $1", email)

	u := entity.User{}

	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Password)

	if err != nil {
		return u, errors.New("username or Password is incorrect")
	}
	return u, nil
}

// Category returns a category with a given id
func (ur *UserRepositoryImpl) User(id int) (entity.User, error) {
	//row := ur.conn.QueryRow("SELECT * FROM users WHERE id = ?", id)
	row := ur.conn.QueryRow("SELECT * FROM users WHERE id = $1", id)

	u := entity.User{}

	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Password)

	//err = bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(u.Password))

	if err != nil {
		return u, err
	}

	return u, nil
}

// UpdateCategory updates a given object with a new data
func (ur *UserRepositoryImpl) EditAccount(u entity.User) error {

	log.Println("Edit Account", u.ID)
	log.Println("Edit Account", u.Name)
	log.Println("Edit Account", u.Email)
	log.Println("Edit Account", u.Phone)

	_, err := ur.conn.Exec("UPDATE users SET name=$1,email=$2,phone=$3 WHERE id=$4",
		u.Name, u.Email, u.Phone, u.ID)
	if err != nil {
		return errors.New("Update has failed")
	}

	return nil
}

// DeleteCategory removes a category from a database by its id
func (ur *UserRepositoryImpl) DeleteUser(id int) error {

	//_, err := ur.conn.Exec("DELETE FROM users WHERE id=?", id)
	_, err := ur.conn.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return errors.New("Delete has failed")
	}
	return nil
}

// StoreCategory stores new category information to database
func (ur *UserRepositoryImpl) StoreUser(u entity.User) error {
	var useremail string

	//err := ur.conn.QueryRow("SELECT email FROM users WHERE email=?", u.Email).Scan(&useremail)
	row1 := ur.conn.QueryRow("SELECT email FROM users WHERE email=$1", u.Email)
	err := row1.Scan(&useremail)

	switch {
	case err == sql.ErrNoRows:
		hashedpass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			panic(err.Error())
		}
		//_, err = ur.conn.Exec("INSERT INTO users (name,email,password)" +
		//	" values(?, ?, ?)",u.Name, u.Email, u.Password)
		_, err = ur.conn.Exec("INSERT INTO users (name,email,password)"+
			" values($1, $2, $3)", u.Name, u.Email, hashedpass)
		if err != nil {
			//panic(err)
			return errors.New("Insertion has failed")
		}
		return nil
		//info := entity.User{
		//	Name:     name,
		//	Email:    email,
		//	Password: pass,
		//}
		//
		//_ = tpl.ExecuteTemplate(w, "registered.html", info)
		//
		//log.Println("Insert: Name: " + name + " | Email: " + email + " | Password: " + string(hashedpass))
		//log.Println(err)
		//
		//log.Println("user" + name + "created!")
	case err != sql.ErrNoRows:
		return errors.New("This email is already taken")
	case err != nil:
		return errors.New("Server error, unable to create your account.")
	default:
		return errors.New("Server error, unable to create your account.")
	}
}

// UpdateCategory updates a given object with a new data
func (ur *UserRepositoryImpl) ChangePassword(u entity.User) error {

	_, err := ur.conn.Exec("UPDATE users SET password=$1 WHERE id=$2", u.Password, u.ID)
	if err != nil {
		return errors.New("Changing password has failed")
	}

	return nil
}

//func (ur *UserRepositoryImpl) GetUser(req *http.Request) entity.User {
//
//	u := entity.User{}
//	// get cookie
//	c, err := req.Cookie("session")
//
//	if err != nil {
//		return u
//	}
//
//	// if the user exists already, get user
//	if em, ok := dbSessions[c.Value]; ok {
//		row := ur.conn.QueryRow("SELECT * FROM users WHERE email = ?", em)
//
//		err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Password)
//
//		//err := sn.conn.QueryRow("SELECT name, email, password, phone FROM user WHERE email=?", em).Scan(u.Name, u.Email, u.Password, u.Phone)
//		if err!=nil{
//			panic(err.Error())
//		}
//	}
//	return u
//}
//
//func (ur *PsqlUserRepository) AlreadyLoggedIn(req *http.Request) bool {
//	var dbEmail string
//
//	c, err := req.Cookie("session")
//	if err != nil {
//		return false
//	}
//	em := dbSessions[c.Value]
//
//	err = ur.conn.QueryRow("SELECT email FROM users WHERE email=?", em).Scan(&dbEmail)
//
//	if dbEmail==em{
//		return true
//	}
//	if err!=nil{
//		panic(err.Error())
//	}
//	return true
//}
