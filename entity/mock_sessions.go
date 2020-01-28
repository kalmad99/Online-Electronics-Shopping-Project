package entity

// SessionMock mocks sessions of loged in user
var SessionMock = Session{
	ID:         1,
	UUID:       "_session_one",
	SigningKey: []byte("E&E"),
	Expires:    0,
}
