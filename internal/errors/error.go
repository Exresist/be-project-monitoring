package errors

import "errors"

var (
	ErrAbnormalExit = errors.New("abnormal exit")
	ErrTermSig      = errors.New("terminated with signal")
	ErrInvalidToken = errors.New("invalid token")
	ErrUserNotFound = errors.New("user not found")
)
