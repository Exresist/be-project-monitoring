package errors

import "errors"

// TODO:
/*type InternalError struct {
	Code int
	Err  error
}

func (i *InternalError) Error() string {
	return i.Err.Error()
}*/

var (
	ErrInvalidToken                = errors.New("invalid token")
	ErrAccessDenied                = errors.New("access denied")
	ErrUserNotFound                = errors.New("user not found")
	ErrEmailAlreadyExists          = errors.New("email already exists")
	ErrUsernameAlreadyExists       = errors.New("username already exists")
	ErrGithubUsernameAlreadyExists = errors.New("github username already exists")
)
