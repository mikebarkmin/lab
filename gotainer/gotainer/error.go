package gotainer

import "errors"

var (
	ErrImagesDirUnreadable = errors.New("images dir is unreadable")
	ErrImageNotExist       = errors.New("image does not exist")
)
