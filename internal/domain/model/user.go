package model

import "errors"

// ErrUserAlreadyExists error of user already exists
var ErrUserAlreadyExists = errors.New("user already exists")

// User User model
type User struct {
	ID       string
	Password string
	Nickname string
	Comment  string
}
