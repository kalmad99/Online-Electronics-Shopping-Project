package users

import "../entity"

type UserServices interface {
	Registration(user entity.User) error
}
