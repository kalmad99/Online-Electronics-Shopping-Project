package users

import "../entity"

type UserRepo interface {
	Registration(user entity.User) error
}
