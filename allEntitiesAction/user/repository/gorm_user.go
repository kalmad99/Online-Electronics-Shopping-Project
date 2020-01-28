package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/user"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// UserGormRepo Implements the menu.UserRepository interface
type UserGormRepo struct {
	conn *gorm.DB
}

// NewUserGormRepo creates a new object of UserGormRepo
func NewUserGormRepo(db *gorm.DB) user.UserRepository {
	return &UserGormRepo{conn: db}
}

// Users return all users from the database
func (userRepo *UserGormRepo) Users() ([]entity.User, []error) {
	users := []entity.User{}
	errs := userRepo.conn.Find(&users).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return users, errs
}
func (userRepo *UserGormRepo) Login(email string) (*entity.User, []error) {
	log.Println(email)

	u := entity.User{}

	errs := userRepo.conn.First(&u, &entity.User{Email: email}).GetErrors()

	if len(errs) > 0 {
		return nil, errs
	}
	return &u, errs
}

// User retrieves a user by its id from the database
func (userRepo *UserGormRepo) User(id uint) (*entity.User, []error) {
	usr := entity.User{}
	errs := userRepo.conn.First(&usr, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &usr, errs
}

// UpdateUser updates a given user in the database
func (userRepo *UserGormRepo) UpdateUser(user *entity.User) (*entity.User, []error) {
	usr := user
	errs := userRepo.conn.Model(&user).Updates(entity.User{Name: usr.Name, Email: usr.Email, Phone: usr.Phone}).GetErrors()
	//errs := userRepo.conn.Save(usr).GetErrors()

	//errs := userRepo.conn.Exec("UPDATE users SET name=$1, email=$2, phone=$3 WHERE id=$4;",
	//	usr.Name, usr.Email, usr.Phone, usr.ID).Save(usr).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}

// DeleteUser deletes a given user from the database
func (userRepo *UserGormRepo) DeleteUser(id uint) (*entity.User, []error) {
	usr, errs := userRepo.User(id)
	if len(errs) > 0 {
		return nil, errs
	}
	errs = userRepo.conn.Delete(usr, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}

// StoreUser stores a new user into the database
func (userRepo *UserGormRepo) StoreUser(user *entity.User) (*entity.User, []error) {
	usr := user
	hashedpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err.Error())
	}
	user.Password = string(hashedpass)

	errs := userRepo.conn.Create(usr).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}

func (userRepo *UserGormRepo) ChangePassword(user *entity.User) (*entity.User, []error) {
	usr := user
	//hashedpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	//if err != nil {
	//	panic(err.Error())
	//}
	errs := userRepo.conn.Model(&user).Updates(entity.User{Password: usr.Password}).GetErrors()

	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}

// PhoneExists check if a given phone number is found
func (userRepo *UserGormRepo) PhoneExists(phone string) bool {
	user := entity.User{}
	errs := userRepo.conn.Find(&user, "phone=?", phone).GetErrors()
	if len(errs) > 0 {
		return false
	}
	return true
}

// EmailExists check if a given email is found
func (userRepo *UserGormRepo) EmailExists(email string) bool {
	user := entity.User{}
	errs := userRepo.conn.Find(&user, "email=?", email).GetErrors()
	if len(errs) > 0 {
		return false
	}
	return true
}

// BankExists check if a given account is found
func (userRepo *UserGormRepo) BankExists(acc string) bool {
	//user := entity.User{}
	//bank := entity.Bank{}
	/*errs := userRepo.conn.Find(&bank, "accountno=?", acc).GetErrors()
	if len(errs) > 0 {
		return false
	}
	return true*/
	var bankacc string
	row := userRepo.conn.Table("bank").Where("accountno = ?", acc).Select("accountno").Row()
	err := row.Scan(&bankacc)
	if err != nil {
		return false
	}
	return true
}

// UserRoles returns list of application roles that a given user has
func (userRepo *UserGormRepo) UserRoles(user *entity.User) ([]entity.Role, []error) {
	userRoles := []entity.Role{}
	errs := userRepo.conn.Model(user).Related(&userRoles).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return userRoles, errs
}
