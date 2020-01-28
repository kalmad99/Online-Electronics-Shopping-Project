package user

import "github.com/kalmad99/Online-Electronics-Shopping-Project/entity"

// UserService specifies application user related services
type UserService interface {
	Users() ([]entity.User, []error)
	User(id uint) (*entity.User, []error)
	Login(email string) (*entity.User, []error)
	UpdateUser(user *entity.User) (*entity.User, []error)
	DeleteUser(id uint) (*entity.User, []error)
	StoreUser(user *entity.User) (*entity.User, []error)
	ChangePassword(user *entity.User) (*entity.User, []error)
	PhoneExists(phone string) bool
	EmailExists(email string) bool
	BankExists(acc string) bool
	UserRoles(*entity.User) ([]entity.Role, []error)
}

// RoleService speifies application user role related services
type RoleService interface {
	Roles() ([]entity.Role, []error)
	Role(id uint) (*entity.Role, []error)
	RoleByName(name string) (*entity.Role, []error)
	UpdateRole(role *entity.Role) (*entity.Role, []error)
	DeleteRole(id uint) (*entity.Role, []error)
	StoreRole(role *entity.Role) (*entity.Role, []error)
}

// SessionService specifies logged in user session related service
type SessionService interface {
	Session(sessionID string) (*entity.Session, []error)
	StoreSession(session *entity.Session) (*entity.Session, []error)
	DeleteSession(sessionID string) (*entity.Session, []error)
}
