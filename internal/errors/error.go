package errors

import "errors"

var (
	ErrInvalidToken          = errors.New("invalid token")
	ErrAccessDenied          = errors.New("access denied")
	ErrUserNotFound          = errors.New("user not found")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
)
