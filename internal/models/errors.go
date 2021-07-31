package models

import "errors"

var (
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrConfirmLinkNotExists = errors.New("confirm link not found")
	ErrBadUser              = errors.New("invalid email or password")
	ErrUserNotActivated     = errors.New("user is not activated")
	ErrMaxLoggedIn          = errors.New("maximum logged in users")
	ErrBadToken             = errors.New("invalid token")
	ErrUserUnauthorized     = errors.New("user is unauthorized")
	ErrTokenExpired         = errors.New("token is expired")
)
