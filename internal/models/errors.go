package models

import "errors"

var (
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrConfirmLinkNotExists = errors.New("confirm link not found")
	ErrUserNotFound         = errors.New("user not found")
)
