package models

import "errors"

var (
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrConfirmLinkNotExists = errors.New("confirm link not found")
	ErrBadUser              = errors.New("invalid email or password")
)
