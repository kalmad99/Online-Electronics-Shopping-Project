package users

import "../entity"

type UserRepository interface {
	StoreUser(user entity.User) error
}
