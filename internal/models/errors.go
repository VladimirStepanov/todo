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
	ErrNoAuthHeader         = errors.New("no authorization header")
	ErrInvalidAuthHeader    = errors.New("invalid authorization header")
	ErrNoList               = errors.New("list not found")
	ErrNoItem               = errors.New("item not found")
	ErrBadParam             = errors.New("bad url parameter")
	ErrNoListAccess         = errors.New("no access to this list")
	ErrUserNotFound         = errors.New("user not found")
	ErrUpdateEmptyArgs      = errors.New("empty title and description")
	ErrTitleTooShort        = errors.New("title too short. min length is 5")
)
