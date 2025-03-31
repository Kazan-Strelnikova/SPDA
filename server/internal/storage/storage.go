package storage

import "errors"

var (
	ErrorNoUser        = errors.New("no user with this email")
	ErrorUserExists    = errors.New("user with this email already exists")
	ErrorEventNotFound = errors.New("no event with this id")
	ErrorLocationNotFound = errors.New("no location in cache")
)
