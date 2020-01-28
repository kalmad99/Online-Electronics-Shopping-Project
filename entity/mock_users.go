package entity

// UserMock mocks application user
var UserMock = User{
	ID:          1,
	Name:        "Mock User 01",
	Email:       "mockuser@example.com",
	Phone:       "0900000000",
	Password:    "P@$$w0rd",
	RoleID:      1,
	ItemsInCart: []Cart{},
}
