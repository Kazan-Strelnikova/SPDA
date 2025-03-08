package storage

import "errors"

var (
	ErrorNoUser = errors.New("no user with this email")
	ErrorUserExists = errors.New("user with this email already exists")
)